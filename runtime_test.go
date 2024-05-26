package t1_test

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
	"github.com/senforsce/t1"
)

func TestCSSHandler(t *testing.T) {
	tests := []struct {
		name             string
		input            []t1.CSSClass
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
			input:            []t1.CSSClass{t1.ComponentCSSClass{ID: "className", Class: t1.SafeCSS(".className{background-color:white;}")}},
			expectedMIMEType: "text/css",
			expectedBody:     ".className{background-color:white;}",
		},
		{
			name: "classes are rendered",
			input: []t1.CSSClass{
				t1.ComponentCSSClass{ID: "classA", Class: t1.SafeCSS(".classA{background-color:white;}")},
				t1.ComponentCSSClass{ID: "classB", Class: t1.SafeCSS(".classB{background-color:green;}")},
			},
			expectedMIMEType: "text/css",
			expectedBody:     ".classA{background-color:white;}.classB{background-color:green;}",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			h := t1.NewCSSHandler(tt.input...)
			h.ServeHTTP(w, &http.Request{})
			if diff := cmp.Diff(tt.expectedMIMEType, w.Header().Get("Content-Type")); diff != "" {
				t.Errorf(diff)
			}
			if diff := cmp.Diff(tt.expectedBody, w.Body.String()); diff != "" {
				t.Errorf(diff)
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
	c1 := t1.ComponentCSSClass{
		ID:    "c1",
		Class: ".c1{color:red}",
	}
	c2 := t1.ComponentCSSClass{
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
			name:             "accessing /style/t1.css renders CSS, even if it's empty",
			input:            httptest.NewRequest("GET", "/styles/t1.css", nil),
			handler:          t1.NewCSSMiddleware(pageHandler),
			expectedMIMEType: "text/css",
			expectedBody:     "",
		},
		{
			name:             "accessing /style/t1.css renders CSS that includes the classes",
			input:            httptest.NewRequest("GET", "/styles/t1.css", nil),
			handler:          t1.NewCSSMiddleware(pageHandler, c1, c2),
			expectedMIMEType: "text/css",
			expectedBody:     ".c1{color:red}.c2{color:blue}",
		},
		{
			name:             "the pageHandler is rendered",
			input:            httptest.NewRequest("GET", "/index.html", nil),
			handler:          t1.NewCSSMiddleware(pageHandler, c1, c2),
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
				t.Errorf(diff)
			}
			if diff := cmp.Diff(tt.expectedBody, w.Body.String()); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

var cssInputs = []any{
	[]string{"a", "b"},       // []string
	"c",                      // string
	t1.ConstantCSSClass("d"), // ConstantCSSClass
	t1.ComponentCSSClass{ID: "e", Class: ".e{color:red}"}, // ComponentCSSClass
	map[string]bool{"f": true, "ff": false},               // map[string]bool
	t1.KV[string, bool]("g", true),                        // KeyValue[string, bool]
	t1.KV[string, bool]("gg", false),                      // KeyValue[string, bool]
	[]t1.KeyValue[string, bool]{
		t1.KV("h", true),
		t1.KV("hh", false),
	}, // []KeyValue[string, bool]
	t1.KV[t1.CSSClass, bool](t1.ConstantCSSClass("i"), true),   // KeyValue[CSSClass, bool]
	t1.KV[t1.CSSClass, bool](t1.ConstantCSSClass("ii"), false), // KeyValue[CSSClass, bool]
	t1.KV[t1.ComponentCSSClass, bool](t1.ComponentCSSClass{
		ID:    "j",
		Class: ".j{color:red}",
	}, true), // KeyValue[ComponentCSSClass, bool]
	t1.KV[t1.ComponentCSSClass, bool](t1.ComponentCSSClass{
		ID:    "jj",
		Class: ".jj{color:red}",
	}, false), // KeyValue[ComponentCSSClass, bool]
	t1.CSSClasses{t1.ConstantCSSClass("k")},                          // CSSClasses
	func() t1.CSSClass { return t1.ConstantCSSClass("l") },           // func() CSSClass
	t1.CSSClass(t1.ConstantCSSClass("m")),                            // CSSClass
	customClass{name: "n"},                                           // CSSClass
	t1.KV[t1.ConstantCSSClass, bool](t1.ConstantCSSClass("o"), true), // KeyValue[ConstantCSSClass, bool]
	[]t1.KeyValue[t1.ConstantCSSClass, bool]{
		t1.KV(t1.ConstantCSSClass("p"), true),
		t1.KV(t1.ConstantCSSClass("pp"), false),
	}, // []KeyValue[ConstantCSSClass, bool]
}

type customClass struct {
	name string
}

func (cc customClass) ClassName() string {
	return cc.name
}

func TestRenderCSS(t *testing.T) {
	c1 := t1.ComponentCSSClass{
		ID:    "c1",
		Class: ".c1{color:red}",
	}
	c2 := t1.ComponentCSSClass{
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
				t1.ComponentCSSClass{
					ID:    "c3",
					Class: t1.SafeCSS(".c3{color:yellow}"),
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
				t1.ComponentCSSClass{
					ID:    "c3",
					Class: t1.SafeCSS(".c3{color:yellow}"),
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
			ctx = t1.InitializeContext(ctx)
			err := t1.RenderCSSItems(ctx, b, tt.toIgnore...)
			if err != nil {
				t.Fatalf("failed to render initial CSS: %v", err)
			}

			// Now render again to check that only the expected classes were rendered.
			b.Reset()
			err = t1.RenderCSSItems(ctx, b, tt.toRender...)
			if err != nil {
				t.Fatalf("failed to render CSS: %v", err)
			}

			if diff := cmp.Diff(tt.expected, b.String()); diff != "" {
				t.Error(diff)
			}
		})
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
			input:    []any{"a", t1.SafeClass("b"), t1.Class("c")},
			expected: "a b c",
		},
		{
			name: "CSS components are included in the output",
			input: []any{
				t1.ComponentCSSClass{ID: "classA", Class: t1.SafeCSS(".classA{background-color:white;}")},
				t1.ComponentCSSClass{ID: "classB", Class: t1.SafeCSS(".classB{background-color:green;}")},
				"c",
			},
			expected: "classA classB c",
		},
		{
			name: "optional classes can be applied with expressions",
			input: []any{
				"a",
				t1.ComponentCSSClass{ID: "classA", Class: t1.SafeCSS(".classA{background-color:white;}")},
				t1.ComponentCSSClass{ID: "classB", Class: t1.SafeCSS(".classB{background-color:green;}")},
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
			name: "unknown types for classes get rendered as --t1-css-class-unknown-type",
			input: []any{
				123,
				map[string]string{"test": "no"},
				false,
				"c",
			},
			expected: "--t1-css-class-unknown-type c",
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
			name: "if a t1.CSSClasses is passed in, the nested CSSClasses are extracted",
			input: []any{
				t1.Classes(
					"a",
					t1.SafeClass("b"),
					t1.Class("c"),
					t1.ComponentCSSClass{
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
				t1.KV("b", true),
				"c",
				t1.KV("c", false),
				t1.KV(t1.SafeClass("d"), true),
				t1.KV(t1.SafeClass("e"), false),
			},
			expected: "a b d",
		},
		{
			name: "an array of KV types can be used to show or hide classes",
			input: []any{
				"a",
				"c",
				[]t1.KeyValue[string, bool]{
					t1.KV("b", true),
					t1.KV("c", false),
					{"d", true},
				},
			},
			expected: "a b d",
		},
		{
			name: "the brackets on component CSS function calls can be elided",
			input: []any{
				func() t1.CSSClass {
					return t1.ComponentCSSClass{
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
			actual := t1.Classes(test.input...).String()
			if actual != test.expected {
				t.Errorf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}

func TestHandler(t *testing.T) {
	hello := t1.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, "Hello"); err != nil {
			t.Fatalf("failed to write string: %v", err)
		}
		return nil
	})
	errorComponent := t1.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, "Hello"); err != nil {
			t.Fatalf("failed to write string: %v", err)
		}
		return errors.New("handler error")
	})

	tests := []struct {
		name             string
		input            *t1.ComponentHandler
		expectedStatus   int
		expectedMIMEType string
		expectedBody     string
	}{
		{
			name:             "handlers return OK by default",
			input:            t1.Handler(hello),
			expectedStatus:   http.StatusOK,
			expectedMIMEType: "text/html; charset=utf-8",
			expectedBody:     "Hello",
		},
		{
			name:             "handlers return OK by default",
			input:            t1.Handler(t1.Raw(`♠ ‘ &spades; &#8216;`)),
			expectedStatus:   http.StatusOK,
			expectedMIMEType: "text/html; charset=utf-8",
			expectedBody:     "♠ ‘ &spades; &#8216;",
		},
		{
			name:             "handlers can be configured to return an alternative status code",
			input:            t1.Handler(hello, t1.WithStatus(http.StatusNotFound)),
			expectedStatus:   http.StatusNotFound,
			expectedMIMEType: "text/html; charset=utf-8",
			expectedBody:     "Hello",
		},
		{
			name:             "handlers can be configured to return an alternative status code and content type",
			input:            t1.Handler(hello, t1.WithStatus(http.StatusOK), t1.WithContentType("text/csv")),
			expectedStatus:   http.StatusOK,
			expectedMIMEType: "text/csv",
			expectedBody:     "Hello",
		},
		{
			name:             "handlers that fail return a 500 error",
			input:            t1.Handler(errorComponent),
			expectedStatus:   http.StatusInternalServerError,
			expectedMIMEType: "text/plain; charset=utf-8",
			expectedBody:     "t1: failed to render template\n",
		},
		{
			name: "error handling can be customised",
			input: t1.Handler(errorComponent, t1.WithErrorHandler(func(r *http.Request, err error) http.Handler {
				// Because the error is received, it's possible to log the detail of the request.
				// log.Printf("template render error for %v %v: %v", r.Method, r.URL.String(), err)
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
					if _, err := io.WriteString(w, "custom body"); err != nil {
						t.Fatalf("failed to write string: %v", err)
					}
				})
			})),
			expectedStatus:   http.StatusBadRequest,
			expectedMIMEType: "text/html; charset=utf-8",
			expectedBody:     "custom body",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/test", nil)
			tt.input.ServeHTTP(w, r)
			if got := w.Result().StatusCode; tt.expectedStatus != got {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, got)
			}
			if mimeType := w.Result().Header.Get("Content-Type"); tt.expectedMIMEType != mimeType {
				t.Errorf("expected content-type %s, got %s", tt.expectedMIMEType, mimeType)
			}
			body, err := io.ReadAll(w.Result().Body)
			if err != nil {
				t.Errorf("failed to read body: %v", err)
			}
			if diff := cmp.Diff(tt.expectedBody, string(body)); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestRenderScriptItems(t *testing.T) {
	s1 := t1.ComponentScript{
		Name:     "s1",
		Function: "function s1() { return 'hello1'; }",
	}
	s2 := t1.ComponentScript{
		Name:     "s2",
		Function: "function s2() { return 'hello2'; }",
	}
	tests := []struct {
		name     string
		toIgnore []t1.ComponentScript
		toRender []t1.ComponentScript
		expected string
	}{
		{
			name:     "if none are ignored, everything is rendered",
			toIgnore: nil,
			toRender: []t1.ComponentScript{s1, s2},
			expected: `<script type="text/javascript">` + s1.Function + s2.Function + `</script>`,
		},
		{
			name: "if something outside the expected is ignored, if has no effect",
			toIgnore: []t1.ComponentScript{
				{
					Name:     "s3",
					Function: "function s3() { return 'hello3'; }",
				},
			},
			toRender: []t1.ComponentScript{s1, s2},
			expected: `<script type="text/javascript">` + s1.Function + s2.Function + `</script>`,
		},
		{
			name:     "if one is ignored, it's not rendered",
			toIgnore: []t1.ComponentScript{s1},
			toRender: []t1.ComponentScript{s1, s2},
			expected: `<script type="text/javascript">` + s2.Function + `</script>`,
		},
		{
			name: "if all are ignored, not even style tags are rendered",
			toIgnore: []t1.ComponentScript{
				s1,
				s2,
				{
					Name:     "s3",
					Function: "function s3() { return 'hello3'; }",
				},
			},
			toRender: []t1.ComponentScript{s1, s2},
			expected: ``,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			b := new(bytes.Buffer)

			// Render twice, reusing the same context so that there's a memory of which classes have been rendered.
			ctx = t1.InitializeContext(ctx)
			err := t1.RenderScriptItems(ctx, b, tt.toIgnore...)
			if err != nil {
				t.Fatalf("failed to render initial scripts: %v", err)
			}

			// Now render again to check that only the expected classes were rendered.
			b.Reset()
			err = t1.RenderScriptItems(ctx, b, tt.toRender...)
			if err != nil {
				t.Fatalf("failed to render scripts: %v", err)
			}

			if diff := cmp.Diff(tt.expected, b.String()); diff != "" {
				t.Error(diff)
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
	wrappedErr := t1.Error{Err: baseErr, Line: 1, Col: 2}
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
		input       t1.Component
		expected    string
		expectedErr error
	}{
		{
			name:     "Raw content is not escaped",
			input:    t1.Raw("<div>Test &</div>"),
			expected: `<div>Test &</div>`,
		},
		{
			name:        "Raw will return errors first",
			input:       t1.Raw("", nil, errors.New("test error")),
			expected:    `<div>Test &</div>`,
			expectedErr: errors.New("test error"),
		},
		{
			name:     "Strings marked as safe are rendered without escaping",
			input:    t1.Raw(template.HTML("<div>")),
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
			c := t1.Raw("<div>")
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
	t.Run("Go templates can be rendered as t1 components", func(t *testing.T) {
		b := new(bytes.Buffer)
		err := t1.FromGoHTML(goTemplate, "Test &").Render(context.Background(), b)
		if err != nil {
			t.Fatalf("failed to render content: %v", err)
		}
		if diff := cmp.Diff("<div>Test &amp;</div>", b.String()); diff != "" {
			t.Error(diff)
		}
	})
	t.Run("t1 components can be rendered in Go templates", func(t *testing.T) {
		b := new(bytes.Buffer)
		c := t1.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			_, err = io.WriteString(w, "<div>Unsanitized &</div>")
			return err
		})
		h, err := t1.ToGoHTML(context.Background(), c)
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
		c := t1.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			return expectedErr
		})
		_, err := t1.ToGoHTML(context.Background(), c)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if err != expectedErr {
			t.Fatalf("expected error %q, got %q", expectedErr, err)
		}
	})
	t.Run("FromGoHTML does not require allocations", func(t *testing.T) {
		actualAllocs := testing.AllocsPerRun(4, func() {
			c := t1.FromGoHTML(goTemplate, "test &")
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
		c := t1.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			_, err = io.WriteString(w, expected)
			return err
		})
		actualAllocs := testing.AllocsPerRun(4, func() {
			h, err := t1.ToGoHTML(context.Background(), c)
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
