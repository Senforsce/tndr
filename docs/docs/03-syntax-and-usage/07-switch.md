# Switch

tndr uses standard Go `switch` statements which can be used to conditionally render components and elements.

```tndr title="component.t1"
package main

t1 userTypeDisplay(userType string) {
	switch userType {
		case "test":
			<span>{ "Test user" }</span>
		case "admin":
			<span>{ "Admin user" }</span>
		default:
			<span>{ "Unknown user" }</span>
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
	userTypeDisplay("Other").Render(context.Background(), os.Stdout)
}
```

```html title="Output"
<span>
 Unknown user
</span>
```
