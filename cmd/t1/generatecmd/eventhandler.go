package generatecmd

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"go/format"
	"go/scanner"
	"go/token"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/senforsce/level0/generator"
	tf "github.com/senforsce/level0/templatefile"

	"github.com/senforsce/tndr/cmd/t1/visualize"
	"github.com/senforsce/toolbelt/sourcemap"
	"github.com/senforsce/toolbelt/templatefile"
)

func NewFSEventHandler(
	log *slog.Logger,
	dir string,
	devMode bool,
	genOpts []generator.GenerateOpt,
	genSourceMapVis bool,
	keepOrphanedFiles bool,
	toStdout bool,
) *FSEventHandler {
	if !path.IsAbs(dir) {
		dir, _ = filepath.Abs(dir)
	}
	fseh := &FSEventHandler{
		Log:                        log,
		dir:                        dir,
		fileNameToLastModTime:      make(map[string]time.Time),
		fileNameToLastModTimeMutex: &sync.Mutex{},
		fileNameToError:            make(map[string]struct{}),
		fileNameToErrorMutex:       &sync.Mutex{},
		hashes:                     make(map[string][sha256.Size]byte),
		hashesMutex:                &sync.Mutex{},
		genOpts:                    genOpts,
		genSourceMapVis:            genSourceMapVis,
		DevMode:                    devMode,
		keepOrphanedFiles:          keepOrphanedFiles,
		writer:                     writeToFile,
	}
	if toStdout {
		fseh.writer = writeToStdout
	}
	if devMode {
		fseh.genOpts = append(fseh.genOpts, generator.WithExtractStrings())
	}
	return fseh
}

type FSEventHandler struct {
	Log *slog.Logger
	// dir is the root directory being processed.
	dir                        string
	fileNameToLastModTime      map[string]time.Time
	fileNameToLastModTimeMutex *sync.Mutex
	fileNameToError            map[string]struct{}
	fileNameToErrorMutex       *sync.Mutex
	hashes                     map[string][sha256.Size]byte
	hashesMutex                *sync.Mutex
	genOpts                    []generator.GenerateOpt
	genSourceMapVis            bool
	DevMode                    bool
	Errors                     []error
	keepOrphanedFiles          bool
	writer                     func(string, []byte) error
}

func writeToFile(fileName string, contents []byte) error {
	return os.WriteFile(fileName, contents, 0o644)
}

func writeToStdout(_ string, contents []byte) error {
	_, err := os.Stdout.Write(contents)
	return err
}

func (h *FSEventHandler) HandleEvent(ctx context.Context, event fsnotify.Event) (goUpdated, textUpdated bool, err error) {
	// Handle _t1.go files.
	if !event.Has(fsnotify.Remove) && strings.HasSuffix(event.Name, "_t1.go") {
		_, err = os.Stat(strings.TrimSuffix(event.Name, "_t1.go") + ".t1")
		if !os.IsNotExist(err) {
			return false, false, err
		}
		// File is orphaned.
		if h.keepOrphanedFiles {
			return false, false, nil
		}
		h.Log.Debug("Deleting orphaned Go file", slog.String("file", event.Name))
		if err = os.Remove(event.Name); err != nil {
			h.Log.Warn("Failed to remove orphaned file", slog.Any("error", err))
		}
		return true, false, nil
	}
	// Handle _t1.txt files.
	if !event.Has(fsnotify.Remove) && strings.HasSuffix(event.Name, "_t1.txt") {
		if h.DevMode {
			// Don't delete the file if we're in dev mode, but mark that text was updated.
			return false, true, nil
		}
		h.Log.Debug("Deleting watch mode file", slog.String("file", event.Name))
		if err = os.Remove(event.Name); err != nil {
			h.Log.Warn("Failed to remove watch mode text file", slog.Any("error", err))
			return false, false, nil
		}
		return false, false, nil
	}

	// Handle .t1 files.
	if !strings.HasSuffix(event.Name, ".t1") {
		return false, false, nil
	}

	// If the file hasn't been updated since the last time we processed it, ignore it.
	if !h.UpsertLastModTime(event.Name) {
		h.Log.Debug("Skipping file because it wasn't updated", slog.String("file", event.Name))
		return false, false, nil
	}

	// Start a processor.
	start := time.Now()
	goUpdated, textUpdated, diag, err := h.generate(ctx, event.Name)
	if err != nil {
		h.Log.Error(
			"Error generating code",
			slog.String("file", event.Name),
			slog.Any("error", err),
		)
		h.SetError(event.Name, true)
		return goUpdated, textUpdated, fmt.Errorf("failed to generate code for %q: %w", event.Name, err)
	}
	if len(diag) > 0 {
		for _, d := range diag {
			h.Log.Warn(d.Message,
				slog.String("from", fmt.Sprintf("%d:%d", d.Range.From.Line, d.Range.From.Col)),
				slog.String("to", fmt.Sprintf("%d:%d", d.Range.To.Line, d.Range.To.Col)),
			)
		}
		return
	}
	if errorCleared, errorCount := h.SetError(event.Name, false); errorCleared {
		h.Log.Info("Error cleared", slog.String("file", event.Name), slog.Int("errors", errorCount))
	}
	h.Log.Debug("Generated code", slog.String("file", event.Name), slog.Duration("in", time.Since(start)))

	return goUpdated, textUpdated, nil
}

func (h *FSEventHandler) SetError(fileName string, hasError bool) (previouslyHadError bool, errorCount int) {
	h.fileNameToErrorMutex.Lock()
	defer h.fileNameToErrorMutex.Unlock()
	_, previouslyHadError = h.fileNameToError[fileName]
	delete(h.fileNameToError, fileName)
	if hasError {
		h.fileNameToError[fileName] = struct{}{}
	}
	return previouslyHadError, len(h.fileNameToError)
}

func (h *FSEventHandler) UpsertLastModTime(fileName string) (updated bool) {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return false
	}
	h.fileNameToLastModTimeMutex.Lock()
	defer h.fileNameToLastModTimeMutex.Unlock()
	lastModTime := h.fileNameToLastModTime[fileName]
	if !fileInfo.ModTime().After(lastModTime) {
		return false
	}
	h.fileNameToLastModTime[fileName] = fileInfo.ModTime()
	return true
}

func (h *FSEventHandler) UpsertHash(fileName string, hash [sha256.Size]byte) (updated bool) {
	h.hashesMutex.Lock()
	defer h.hashesMutex.Unlock()
	lastHash := h.hashes[fileName]
	if lastHash == hash {
		return false
	}
	h.hashes[fileName] = hash
	return true
}

// generate Go code for a single template.
// If a basePath is provided, the filename included in error messages is relative to it.
func (h *FSEventHandler) generate(ctx context.Context, fileName string) (goUpdated, textUpdated bool, diagnostics []templatefile.Diagnostic, err error) {
	tfp := tf.NewTemplateFileParser("main")
	t, err := templatefile.ReadByNameAndParsers(fileName, tfp, tf.TemplateNodeParserList)

	if err != nil {
		return false, false, nil, fmt.Errorf("%s parsing error: %w", fileName, err)
	}
	targetFileName := strings.TrimSuffix(fileName, ".t1") + "_t1.go"

	// Only use relative filenames to the basepath for filenames in runtime error messages.
	absFilePath, err := filepath.Abs(fileName)
	if err != nil {
		return false, false, nil, fmt.Errorf("failed to get absolute path for %q: %w", fileName, err)
	}
	relFilePath, err := filepath.Rel(h.dir, absFilePath)
	if err != nil {
		return false, false, nil, fmt.Errorf("failed to get relative path for %q: %w", fileName, err)
	}
	// Convert Windows file paths to Unix-style for consistency.
	relFilePath = filepath.ToSlash(relFilePath)

	var b bytes.Buffer
	sourceMap, literals, err := generator.Generate(t, &b, append(h.genOpts, generator.WithFileName(relFilePath))...)
	if err != nil {
		return false, false, nil, fmt.Errorf("%s generation error: %w", fileName, err)
	}

	formattedGoCode, err := format.Source(b.Bytes())
	if err != nil {
		err = remapErrorList(err, sourceMap, fileName, targetFileName)
		return false, false, nil, fmt.Errorf("%s source formatting error %w => %s", fileName, err, b.String())
	}

	// Hash output, and write out the file if the goCodeHash has changed.
	goCodeHash := sha256.Sum256(formattedGoCode)
	if h.UpsertHash(targetFileName, goCodeHash) {
		goUpdated = true
		if err = h.writer(targetFileName, formattedGoCode); err != nil {
			return false, false, nil, fmt.Errorf("failed to write target file %q: %w", targetFileName, err)
		}
	}

	// Add the txt file if it has changed.
	if len(literals) > 0 {
		txtFileName := strings.TrimSuffix(fileName, ".t1") + "_t1.txt"
		txtHash := sha256.Sum256([]byte(literals))
		if h.UpsertHash(txtFileName, txtHash) {
			textUpdated = true
			if err = os.WriteFile(txtFileName, []byte(literals), 0o644); err != nil {
				return false, false, nil, fmt.Errorf("failed to write string literal file %q: %w", txtFileName, err)
			}
		}
	}

	parsedDiagnostics, err := templatefile.Diagnose(t)
	if err != nil {
		return goUpdated, textUpdated, nil, fmt.Errorf("%s diagnostics error: %w", fileName, err)
	}

	if h.genSourceMapVis {
		err = generateSourceMapVisualisation(ctx, fileName, targetFileName, sourceMap)
	}

	return goUpdated, textUpdated, parsedDiagnostics, err
}

// Takes an error from the formatter and attempts to convert the positions reported in the target file to their positions
// in the source file.
func remapErrorList(err error, sourceMap *sourcemap.SourceMap, fileName string, targetFileName string) error {
	list, ok := err.(scanner.ErrorList)
	if !ok || len(list) == 0 {
		return err
	}
	for i, e := range list {
		// The positions in the source map are off by one line because of the package definition.
		srcPos, ok := sourceMap.SourcePositionFromTarget(int(e.Pos.Line-1), int(e.Pos.Column))
		if !ok {
			continue
		}
		list[i].Pos = token.Position{
			Filename: fileName,
			Offset:   int(srcPos.Index),
			Line:     int(srcPos.Line) + 1,
			Column:   int(srcPos.Col),
		}
	}
	return list
}

func generateSourceMapVisualisation(ctx context.Context, t1FileName, goFileName string, sourceMap *sourcemap.SourceMap) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	var t1Contents, goContents []byte
	var t1Err, goErr error
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		t1Contents, t1Err = os.ReadFile(t1FileName)
	}()
	go func() {
		defer wg.Done()
		goContents, goErr = os.ReadFile(goFileName)
	}()
	wg.Wait()
	if t1Err != nil {
		return t1Err
	}
	if goErr != nil {
		return t1Err
	}

	targetFileName := strings.TrimSuffix(t1FileName, ".t1") + "_t1_sourcemap.html"
	w, err := os.Create(targetFileName)
	if err != nil {
		return fmt.Errorf("%s sourcemap visualisation error: %w", t1FileName, err)
	}
	defer w.Close()
	b := bufio.NewWriter(w)
	defer b.Flush()

	return visualize.HTML(t1FileName, string(t1Contents), string(goContents), sourceMap).Render(ctx, b)
}
