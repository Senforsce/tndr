package proxy

import (
	"path"
	"strings"

	lsp "github.com/a-h/protocol"
)

func convertT1ToGoURI(t1URI lsp.DocumentURI) (isT1File bool, goURI lsp.DocumentURI) {
	base, fileName := path.Split(string(t1URI))
	if !strings.HasSuffix(fileName, ".t1") {
		return
	}
	return true, lsp.DocumentURI(base + (strings.TrimSuffix(fileName, ".t1") + "_t1.go"))
}

func convertT1GoToT1URI(goURI lsp.DocumentURI) (isT1GoFile bool, t1URI lsp.DocumentURI) {
	base, fileName := path.Split(string(goURI))
	if !strings.HasSuffix(fileName, "_t1.go") {
		return
	}
	return true, lsp.DocumentURI(base + (strings.TrimSuffix(fileName, "_t1.go") + ".t1"))
}
