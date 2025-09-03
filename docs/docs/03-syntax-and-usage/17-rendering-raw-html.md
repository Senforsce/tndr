# Rendering raw HTML

To render HTML that has come from a trusted source, bypassing all HTML escaping and security mechanisms that tndr includes, use the `tndr.Raw` function.

:::info
Only include HTML that comes from a trusted source.
:::

:::warning
Use of this function may introduce security vulnerabilities to your program.
:::

```tndr title="component.t1"
t1 Example() {
	<!DOCTYPE html>
	<html>
		<body>
			@tndr.Raw("<div>Hello, World!</div>")
		</body>
	</html>
}
```

```html title="Output"
<!DOCTYPE html>
<html>
	<body>
		<div>Hello, World!</div>
	</body>
</html>
```
