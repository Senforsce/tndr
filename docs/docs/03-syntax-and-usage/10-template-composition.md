# Template composition

Templates can be composed using the import expression.

```tndr
t1 showAll() {
	@left()
	@middle()
	@right()
}

t1 left() {
	<div>Left</div>
}

t1 middle() {
	<div>Middle</div>
}

t1 right() {
	<div>Right</div>
}
```

```html title="Output"
<div>
 Left
</div>
<div>
 Middle
</div>
<div>
 Right
</div>
```

## Children

Children can be passed to a component for it to wrap.

```tndr
t1 showAll() {
	@wrapChildren() {
		<div>Inserted from the top</div>
	}
}

t1 wrapChildren() {
	<div id="wrapper">
		{ children... }
	</div>
}
```

:::note
The use of the `{ children... }` expression in the child component.
:::

```html title="output"
<div id="wrapper">
 <div>
  Inserted from the top
 </div>
</div>
```

### Using children in code components

Children are passed to a component using the Go context. To pass children to a component using Go code, use the `tndr.WithChildren` function.

```tndr
package main

import (
  "context"
  "os"

  "github.com/senforsce/tndr"
)

t1 wrapChildren() {
	<div id="wrapper">
		{ children... }
	</div>
}

func main() {
  contents := tndr.ComponentFunc(func(ctx context.Context, w io.Writer) error {
    _, err := io.WriteString(w, "<div>Inserted from Go code</div>")
    return err
  })
  ctx := tndr.WithChildren(context.Background(), contents)
  wrapChildren().Render(ctx, os.Stdout)
}
```

```html title="output"
<div id="wrapper">
 <div>
  Inserted from Go code
 </div>
</div>
```

To get children from the context, use the `tndr.GetChildren` function.

```tndr
package main

import (
  "context"
  "os"

  "github.com/senforsce/tndr"
)

func main() {
  contents := tndr.ComponentFunc(func(ctx context.Context, w io.Writer) error {
    _, err := io.WriteString(w, "<div>Inserted from Go code</div>")
    return err
  })
  wrapChildren := tndr.ComponentFunc(func(ctx context.Context, w io.Writer) error {
    children := tndr.GetChildren(ctx)
    ctx = tndr.ClearChildren(ctx)
    _, err := io.WriteString(w, "<div id=\"wrapper\">")
    if err != nil {
      return err
    }
    err = children.Render(ctx, w)
    if err != nil {
      return err
    }
    _, err = io.WriteString(w, "</div>")
    return err
  })
```

:::note
The `tndr.ClearChildren` function is used to stop passing the children down the tree.
:::

## Components as parameters

Components can also be passed as parameters and rendered using the `@component` expression.

```tndr
package main

t1 heading() {
    <h1>Heading</h1>
}

t1 layout(contents tndr.Component) {
	<div id="heading">
		@heading()
	</div>
	<div id="contents">
		@contents
	</div>
}

t1 paragraph(contents string) {
	<p>{ contents }</p>
}
```

```go title="main.go"
package main

import (
	"context"
	"os"
)

func main() {
	c := paragraph("Dynamic contents")
	layout(c).Render(context.Background(), os.Stdout)
}
```

```html title="output"
<div id="heading">
	<h1>Heading</h1>
</div>
<div id="contents">
	<p>Dynamic contents</p>
</div>
```

You can pass `tndr` components as parameters to other components within tndrates using standard Go function call syntax.

```tndr
package main

t1 heading() {
    <h1>Heading</h1>
}

t1 layout(contents tndr.Component) {
	<div id="heading">
		@heading()
	</div>
	<div id="contents">
		@contents
	</div>
}

t1 paragraph(contents string) {
	<p>{ contents }</p>
}

t1 root() {
	@layout(paragraph("Dynamic contents"))
}
```

```go title="main.go"
package main

import (
	"context"
	"os"
)

func main() {
	root().Render(context.Background(), os.Stdout)
}
```

```html title="output"
<div id="heading">
	<h1>Heading</h1>
</div>
<div id="contents">
	<p>Dynamic contents</p>
</div>
```

## Joining Components

Components can be aggregated into a single Component using `tndr.Join`.

```tndr
package main

t1 hello() {
	<span>hello</span>
}

t1 world() {
	<span>world</span>
}

t1 helloWorld() {
	@tndr.Join(hello(), world())
}
```

```go title="main.go"
package main

import (
	"context"
	"os"
)

func main() {
	helloWorld().Render(context.Background(), os.Stdout)
}
```

```html title="output"
<span>hello</span><span>world</span>
```

## Sharing and re-using components

Since tndr components are compiled into Go functions by the `go generate` command, tndr components follow the rules of Go, and are shared in exactly the same way as Go code.

t1 files in the same directory can access each other's components. Components in different directories can be accessed by importing the package that contains the component, so long as the component is exported by capitalizing its name.

:::tip
In Go, a _package_ is a collection of Go source files in the same directory that are compiled together. All of the functions, types, variables, and constants defined in one source file in a package are available to all other source files in the same package.

Packages exist within a Go _module_, defined by the `go.mod` file.
:::

:::note
Go is structured differently to JavaScript, but uses similar terminology. A single `.js` or `.ts` _file_ is like a Go package, and an NPM package is like a Go module.
:::

### Exporting components

To make a tndr component available to other packages, export it by capitalizing its name.

```tndr
package components

t1 Hello() {
	<div>Hello</div>
}
```

### Importing components

To use a component in another package, import the package and use the component as you would any other Go function or type.

```tndr
package main

import "github.com/senforsce/tndr/examples/counter/components"

t1 Home() {
	@components.Hello()
}
```

:::tip
To import a component from another Go module, you must first import the module by using the `go get <module>` command. Then, you can import the component as you would any other Go package.
:::
