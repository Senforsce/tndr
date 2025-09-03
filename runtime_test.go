package tndr_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/senforsce/tndr"
)

func TestCSSID(t *testing.T) {
	t.Run("minimum hash suffix length is 8", func(t *testing.T) {
		// See issue #978.
		name := "classA"
		css := "background-color:white;"
		actual := len(tndr.CSSID(name, css))
		expected := len(name) + 1 + 8
		if expected != actual {
			t.Errorf("expected length %d, got %d", expected, actual)
		}
	})
	t.Run("known hash collisions are avoided", func(t *testing.T) {
		name := "classA"
		// Note that the first 4 characters of the hash are the same.
		css1 := "grid-column:1;grid-row:1;"  // After hash: f781266f
		css2 := "grid-column:13;grid-row:6;" // After hash: f781f18b
		id1 := tndr.CSSID(name, css1)
		id2 := tndr.CSSID(name, css2)
		if id1 == id2 {
			t.Errorf("hash collision: %s == %s", id1, id2)
		}
	})
}

func TestCSSHandler(t *testing.T) {
	tests := []struct {
		name             string
		input            []tndr.CSSClass
		expectedMIMEType string
		expectedBody     string
	}{
		{
			name:             "no classes",
			input:            nil,
			expectedMIMEType: "text/css",
			expectedBody:     "",
		},
		{
			name:             "classes are rendered",
			input:            []tndr.CSSClass{tndr.ComponentCSSClass{ID: "className", Class: tndr.SafeCSS(".className{background-color:white;}")}},
			expectedMIMEType: "text/css",
			expectedBody:     ".className{background-color:white;}",
		},
		{
			name: "classes are rendered",
			input: []tndr.CSSClass{
				tndr.ComponentCSSClass{ID: "classA", Class: tndr.SafeCSS(".classA{background-color:white;}")},
				tndr.ComponentCSSClass{ID: "classB", Class: tndr.SafeCSS(".classB{background-color:green;}")},
			},
			expectedMIMEType: "text/css",
			expectedBody:     ".classA{background-color:white;}.classB{background-color:green;}",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			h := tndr.NewCSSHandler(tt.input...)
			h.ServeHTTP(w, &http.Request{})
			if diff := cmp.Diff(tt.expectedMIMEType, w.Header().Get("Content-Type")); diff != "" {
				t.Error(diff)
			}
			if diff := cmp.Diff(tt.expectedBody, w.Body.String()); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestCSSMiddleware(t *testing.T) {
	pageHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.WriteString(w, "Hello, World!"); err != nil {
			t.Fatalf("failed to write string: %v", err)
		}
	})
	c1 := tndr.ComponentCSSClass{
		ID:    "c1",
		Class: ".c1{color:red}",
	}
	c2 := tndr.ComponentCSSClass{
		ID:    "c2",
		Class: ".c2{color:blue}",
	}

	tests := []struct {
		name             string
		input            *http.Request
		handler          http.Handler
		expectedMIMEType string
		expectedBody     string
	}{
		{
			name:             "accessing /style/tndr.css renders CSS, even if it's empty",
			input:            httptest.NewRequest("GET", "/styles/tndr.css", nil),
			handler:          tndr.NewCSSMiddleware(pageHandler),
			expectedMIMEType: "text/css",
			expectedBody:     "",
		},
		{
			name:             "accessing /style/tndr.css renders CSS that includes the classes",
			input:            httptest.NewRequest("GET", "/styles/tndr.css", nil),
			handler:          tndr.NewCSSMiddleware(pageHandler, c1, c2),
			expectedMIMEType: "text/css",
			expectedBody:     ".c1{color:red}.c2{color:blue}",
		},
		{
			name:             "the pageHandler is rendered",
			input:            httptest.NewRequest("GET", "/index.html", nil),
			handler:          tndr.NewCSSMiddleware(pageHandler, c1, c2),
			expectedMIMEType: "text/plain; charset=utf-8",
			expectedBody:     "Hello, World!",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.handler.ServeHTTP(w, tt.input)
			if diff := cmp.Diff(tt.expectedMIMEType, w.Header().Get("Content-Type")); diff != "" {
				t.Error(diff)
			}
			if diff := cmp.Diff(tt.expectedBody, w.Body.String()); diff != "" {
				t.Error(diff)
			}
		})
	}
}

var cssInputs = []any{
	[]string{"a", "b"},         // []string
	"c",                        // string
	tndr.ConstantCSSClass("d"), // ConstantCSSClass
	tndr.ComponentCSSClass{ID: "e", Class: ".e{color:red}"}, // ComponentCSSClass
	map[string]bool{"f": true, "ff": false},                 // map[string]bool
	tndr.KV("g", true),                                      // KeyValue[string, bool]
	tndr.KV("gg", false),                                    // KeyValue[string, bool]
	[]tndr.KeyValue[string, bool]{
		tndr.KV("h", true),
		tndr.KV("hh", false),
	}, // []KeyValue[string, bool]
	tndr.KV(tndr.ConstantCSSClass("i"), true),   // KeyValue[CSSClass, bool]
	tndr.KV(tndr.ConstantCSSClass("ii"), false), // KeyValue[CSSClass, bool]
	tndr.KV(tndr.ComponentCSSClass{
		ID:    "j",
		Class: ".j{color:red}",
	}, true), // KeyValue[ComponentCSSClass, bool]
	tndr.KV(tndr.ComponentCSSClass{
		ID:    "jj",
		Class: ".jj{color:red}",
	}, false), // KeyValue[ComponentCSSClass, bool]
	tndr.CSSClasses{tndr.ConstantCSSClass("k")},                // CSSClasses
	func() tndr.CSSClass { return tndr.ConstantCSSClass("l") }, // func() CSSClass
	tndr.CSSClass(tndr.ConstantCSSClass("m")),                  // CSSClass
	customClass{name: "n"},                                     // CSSClass
	[]tndr.CSSClass{customClass{name: "n"}},                    // []CSSClass
	tndr.KV(tndr.ConstantCSSClass("o"), true),                  // KeyValue[ConstantCSSClass, bool]
	[]tndr.KeyValue[tndr.ConstantCSSClass, bool]{
		tndr.KV(tndr.ConstantCSSClass("p"), true),
		tndr.KV(tndr.ConstantCSSClass("pp"), false),
	}, // []KeyValue[ConstantCSSClass, bool]
}

type customClass struct {
	name string
}

func (cc customClass) ClassName() string {
	return cc.name
}

func TestRenderCSS(t *testing.T) {
	c1 := tndr.ComponentCSSClass{
		ID:    "c1",
		Class: ".c1{color:red}",
	}
	c2 := tndr.ComponentCSSClass{
		ID:    "c2",
		Class: ".c2{color:blue}",
	}

	tests := []struct {
		name     string
		toIgnore []any
		toRender []any
		expected string
	}{
		{
			name:     "if none are ignored, everything is rendered",
			toIgnore: nil,
			toRender: []any{c1, c2},
			expected: `<style type="text/css">.c1{color:red}.c2{color:blue}</style>`,
		},
		{
			name: "if something outside the expected is ignored, if has no effect",
			toIgnore: []any{
				tndr.ComponentCSSClass{
					ID:    "c3",
					Class: tndr.SafeCSS(".c3{color:yellow}"),
				},
			},
			toRender: []any{c1, c2},
			expected: `<style type="text/css">.c1{color:red}.c2{color:blue}</style>`,
		},
		{
			name:     "if one is ignored, it's not rendered",
			toIgnore: []any{c1},
			toRender: []any{c1, c2},
			expected: `<style type="text/css">.c2{color:blue}</style>`,
		},
		{
			name: "if all are ignored, not even style tags are rendered",
			toIgnore: []any{
				c1,
				c2,
				tndr.ComponentCSSClass{
					ID:    "c3",
					Class: tndr.SafeCSS(".c3{color:yellow}"),
				},
			},
			toRender: []any{c1, c2},
			expected: ``,
		},
		{
			name:     "CSS classes are rendered",
			toIgnore: nil,
			toRender: cssInputs,
			expected: `<style type="text/css">.e{color:red}.j{color:red}</style>`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			b := new(bytes.Buffer)

			// Render twice, reusing the same context so that there's a memory of which classes have been rendered.
			ctx = tndr.InitializeContext(ctx)
			err := tndr.RenderCSSItems(ctx, b, tt.toIgnore...)
			if err != nil {
				t.Fatalf("failed to render initial CSS: %v", err)
			}

			// Now render again to check that only the expected classes were rendered.
			b.Reset()
			err = tndr.RenderCSSItems(ctx, b, tt.toRender...)
			if err != nil {
				t.Fatalf("failed to render CSS: %v", err)
			}

			if diff := cmp.Diff(tt.expected, b.String()); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestRenderCSSItemsWithNonce(t *testing.T) {
	ctx := tndr.WithNonce(context.Background(), "testnonce")
	b := new(bytes.Buffer)
	c := tndr.ComponentCSSClass{
		ID:    "c1",
		Class: ".c1{color:red}",
	}
	err := tndr.RenderCSSItems(ctx, b, c)
	if err != nil {
		t.Fatalf("failed to render CSS: %v", err)
	}
	actual := b.String()
	// Should include nonce attribute on <style> tag.
	expected := `<style type="text/css" nonce="testnonce">.c1{color:red}</style>`
	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestClassesFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    []any
		expected string
	}{
		{
			name:     "constants are allowed",
			input:    []any{"a", "b", "c", "</style>"},
			expected: "a b c </style>",
		},
		{
			name:     "legacy CSS types are supported",
			input:    []any{"a", tndr.SafeClass("b"), tndr.Class("c")},
			expected: "a b c",
		},
		{
			name: "CSS components are included in the output",
			input: []any{
				tndr.ComponentCSSClass{ID: "classA", Class: tndr.SafeCSS(".classA{background-color:white;}")},
				tndr.ComponentCSSClass{ID: "classB", Class: tndr.SafeCSS(".classB{background-color:green;}")},
				"c",
			},
			expected: "classA classB c",
		},
		{
			name: "optional classes can be applied with expressions",
			input: []any{
				"a",
				tndr.ComponentCSSClass{ID: "classA", Class: tndr.SafeCSS(".classA{background-color:white;}")},
				tndr.ComponentCSSClass{ID: "classB", Class: tndr.SafeCSS(".classB{background-color:green;}")},
				"c",
				map[string]bool{
					"a":      false,
					"classA": false,
					"classB": false,
					"c":      true,
					"d":      false,
				},
			},
			expected: "c",
		},
		{
			name: "unknown types for classes get rendered as --tndr-css-class-unknown-type",
			input: []any{
				123,
				map[string]string{"test": "no"},
				false,
				"c",
			},
			expected: "--tndr-css-class-unknown-type c",
		},
		{
			name: "string arrays are supported",
			input: []any{
				[]string{"a", "b", "c", "</style>"},
				"d",
			},
			expected: "a b c </style> d",
		},
		{
			name: "strings are broken up",
			input: []any{
				"a </style>",
			},
			expected: "a </style>",
		},
		{
			name: "if a tndr.CSSClasses is passed in, the nested CSSClasses are extracted",
			input: []any{
				tndr.Classes(
					"a",
					tndr.SafeClass("b"),
					tndr.Class("c"),
					tndr.ComponentCSSClass{
						ID:    "d",
						Class: "{}",
					},
				),
			},
			expected: "a b c d",
		},
		{
			name: "kv types can be used to show or hide classes",
			input: []any{
				"a",
				tndr.KV("b", true),
				"c",
				tndr.KV("c", false),
				tndr.KV(tndr.SafeClass("d"), true),
				tndr.KV(tndr.SafeClass("e"), false),
			},
			expected: "a b d",
		},
		{
			name: "an array of KV types can be used to show or hide classes",
			input: []any{
				"a",
				"c",
				[]tndr.KeyValue[string, bool]{
					tndr.KV("b", true),
					tndr.KV("c", false),
					{"d", true},
				},
			},
			expected: "a b d",
		},
		{
			name: "the brackets on component CSS function calls can be elided",
			input: []any{
				func() tndr.CSSClass {
					return tndr.ComponentCSSClass{
						ID:    "a",
						Class: "",
					}
				},
			},
			expected: "a",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := tndr.Classes(test.input...).String()
			if actual != test.expected {
				t.Errorf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}

type baseError struct {
	Value int
}

func (baseError) Error() string { return "base error" }

type nonMatchedError struct{}

func (nonMatchedError) Error() string { return "non matched error" }

func TestErrorWrapping(t *testing.T) {
	baseErr := baseError{
		Value: 1,
	}
	wrappedErr := tndr.Error{Err: baseErr, Line: 1, Col: 2}
	t.Run("errors.Is() returns true for the base error", func(t *testing.T) {
		if !errors.Is(wrappedErr, baseErr) {
			t.Error("errors.Is() returned false for the base error")
		}
	})
	t.Run("errors.Is() returns false for a different error", func(t *testing.T) {
		if errors.Is(wrappedErr, errors.New("different error")) {
			t.Error("errors.Is() returned true for a different error")
		}
	})
	t.Run("errors.As() returns true for the base error", func(t *testing.T) {
		var err baseError
		if !errors.As(wrappedErr, &err) {
			t.Error("errors.As() returned false for the base error")
		}
		if err.Value != 1 {
			t.Errorf("errors.As() returned a different value: %v", err.Value)
		}
	})
	t.Run("errors.As() returns false for a different error", func(t *testing.T) {
		var err nonMatchedError
		if errors.As(wrappedErr, &err) {
			t.Error("errors.As() returned true for a different error")
		}
	})
}

func TestRawComponent(t *testing.T) {
	tests := []struct {
		name        string
		input       tndr.Component
		expected    string
		expectedErr error
	}{
		{
			name:     "Raw content is not escaped",
			input:    tndr.Raw("<div>Test &</div>"),
			expected: `<div>Test &</div>`,
		},
		{
			name:        "Raw will return errors first",
			input:       tndr.Raw("", nil, errors.New("test error")),
			expected:    `<div>Test &</div>`,
			expectedErr: errors.New("test error"),
		},
		{
			name:     "Strings marked as safe are rendered without escaping",
			input:    tndr.Raw(template.HTML("<div>")),
			expected: `<div>`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			err := tt.input.Render(context.Background(), b)
			if tt.expectedErr != nil {
				expected := tt.expectedErr.Error()
				actual := fmt.Sprintf("%v", err)
				if actual != expected {
					t.Errorf("expected error %q, got %q", expected, actual)
				}
				return
			}
			if err != nil {
				t.Fatalf("failed to render content: %v", err)
			}
			if diff := cmp.Diff(tt.expected, b.String()); diff != "" {
				t.Error(diff)
			}
		})
	}
	t.Run("Raw does not require allocations", func(t *testing.T) {
		actualAllocs := testing.AllocsPerRun(4, func() {
			c := tndr.Raw("<div>")
			if c == nil {
				t.Fatalf("unexpected nil value")
			}
		})
		if actualAllocs > 0 {
			t.Errorf("expected no allocs, got %v", actualAllocs)
		}
	})
}

var goTemplate = template.Must(template.New("example").Parse("<div>{{ . }}</div>"))

func TestGoHTMLComponents(t *testing.T) {
	t.Run("Go templates can be rendered as tndr components", func(t *testing.T) {
		b := new(bytes.Buffer)
		err := tndr.FromGoHTML(goTemplate, "Test &").Render(context.Background(), b)
		if err != nil {
			t.Fatalf("failed to render content: %v", err)
		}
		if diff := cmp.Diff("<div>Test &amp;</div>", b.String()); diff != "" {
			t.Error(diff)
		}
	})
	t.Run("tndr components can be rendered in Go templates", func(t *testing.T) {
		b := new(bytes.Buffer)
		c := tndr.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			_, err = io.WriteString(w, "<div>Unsanitized &</div>")
			return err
		})
		h, err := tndr.ToGoHTML(context.Background(), c)
		if err != nil {
			t.Fatalf("failed to convert to Go HTML: %v", err)
		}
		if err = goTemplate.Execute(b, h); err != nil {
			t.Fatalf("failed to render content: %v", err)
		}
		if diff := cmp.Diff("<div><div>Unsanitized &</div></div>", b.String()); diff != "" {
			t.Error(diff)
		}
	})
	t.Run("errors in ToGoHTML are returned", func(t *testing.T) {
		expectedErr := errors.New("test error")
		c := tndr.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			return expectedErr
		})
		_, err := tndr.ToGoHTML(context.Background(), c)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if err != expectedErr {
			t.Fatalf("expected error %q, got %q", expectedErr, err)
		}
	})
	t.Run("FromGoHTML does not require allocations", func(t *testing.T) {
		actualAllocs := testing.AllocsPerRun(4, func() {
			c := tndr.FromGoHTML(goTemplate, "test &")
			if c == nil {
				t.Fatalf("unexpected nil value")
			}
		})
		if actualAllocs > 0 {
			t.Errorf("expected no allocs, got %v", actualAllocs)
		}
	})
	t.Run("ToGoHTML requires one allocation", func(t *testing.T) {
		expected := "<div>Unsanitized &</div>"
		c := tndr.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			_, err = io.WriteString(w, expected)
			return err
		})
		actualAllocs := testing.AllocsPerRun(4, func() {
			h, err := tndr.ToGoHTML(context.Background(), c)
			if err != nil {
				t.Fatalf("failed to convert to Go HTML: %v", err)
			}
			if h != template.HTML(expected) {
				t.Fatalf("unexpected value")
			}
		})
		if actualAllocs > 1 {
			t.Errorf("expected 1 alloc, got %v", actualAllocs)
		}
	})
}

func TestNonce(t *testing.T) {
	ctx := context.Background()
	t.Run("returns empty string if not set", func(t *testing.T) {
		actual := tndr.GetNonce(ctx)
		if actual != "" {
			t.Errorf("expected empty string got %q", actual)
		}
	})
	t.Run("returns value if one has been set", func(t *testing.T) {
		expected := "abc123"
		ctx := tndr.WithNonce(context.Background(), expected)
		actual := tndr.GetNonce(ctx)
		if actual != expected {
			t.Errorf("expected %q got %q", expected, actual)
		}
	})
}

func TestRenderAttributes(t *testing.T) {
	tests := []struct {
		name       string
		attributes tndr.Attributes
		expected   string
	}{
		{
			name: "string attributes are rendered",
			attributes: tndr.Attributes{
				"class": "test-class",
				"id":    "test-id",
			},
			expected: ` class="test-class" id="test-id"`,
		},
		{
			name: "integer types are rendered as strings",
			attributes: tndr.Attributes{
				"int":   42,
				"int8":  int8(8),
				"int16": int16(16),
				"int32": int32(32),
				"int64": int64(64),
			},
			expected: ` int="42" int16="16" int32="32" int64="64" int8="8"`,
		},
		{
			name: "unsigned integer types are rendered as strings",
			attributes: tndr.Attributes{
				"uint":    uint(42),
				"uint8":   uint8(8),
				"uint16":  uint16(16),
				"uint32":  uint32(32),
				"uint64":  uint64(64),
				"uintptr": uintptr(100),
			},
			expected: ` uint="42" uint16="16" uint32="32" uint64="64" uint8="8" uintptr="100"`,
		},
		{
			name: "float types are rendered as strings",
			attributes: tndr.Attributes{
				"float32": float32(3.14),
				"float64": float64(2.718),
			},
			expected: ` float32="3.14" float64="2.718"`,
		},
		{
			name: "complex types are rendered as strings",
			attributes: tndr.Attributes{
				"complex64":  complex64(1 + 2i),
				"complex128": complex128(3 + 4i),
			},
			expected: ` complex128="(3+4i)" complex64="(1+2i)"`,
		},
		{
			name: "boolean attributes are rendered correctly",
			attributes: tndr.Attributes{
				"checked":  true,
				"disabled": false,
			},
			expected: ` checked`,
		},
		{
			name: "mixed types are rendered correctly",
			attributes: tndr.Attributes{
				"class":  "button",
				"value":  42,
				"width":  float64(100.5),
				"hidden": false,
				"active": true,
			},
			expected: ` active class="button" value="42" width="100.5"`,
		},
		{
			name: "nil pointer attributes are not rendered",
			attributes: tndr.Attributes{
				"optional": (*string)(nil),
				"visible":  (*bool)(nil),
			},
			expected: ``,
		},
		{
			name: "non-nil pointer attributes are rendered",
			attributes: tndr.Attributes{
				"title":   ptr("test title"),
				"enabled": ptr(true),
			},
			expected: ` enabled title="test title"`,
		},
		{
			name: "numeric pointer types are rendered as strings",
			attributes: tndr.Attributes{
				"int-ptr":        ptr(42),
				"int8-ptr":       ptr(int8(8)),
				"int16-ptr":      ptr(int16(16)),
				"int32-ptr":      ptr(int32(32)),
				"int64-ptr":      ptr(int64(64)),
				"uint-ptr":       ptr(uint(42)),
				"uint8-ptr":      ptr(uint8(8)),
				"uint16-ptr":     ptr(uint16(16)),
				"uint32-ptr":     ptr(uint32(32)),
				"uint64-ptr":     ptr(uint64(64)),
				"uintptr-ptr":    ptr(uintptr(100)),
				"float32-ptr":    ptr(float32(3.14)),
				"float64-ptr":    ptr(float64(2.718)),
				"complex64-ptr":  ptr(complex64(1 + 2i)),
				"complex128-ptr": ptr(complex128(3 + 4i)),
			},
			expected: ` complex128-ptr="(3+4i)" complex64-ptr="(1+2i)" float32-ptr="3.14" float64-ptr="2.718" int-ptr="42" int16-ptr="16" int32-ptr="32" int64-ptr="64" int8-ptr="8" uint-ptr="42" uint16-ptr="16" uint32-ptr="32" uint64-ptr="64" uint8-ptr="8" uintptr-ptr="100"`,
		},
		{
			name: "nil numeric pointer attributes are not rendered",
			attributes: tndr.Attributes{
				"int-ptr":       (*int)(nil),
				"float32-ptr":   (*float32)(nil),
				"complex64-ptr": (*complex64)(nil),
			},
			expected: ``,
		},
		{
			name: "KeyValue[string, bool] attributes are rendered correctly",
			attributes: tndr.Attributes{
				"data-value":  tndr.KV("test-string", true),
				"data-hidden": tndr.KV("ignored", false),
			},
			expected: ` data-value="test-string"`,
		},
		{
			name: "KeyValue[bool, bool] attributes are rendered correctly",
			attributes: tndr.Attributes{
				"checked":  tndr.KV(true, true),
				"disabled": tndr.KV(false, true),
				"hidden":   tndr.KV(true, false),
			},
			expected: ` checked`,
		},
		{
			name: "function bool attributes are rendered correctly",
			attributes: tndr.Attributes{
				"enabled": func() bool { return true },
				"hidden":  func() bool { return false },
			},
			expected: ` enabled`,
		},
		{
			name: "mixed KeyValue and function attributes",
			attributes: tndr.Attributes{
				"data-name": tndr.KV("value", true),
				"active":    tndr.KV(true, true),
				"dynamic":   func() bool { return true },
				"ignored":   tndr.KV("ignored", false),
			},
			expected: ` active data-name="value" dynamic`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			err := tndr.RenderAttributes(context.Background(), &buf, tt.attributes)
			if err != nil {
				t.Fatalf("RenderAttributes failed: %v", err)
			}

			actual := buf.String()
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func ptr[T any](x T) *T {
	return &x
}
