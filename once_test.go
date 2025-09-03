package tndr_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/senforsce/tndr"
)

type onceHandleTest struct {
	ctx      context.Context
	expected string
}

func TestOnceHandle(t *testing.T) {
	withHello := tndr.WithChildren(context.Background(), tndr.Raw("hello"))
	tests := []struct {
		name  string
		tests []onceHandleTest
	}{
		{
			name: "renders nothing without children",
			tests: []onceHandleTest{
				{
					ctx:      context.Background(),
					expected: "",
				},
			},
		},
		{
			name: "children are rendered",
			tests: []onceHandleTest{
				{
					ctx:      tndr.WithChildren(context.Background(), tndr.Raw("hello")),
					expected: "hello",
				},
			},
		},
		{
			name: "children are rendered once per context",
			tests: []onceHandleTest{
				{
					ctx:      withHello,
					expected: "hello",
				},
				{
					ctx:      withHello,
					expected: "",
				},
			},
		},
		{
			name: "different contexts have different once state",
			tests: []onceHandleTest{
				{
					ctx:      tndr.WithChildren(context.Background(), tndr.Raw("hello")),
					expected: "hello",
				},
				{
					ctx:      tndr.WithChildren(context.Background(), tndr.Raw("hello2")),
					expected: "hello2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tndr.NewOnceHandle().Once()
			for i, test := range tt.tests {
				t.Run(fmt.Sprintf("render %d/%d", i+1, len(tt.tests)), func(t *testing.T) {
					html, err := tndr.ToGoHTML(test.ctx, c)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}
					if diff := cmp.Diff(test.expected, string(html)); diff != "" {
						t.Errorf("unexpected diff:\n%v", diff)
					}
				})
			}
		})
	}
	t.Run("each new handle manages different state", func(t *testing.T) {
		ctx := tndr.WithChildren(context.Background(), tndr.Raw("hello"))
		h1 := tndr.NewOnceHandle()
		c1 := h1.Once()
		h2 := tndr.NewOnceHandle()
		c2 := h2.Once()
		c3 := h2.Once()
		var w strings.Builder
		if err := c1.Render(ctx, &w); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := c2.Render(ctx, &w); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := c3.Render(ctx, &w); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if diff := cmp.Diff("hellohello", w.String()); diff != "" {
			t.Errorf("unexpected diff:\n%v", diff)
		}
	})
	t.Run("a handle can be used to render a specific component", func(t *testing.T) {
		ctx := tndr.WithChildren(context.Background(), tndr.Raw("child"))
		o := tndr.NewOnceHandle(tndr.WithComponent(tndr.Raw("c"))).Once()
		var w strings.Builder
		if err := o.Render(ctx, &w); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if diff := cmp.Diff("c", w.String()); diff != "" {
			t.Errorf("unexpected diff:\n%v", diff)
		}
	})
}
