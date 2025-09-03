# Statements

## Control flow

Within a tndr element, a subset of Go statements can be used directly.

These Go statements can be used to conditionally render child elements, or to iterate variables.

For individual implementation guides see:

* [if/else](/syntax-and-usage/if-else)
* [switch](/syntax-and-usage/switch)
* [for loops](/syntax-and-usage/loops)

## if/switch/for within text

Go statements can be used without any escaping to make it simple for developers to include them.

The tndr parser assumes that text that starts with `if`, `switch` or `for` denotes the start of one of those expressions as per this example.

```t1 title="show-hello.t1"
package main

t1 showHelloIfTrue(b bool) {
	<div>
		if b {
			<p>Hello</p>
		}
	</div>
}
```

If you need to start a text block with the words `if`, `switch`, or `for`:

* Use a Go string expression.
* Capitalise `if`, `switch`, or `for`.

```t1 title="paragraph.t1"
package main

t1 display(price float64, count int) {
	<p>Switch to Linux</p>
	<p>{ `switch to Linux` }</p>
	<p>{ "for a day" }</p>
	<p>{ fmt.Sprintf("%f", price) }{ "for" }{ fmt.Sprintf("%d", count) }</p>
	<p>{ fmt.Sprintf("%f for %d", price, count) }</p>
}
```

## Design considerations

We decided to not require a special prefix for `if`, `switch` and `for` expressions on the basis that we were more likely to want to use a Go control statement than start a text run with those strings.

To reduce the risk of a broken control statement, resulting in printing out the source code of the application, tndr will complain if a text run starts with `if`, `switch` or `for`, but no opening brace `{` is found.

For example, the following code causes the tndr parser to return an error:

```t1 title="broken-if.t1"
package main

t1 showIfTrue(b bool) {
	if b 
	  <p>Hello</p>
	}
}
```

:::note
Note the missing `{` on line 4.
:::

The following code also produces an error, since the text run starts with `if`, but no opening `{` is found.

```t1 title="paragraph.t1"
package main

t1 text(b bool) {
	<p>if a tree fell in the woods</p>
}
```

:::note
This also applies to `for` and `switch` statements.
:::

To resolve the issue:

* Use a Go string expression.
* Capitalise `if`, `switch`, or `for`.

```t1 title="paragraph.t1"
package main

t1 display(price float64, count int) {
	<p>Switch to Linux</p>
	<p>{ `switch to Linux` }</p>
	<p>{ "for a day" }</p>
	<p>{ fmt.Sprintf("%f", price) }{ "for" }{ fmt.Sprintf("%d", count) }</p>
	<p>{ fmt.Sprintf("%f for %d", price, count) }</p>
}
```
