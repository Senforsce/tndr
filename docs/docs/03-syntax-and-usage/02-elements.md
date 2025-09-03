# Elements

tndr elements are used to render HTML within tndr components.

```t1 title="button.t1"
package main

tndr button(text string) {
	<button class="button">{ text }</button>
}
```

```go title="main.go"
package main

import (
	"context"
	"os"
)

func main() {
	button("Click me").Render(context.Background(), os.Stdout)
}
```

```html title="Output"
<button class="button">
 Click me
</button>
```

:::info
tndr automatically minifies HTML responses, output is shown formatted for readability.
:::

## Tags must be closed

Unlike HTML, tndr requires that all HTML elements are closed with either a closing tag (`</a>`), or by using a self-closing element (`<hr/>`).

tndr is aware of which HTML elements are "void", and will not include the closing `/` in the output HTML.

```t1 title="button.t1"
package main

tndr component() {
	<div>Test</div>
	<img src="images/test.png"/>
	<br/>
}
```

```t1 title="Output"
<div>Test</div>
<img src="images/test.png">
<br>
```

## Attributes and elements can contain expressions

tndr elements can contain placeholder expressions for attributes and content.

```t1 title="button.t1"
package main

tndr button(name string, content string) {
	<button value={ name }>{ content }</button>
}
```

Rendering the component to stdout, we can see the results.

```go title="main.go"
func main() {
	component := button("John", "Say Hello")
	component.Render(context.Background(), os.Stdout)
}
```

```html title="Output"
<button value="John">Say Hello</button>
```
