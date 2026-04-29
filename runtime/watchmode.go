package runtime

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var developmentMode = os.Getenv("TNDR_DEV_MODE") == "true"

func GetDevModeTextFileName(tndrFilename string) string {
	if strings.HasSuffix(tndrFilename, "_t1.go") {
		tndrFilename = strings.TrimSuffix(tndrFilename, "_t1.go") + ".t1"
	}
	absFileName, err := filepath.Abs(tndrFilename)
	if err != nil {
		absFileName = tndrFilename
	}
	absFileName, err = filepath.EvalSymlinks(absFileName)
	if err != nil {
		absFileName = tndrFilename
	}
	absFileName = normalizePath(absFileName)

	hashedFileName := sha256.Sum256([]byte(absFileName))
	outputFileName := fmt.Sprintf("t1_%s.txt", hex.EncodeToString(hashedFileName[:]))

	root := os.TempDir()
	if os.Getenv("TNDR_DEV_MODE_ROOT") != "" {
		root = os.Getenv("TNDR_DEV_MODE_ROOT")
	}

	return filepath.Join(root, outputFileName)
}

// normalizePath converts Windows paths to Unix style paths.
func normalizePath(p string) string {
	p = strings.ReplaceAll(filepath.Clean(p), `\`, `/`)
	parts := strings.SplitN(p, ":", 2)
	if len(parts) == 2 && len(parts[0]) == 1 {
		drive := strings.ToLower(parts[0])
		p = "/" + drive + parts[1]
	}
	return p
}

// WriteString writes the string to the writer. If development mode is enabled
// s is replaced with the string at the index in the _tndr.txt file.
func WriteString(w io.Writer, index int, s string) (err error) {
	if developmentMode {
		_, path, _, _ := runtime.Caller(1)
		if !strings.HasSuffix(path, "_t1.go") {
			return errors.New("tndr: attempt to use WriteString from a non t1 file")
		}
		path, err := filepath.EvalSymlinks(path)
		if err != nil {
			return fmt.Errorf("tndr: failed to eval symlinks for %q: %w", path, err)
		}

		txtFilePath := GetDevModeTextFileName(path)
		literals, err := getWatchedStrings(txtFilePath)
		if err != nil {
			return fmt.Errorf("tndr: failed to get watched strings for %q: %w", path, err)
		}
		if index > len(literals) {
			return fmt.Errorf("tndr: failed to find line %d in %s", index, txtFilePath)
		}

		s, err = strconv.Unquote(`"` + literals[index-1] + `"`)
		if err != nil {
			return err
		}
	}
	_, err = io.WriteString(w, s)
	return err
}

var (
	watchModeCache  = map[string]watchState{}
	watchStateMutex sync.Mutex
)

type watchState struct {
	modTime time.Time
	strings []string
}

func getWatchedStrings(txtFilePath string) ([]string, error) {
	watchStateMutex.Lock()
	defer watchStateMutex.Unlock()

	state, cached := watchModeCache[txtFilePath]
	if !cached {
		return cacheStrings(txtFilePath)
	}

	if time.Since(state.modTime) < time.Millisecond*100 {
		return state.strings, nil
	}

	info, err := os.Stat(txtFilePath)
	if err != nil {
		return nil, fmt.Errorf("tndr: failed to stat %s: %w", txtFilePath, err)
	}

	if !info.ModTime().After(state.modTime) {
		return state.strings, nil
	}

	return cacheStrings(txtFilePath)
}

func cacheStrings(txtFilePath string) ([]string, error) {
	txtFile, err := os.Open(txtFilePath)
	if err != nil {
		return nil, fmt.Errorf("tndr: failed to open %s: %w", txtFilePath, err)
	}
	defer func() {
		_ = txtFile.Close()
	}()

	info, err := txtFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("tndr: failed to stat %s: %w", txtFilePath, err)
	}

	all, err := io.ReadAll(txtFile)
	if err != nil {
		return nil, fmt.Errorf("tndr: failed to read %s: %w", txtFilePath, err)
	}

	literals := strings.Split(string(all), "\n")
	watchModeCache[txtFilePath] = watchState{
		modTime: info.ModTime(),
		strings: literals,
	}

	return literals, nil
}
