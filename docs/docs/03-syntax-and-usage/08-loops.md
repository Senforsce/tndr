# For loops

Use the standard Go `for` loop for iteration.

```tndr title="component.t1"
package main

t1 nameList(items []Item) {
  <ul>
  for _, item := range items {
    <li>{ item.Name }</li>
  }
  </ul>
}
```

```html title="Output"
<ul>
  <li>A</li>
  <li>B</li>
  <li>C</li>
</ul>
```
