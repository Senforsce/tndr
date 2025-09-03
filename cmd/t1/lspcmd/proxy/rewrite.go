package proxy

import (
	"path"
	"strings"

	lsp "github.com/senforsce/tndr/lsp/protocol"
)

func convertTndrToGoURI(templURI lsp.DocumentURI) (isTemplFile bool, goURI lsp.DocumentURI) {
	base, fileName := path.Split(string(templURI))
	if !strings.HasSuffix(fileName, ".t1") {
		return
	}
	return true, lsp.DocumentURI(base + (strings.TrimSuffix(fileName, ".t1") + "_t1.go"))
}

func convertTndrGoToTndrURI(goURI lsp.DocumentURI) (isTemplGoFile bool, t1URI lsp.DocumentURI) {
	base, fileName := path.Split(string(goURI))
	if !strings.HasSuffix(fileName, "_t1.go") {
		return
	}
	return true, lsp.DocumentURI(base + (strings.TrimSuffix(fileName, "_t1.go") + ".t1"))
}
