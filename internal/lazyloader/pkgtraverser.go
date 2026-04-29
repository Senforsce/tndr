package lazyloader

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	lsp "github.com/senforsce/tndr/lsp/protocol"
	"github.com/senforsce/tndr/lsp/uri"
	"golang.org/x/tools/go/packages"
)

const (
	_tndrExt = ".t1"
)

type pkgTraverser interface {
	openTopologically(ctx context.Context, pkg *packages.Package) error
	closeTopologically(ctx context.Context, pkg *packages.Package) error
}

type goPkgTraverser struct {
	tndrDocHandler TndrDocHandler
	pkgsRefCount   map[string]int
	fileReader     fileReader
}

type TndrDocHandler interface {
	HandleDidOpen(ctx context.Context, params *lsp.DidOpenTextDocumentParams) error
	HandleDidClose(ctx context.Context, params *lsp.DidCloseTextDocumentParams) error
}

type fileReader interface {
	read(file string) ([]byte, error)
}

type tndrFileReader struct{}

func (tndrFileReader) read(file string) ([]byte, error) {
	return os.ReadFile(file)
}

func (t *goPkgTraverser) openTopologically(ctx context.Context, pkg *packages.Package) error {
	if t.pkgsRefCount[pkg.PkgPath] > 0 {
		t.pkgsRefCount[pkg.PkgPath]++
		return nil
	}

	for _, imp := range pkg.Imports {
		if err := t.openTopologically(ctx, imp); err != nil {
			return fmt.Errorf("open topologically %q: %w", imp.PkgPath, err)
		}
	}

	for _, otherFile := range pkg.OtherFiles {
		if filepath.Ext(otherFile) != _tndrExt {
			continue
		}

		text, err := t.fileReader.read(otherFile)
		if err != nil {
			return fmt.Errorf("read file %q: %w", otherFile, err)
		}

		if err := t.tndrDocHandler.HandleDidOpen(ctx, &lsp.DidOpenTextDocumentParams{
			TextDocument: lsp.TextDocumentItem{
				URI:        uri.File(otherFile),
				Text:       string(text),
				Version:    1,
				LanguageID: "go",
			},
		}); err != nil {
			return fmt.Errorf("did open file %q: %w", otherFile, err)
		}
	}
	t.pkgsRefCount[pkg.PkgPath]++

	return nil
}

func (t *goPkgTraverser) closeTopologically(ctx context.Context, pkg *packages.Package) error {
	if t.pkgsRefCount[pkg.PkgPath] > 1 {
		t.pkgsRefCount[pkg.PkgPath]--
		return nil
	}

	for _, otherFile := range pkg.OtherFiles {
		if filepath.Ext(otherFile) != _tndrExt {
			continue
		}

		if err := t.tndrDocHandler.HandleDidClose(ctx, &lsp.DidCloseTextDocumentParams{
			TextDocument: lsp.TextDocumentIdentifier{
				URI: uri.File(otherFile),
			},
		}); err != nil {
			return fmt.Errorf("did close file %q: %w", otherFile, err)
		}
	}
	delete(t.pkgsRefCount, pkg.PkgPath)

	for _, imp := range pkg.Imports {
		if err := t.closeTopologically(ctx, imp); err != nil {
			return fmt.Errorf("close topologically %q: %w", imp.PkgPath, err)
		}
	}

	return nil
}
