# Expressions

## Interpolation expressions

Within a t1 element, expressions can be used to interpolate Go values. Content is automatically escaped using context-aware HTML encoding rules to protect against XSS and CSS injection attacks.

Literals, variables and functions that return a value can be used. 

The supported types for interpolation are:

- Strings
- Numbers (`int`, `uint`, `float32`, `complex64` etc.)
- Booleans

:::note
Any type based on the above list can be used, for example `type Age int` or `type Name string`.
:::

### Literals

You can use Go literals.

```t1 title="component.t1"
package main

t1 component() {
  <div>{ "print this" }</div>
  <div>{ `and this` }</div>
  <div>Number of the day: { 1 }</div>
}
```

```html title="Output"
<div>print this</div><div>and this</div><div>Number of the day: 1</div>
```

### Variables

Any supported Go variable can be used, for example:

* A function parameter.
* A field on a struct.
* A variable or constant that is in scope.

```t1 title="/main.t1"
package main

t1 greet(prefix string, p Person) {
  <div>{ prefix } { p.Name }{ exclamation }</div>
  <div>Congratulations on being { p.Age }!</div>
}
```

```t1 title="main.go"
package main

type Person struct {
  Name string
  Age  int
}

const exclamation = "!"

func main() {
  p := Person{ Name: "John", Age: 42 }
  component := greet("Hello", p) 
  component.Render(context.Background(), os.Stdout)
}
```

```html title="Output"
<div>Hello John!</div><div>Congratulations on being 42!</div>
```

### Functions

Functions that return a value, or a value-error tuple can be used.

```t1 title="component.t1"
package main

import "strings"
import "strconv"

func getString() (string, error) {
  return "DEF", nil
}

t1 component() {
  <div>{ strings.ToUpper("abc") }</div>
  <div>{ getString() }</div>
}
```

```html title="Output"
<div>ABC</div>
<div>DEF</div>
```

If the function returns an error, the `Render` function will return an error containing the location of the error and the underlying error.

### Escaping

t1 automatically escapes strings using HTML escaping rules.

```t1 title="component.t1"
package main

t1 component() {
  <div>{ `</div><script>alert('hello!')</script><div>` }</div>
}
```

```html title="Output"
<div>&lt;/div&gt;&lt;script&gt;alert(&#39;hello!&#39;)&lt;/script&gt;&lt;div&gt;</div>
```
