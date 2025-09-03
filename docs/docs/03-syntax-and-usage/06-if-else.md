# If/else

tndr uses standard Go `if`/`else` statements which can be used to conditionally render components and elements.

```tndr title="component.t1"
t1 login(isLoggedIn bool) {
  if isLoggedIn {
    <div>Welcome back!</div>
  } else {
    <input name="login" type="button" value="Log in"/>
  }
}
```

```go title="main.go"
package main

import (
	"context"
	"os"
)

func main() {
	login(true).Render(context.Background(), os.Stdout)
}
```

```html title="Output"
<div>
 Welcome back!
</div>
```
