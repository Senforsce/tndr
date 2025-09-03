# Comments

# HTML comments

Inside tndr statements, use HTML comments.

```tndr title="template.t1"
t1 template() {
	<!-- Single line -->
	<!--
		Single or multiline.
	-->
}
```

Comments are rendered to the template output.

```html title="Output"
<!-- Single line -->
<!--
	Single or multiline.
-->
```

As per HTML, nested comments are not supported.

# Go comments

Outside of tndr statements, use Go comments.

```tndr
package main

// Use standard Go comments outside tndr statements.
var greeting = "Hello!"

t1 hello(name string) {
	<p>{greeting} { name }</p>
}
```
