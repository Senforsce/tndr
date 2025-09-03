# Creating a simple tndr component

To create a tndr component, first create a new Go project.

## Setup project

Create a new directory containing our project.

```bash
mkdir hello-world
```

Initialize a new Go project within it.

```bash
cd hello-world
go mod init github.com/senforsce/tndr-examples/hello-world
go get github.com/senforsce/tndr
```

## Create a tndr file

To use it, create a `hello.t1` file containing a component.

Components are functions that contain tndr elements, markup, and `if`, `switch`, and `for` Go expressions.

```tndr title="hello.t1"
package main

t1 hello(name string) {
	<div>Hello, { name }</div>
}
```

## Generate Go code from the tndr file

Run the `t1 generate` command.

```bash
t1 generate
```

tndr will generate a `hello_t1.go` file containing Go code.

This file will contain a function called `hello` which takes `name` as an argument, and returns a `tndr.Component` that renders HTML.

```go
func hello(name string) tndr.Component {
  // ...
}
```

## Write a program that renders to stdout

Create a `main.go` file.

```go title="main.go"
package main

import (
	"context"
	"os"
)

func main() {
	component := hello("John")
	component.Render(context.Background(), os.Stdout)
}
```

## Run the program

Running the code will render the component's HTML to stdout.

```bash
go run .
```

```html title="Output"
<div>Hello, John</div>
```

Instead of passing `os.Stdout` to the component's render function, you can pass any type that implements the `io.Writer` interface. This includes files, `bytes.Buffer`, and HTTP responses.

In this way, tndr can be used to generate HTML files that can be hosted as static content in an S3 bucket, Google Cloud Storage, or used to generate HTML that is fed into PDF conversion processes, or sent via email.
