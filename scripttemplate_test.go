package tndr_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/senforsce/tndr"
)

func TestRenderScriptItems(t *testing.T) {
	s1 := tndr.ComponentScript{
		Name:     "s1",
		Function: "function s1() { return 'hello1'; }",
	}
	s2 := tndr.ComponentScript{
		Name:     "s2",
		Function: "function s2() { return 'hello2'; }",
	}
	tests := []struct {
		name     string
		toIgnore []tndr.ComponentScript
		toRender []tndr.ComponentScript
		expected string
	}{
		{
			name:     "if none are ignored, everything is rendered",
			toIgnore: nil,
			toRender: []tndr.ComponentScript{s1, s2},
			expected: `<script>` + s1.Function + s2.Function + `</script>`,
		},
		{
			name: "if something outside the expected is ignored, if has no effect",
			toIgnore: []tndr.ComponentScript{
				{
					Name:     "s3",
					Function: "function s3() { return 'hello3'; }",
				},
			},
			toRender: []tndr.ComponentScript{s1, s2},
			expected: `<script>` + s1.Function + s2.Function + `</script>`,
		},
		{
			name:     "if one is ignored, it's not rendered",
			toIgnore: []tndr.ComponentScript{s1},
			toRender: []tndr.ComponentScript{s1, s2},
			expected: `<script>` + s2.Function + `</script>`,
		},
		{
			name: "if all are ignored, not even style tags are rendered",
			toIgnore: []tndr.ComponentScript{
				s1,
				s2,
				{
					Name:     "s3",
					Function: "function s3() { return 'hello3'; }",
				},
			},
			toRender: []tndr.ComponentScript{s1, s2},
			expected: ``,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			b := new(bytes.Buffer)

			// Render twice, reusing the same context so that there's a memory of which classes have been rendered.
			ctx = tndr.InitializeContext(ctx)
			err := tndr.RenderScriptItems(ctx, b, tt.toIgnore...)
			if err != nil {
				t.Fatalf("failed to render initial scripts: %v", err)
			}

			// Now render again to check that only the expected classes were rendered.
			b.Reset()
			err = tndr.RenderScriptItems(ctx, b, tt.toRender...)
			if err != nil {
				t.Fatalf("failed to render scripts: %v", err)
			}

			if diff := cmp.Diff(tt.expected, b.String()); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestJSExpression(t *testing.T) {
	expected := "myJSFunction(\"StringValue\",123,event,1 + 2)"
	actual := tndr.SafeScriptInline("myJSFunction", "StringValue", 123, tndr.JSExpression("event"), tndr.JSExpression("1 + 2"))

	if actual != expected {
		t.Fatalf("TestJSExpression: expected %q, got %q", expected, actual)
	}
}
