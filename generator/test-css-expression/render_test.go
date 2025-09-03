package testcssexpression

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/senforsce/tndr"
)

var expected = tndr.ComponentCSSClass{
	ID:    "className_34fc0328",
	Class: tndr.SafeCSS(`.className_34fc0328{background-color:#ffffff;max-height:calc(100vh - 170px);color:#ff0000;}`),
}

func TestCSSExpression(t *testing.T) {
	if diff := cmp.Diff(expected, className()); diff != "" {
		t.Error(diff)
	}
}
