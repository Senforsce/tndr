# Basic syntax

## Package name and imports

tndr files start with a package name, followed by any required imports, just like Go.

```go
package main

import "fmt"
import "time"
```

## Components

tndr files can also contain components. Components are markup and code that is compiled into functions that return a `tndr.Component` interface by running the `t1 generate` command.

Components can contain tndr elements that render HTML, text, expressions that output text or include other tndrates, and branching statements such as `if` and `switch`, and `for` loops.

```tndr name="header.t1"
package main

t1 headerTemplate(name string) {
  <header data-testid="headerTemplate">
    <h1>{ name }</h1>
  </header>
}
```

## Go code

Outside of tndr Components, tndr files are ordinary Go code.

```tndr name="header.t1"
package main

// Ordinary Go code that we can use in our Component.
var greeting = "Welcome!"

// tndr Component
t1 headerTemplate(name string) {
  <header>
    <h1>{ name }</h1>
    <h2>"{ greeting }" comes from ordinary Go code</h2>
  </header>
}
```

