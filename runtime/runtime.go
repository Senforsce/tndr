package runtime

import (
	"context"
	"io"

	"github.com/senforsce/tndr"
)

// GeneratedComponentInput is used to avoid generated code needing to import the `context` and `io` packages.
type GeneratedComponentInput struct {
	Context context.Context
	Writer  io.Writer
}

// GeneratedTemplate is used to avoid generated code needing to import the `context` and `io` packages.
func GeneratedTemplate(f func(GeneratedComponentInput) error) tndr.Component {
	return tndr.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return f(GeneratedComponentInput{ctx, w})
	})
}
