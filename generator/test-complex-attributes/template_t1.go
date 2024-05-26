// Code generated by t1 - DO NOT EDIT.

package testcomplexattributes

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/senforsce/t1"
import "context"
import "io"
import "bytes"

func ComplexAttributes() t1.Component {
	return t1.ComponentFunc(func(ctx context.Context, t1_7745c5c3_W io.Writer) (t1_7745c5c3_Err error) {
		t1_7745c5c3_Buffer, t1_7745c5c3_IsBuffer := t1_7745c5c3_W.(*bytes.Buffer)
		if !t1_7745c5c3_IsBuffer {
			t1_7745c5c3_Buffer = t1.GetBuffer()
			defer t1.ReleaseBuffer(t1_7745c5c3_Buffer)
		}
		ctx = t1.InitializeContext(ctx)
		t1_7745c5c3_Var1 := t1.GetChildren(ctx)
		if t1_7745c5c3_Var1 == nil {
			t1_7745c5c3_Var1 = t1.NopComponent
		}
		ctx = t1.ClearChildren(ctx)
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString("<div x-data=\"{darkMode: localStorage.getItem(&#39;darkMode&#39;) || localStorage.setItem(&#39;darkMode&#39;, &#39;system&#39;)}\" x-init=\"$watch(&#39;darkMode&#39;, val =&gt; localStorage.setItem(&#39;darkMode&#39;, val))\" :class=\"{&#39;dark&#39;: darkMode === &#39;dark&#39; || (darkMode === &#39;system&#39; &amp;&amp; window.matchMedia(&#39;(prefers-color-scheme: dark)&#39;).matches)}\"></div><div x-data=\"{ count: 0 }\"><button x-on:click=\"count++\">Increment</button> <span x-text=\"count\"></span></div><div x-data=\"{ count: 0 }\"><button @click=\"count++\">Increment</button> <span x-text=\"count\"></span></div>")
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		if !t1_7745c5c3_IsBuffer {
			_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteTo(t1_7745c5c3_W)
		}
		return t1_7745c5c3_Err
	})
}
