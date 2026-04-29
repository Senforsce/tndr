package proxy

import (
	"sync"

	lsp "github.com/senforsce/tndr/lsp/protocol"
)

func NewDiagnosticCache() *DiagnosticCache {
	return &DiagnosticCache{
		m:     &sync.Mutex{},
		cache: make(map[string]fileDiagnostic),
	}
}

type fileDiagnostic struct {
	tndrDiagnostics  []lsp.Diagnostic
	goplsDiagnostics []lsp.Diagnostic
}

type DiagnosticCache struct {
	m     *sync.Mutex
	cache map[string]fileDiagnostic
}

func zeroLengthSliceIfNil(diags []lsp.Diagnostic) []lsp.Diagnostic {
	if diags == nil {
		return make([]lsp.Diagnostic, 0)
	}
	return diags
}

func (dc *DiagnosticCache) AddTndrDiagnostics(uri string, goDiagnostics []lsp.Diagnostic) []lsp.Diagnostic {
	goDiagnostics = zeroLengthSliceIfNil(goDiagnostics)
	dc.m.Lock()
	defer dc.m.Unlock()
	diag := dc.cache[uri]
	diag.goplsDiagnostics = goDiagnostics
	diag.tndrDiagnostics = zeroLengthSliceIfNil(diag.tndrDiagnostics)
	dc.cache[uri] = diag
	return append(diag.tndrDiagnostics, goDiagnostics...)
}

func (dc *DiagnosticCache) ClearTndrDiagnostics(uri string) {
	dc.m.Lock()
	defer dc.m.Unlock()
	diag := dc.cache[uri]
	diag.tndrDiagnostics = make([]lsp.Diagnostic, 0)
	dc.cache[uri] = diag
}

func (dc *DiagnosticCache) AddGoDiagnostics(uri string, tndrDiagnostics []lsp.Diagnostic) []lsp.Diagnostic {
	tndrDiagnostics = zeroLengthSliceIfNil(tndrDiagnostics)
	dc.m.Lock()
	defer dc.m.Unlock()
	diag := dc.cache[uri]
	diag.tndrDiagnostics = tndrDiagnostics
	diag.goplsDiagnostics = zeroLengthSliceIfNil(diag.goplsDiagnostics)
	dc.cache[uri] = diag
	return append(diag.goplsDiagnostics, tndrDiagnostics...)
}
