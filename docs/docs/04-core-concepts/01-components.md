# Components

tndr Components are markup and code that is compiled into functions that return a `tndr.Component` interface by running the `tndr generate` command.

Components can contain tndr elements that render HTML, text, expressions that output text or include other templates, and branching statements such as `if` and `switch`, and `for` loops.

```tndr title="header.t1"
package main

t1 headerTemplate(name string) {
  <header data-testid="headerTemplate">
    <h1>{ name }</h1>
  </header>
}
```

The generated code is a Go function that returns a `tndr.Component`.

```go title="header_tndr.go"
func headerTemplate(name string) tndr.Component {
  // Generated contents
}
```

`tndr.Component` is an interface that has a `Render` method on it that is used to render the component to an `io.Writer`.

```go
type Component interface {
	Render(ctx context.Context, w io.Writer) error
}
```

:::tip
Since tndr produces Go code, you can share templates the same way that you share Go code - by sharing your Go module.

tndr follows the same rules as Go. If a `tndr` block starts with an uppercase letter, then it is public, otherwise, it is private.

A `tndr.Component` may write partial output to the `io.Writer` if it returns an error. If you want to ensure you only get complete output or nothing, write to a buffer first and then write the buffer to an `io.Writer`.
:::

## Code-only components

Since tndr Components ultimately implement the `tndr.Component` interface, any code that implements the interface can be used in place of a tndr component generated from a `*.t1` file.

```go
package main

import (
	"context"
	"io"
	"os"

	"github.com/senforsce/tndr"
)

func button(text string) tndr.Component {
	return tndr.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, "<button>"+text+"</button>")
		return err
	})
}

func main() {
	button("Click me").Render(context.Background(), os.Stdout)
}
```

```html title="Output"
<button>
 Click me
</button>
```

:::warning
This code is unsafe! In code-only components, you're responsible for escaping the HTML content yourself, e.g. with the `tndr.EscapeString` function.
:::

## Method components

tndr components can be returned from methods (functions attached to types).

Go code:

```tndr
package main

import "os"

type Data struct {
	message string
}

t1 (d Data) Method() {
	<div>{ d.message }</div>
}

func main() {
	d := Data{
		message: "You can implement methods on a type.",
	}
	d.Method().Render(context.Background(), os.Stdout)
}
```

It is also possible to initialize a struct and call its component method inline.

```tndr
package main

import "os"

type Data struct {
	message string
}

t1 (d Data) Method() {
	<div>{ d.message }</div>
}

t1 Message() {
    <div>
        @Data{
            message: "You can implement methods on a type.",
        }.Method()
    </div>
}

func main() {
	Message().Render(context.Background(), os.Stdout)
}
```

