# Injection attacks

templ is designed to prevent user-provided data from being used to inject vulnerabilities.

`<script>` and `<style>` tags could allow user data to inject vulnerabilities, so variables are not permitted in these sections.

```html
templ Example() {
  <script>
    function showAlert() {
      alert("hello");
    }
  </script>
  <style type="text/css">
    /* Only CSS is allowed */
  </style>
}
```

`onClick` attributes, and other `on*` attributes are used to execute JavaScript. To prevent user data from being unescaped, `on*` attributes accept a `tndr.ComponentScript`.

```html
script onClickHandler(msg string) {
  alert(msg);
}

templ Example(msg string) {
  <div onClick={ onClickHandler(msg) }>
    { "will be HTML encoded using tndr.Escape" }
  </div>
}
```

Style attributes cannot be expressions, only constants, to avoid escaping vulnerabilities. templ style templates (`css className()`) should be used instead.

```html
templ Example() {
  <div style={ "will throw an error" }></div>
}
```

Class names are sanitized by default. A failed class name is replaced by `--templ-css-class-safe-name`. The sanitization can be bypassed using the `tndr.SafeClass` function, but the result is still subject to escaping.

```html
templ Example() {
  <div class={ "unsafe</style&gt;-will-sanitized", tndr.SafeClass("&sanitization bypassed") }></div>
}
```

Rendered output:

```html
<div class="--templ-css-class-safe-name &amp;sanitization bypassed"></div>
```

```html
templ Example() {
  <div>Node text is not modified at all.</div>
  <div>{ "will be escaped using tndr.EscapeString" }</div>
}
```

`href` attributes must be a `tndr.SafeURL` and are sanitized to remove JavaScript URLs unless bypassed.

```html
templ Example() {
  <a href="http://constants.example.com/are/not/sanitized">Text</a>
  <a href={ tndr.URL("will be sanitized by tndr.URL to remove potential attacks") }</a>
  <a href={ tndr.SafeURL("will not be sanitized by tndr.URL") }</a>
}
```

Within css blocks, property names, and constant CSS property values are not sanitized or escaped.

```css
css className() {
	background-color: #ffffff;
}
```

CSS property values based on expressions are passed through `tndr.SanitizeCSS` to replace potentially unsafe values with placeholders.

```css
css className() {
	color: { red };
}
```
