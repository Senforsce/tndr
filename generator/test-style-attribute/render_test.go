package teststyleattribute

import (
	_ "embed"
	"fmt"
	"os"
	"testing"

	"github.com/senforsce/tndr"
	"github.com/senforsce/tndr/generator/htmldiff"
)

//go:embed expected.html
var expected string

func Test(t *testing.T) {
	var stringCSS = "background-color:blue;color:red"
	var safeCSS = tndr.SafeCSS("background-color:blue;color:red;")
	var mapStringString = map[string]string{
		"color":            "red",
		"background-color": "blue",
	}
	var mapStringSafeCSSProperty = map[string]tndr.SafeCSSProperty{
		"color":            tndr.SafeCSSProperty("red"),
		"background-color": tndr.SafeCSSProperty("blue"),
	}
	var kvStringStringSlice = []tndr.KeyValue[string, string]{
		tndr.KV("background-color", "blue"),
		tndr.KV("color", "red"),
	}
	var kvStringBoolSlice = []tndr.KeyValue[string, bool]{
		tndr.KV("background-color:blue", true),
		tndr.KV("color:red", true),
		tndr.KV("color:blue", false),
	}
	var kvSafeCSSBoolSlice = []tndr.KeyValue[tndr.SafeCSS, bool]{
		tndr.KV(tndr.SafeCSS("background-color:blue"), true),
		tndr.KV(tndr.SafeCSS("color:red"), true),
		tndr.KV(tndr.SafeCSS("color:blue"), false),
	}

	tests := []any{
		stringCSS,
		safeCSS,
		mapStringString,
		mapStringSafeCSSProperty,
		kvStringStringSlice,
		kvStringBoolSlice,
		kvSafeCSSBoolSlice,
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%T", test), func(t *testing.T) {
			component := Button(test, "Click me")

			actual, diff, err := htmldiff.Diff(component, expected)
			if err != nil {
				t.Fatal(err)
			}
			if diff != "" {
				if err := os.WriteFile("actual.html", []byte(actual), 0644); err != nil {
					t.Errorf("failed to write actual.html: %v", err)
				}
				t.Error(diff)
			}
		})
	}
}
