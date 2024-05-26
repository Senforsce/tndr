// Code generated by t1 - DO NOT EDIT.

package testcssusage

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/senforsce/t1"
import "context"
import "io"
import "bytes"
import "strings"

import "fmt"

// Constant class.
func StyleTagsAreSupported() t1.Component {
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
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString("<style>\n\t.test {\n\t\tcolor: #ff0000;\n\t}\n\t</style><div class=\"test\">Style tags are supported</div>")
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		if !t1_7745c5c3_IsBuffer {
			_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteTo(t1_7745c5c3_W)
		}
		return t1_7745c5c3_Err
	})
}

// CSS components.

const red = "#00ff00"

func cssComponentGreen() t1.CSSClass {
	var t1_7745c5c3_CSSBuilder strings.Builder
	t1_7745c5c3_CSSBuilder.WriteString(string(t1.SanitizeCSS(`color`, red)))
	t1_7745c5c3_CSSID := t1.CSSID(`cssComponentGreen`, t1_7745c5c3_CSSBuilder.String())
	return t1.ComponentCSSClass{
		ID:    t1_7745c5c3_CSSID,
		Class: t1.SafeCSS(`.` + t1_7745c5c3_CSSID + `{` + t1_7745c5c3_CSSBuilder.String() + `}`),
	}
}

func CSSComponentsAreSupported() t1.Component {
	return t1.ComponentFunc(func(ctx context.Context, t1_7745c5c3_W io.Writer) (t1_7745c5c3_Err error) {
		t1_7745c5c3_Buffer, t1_7745c5c3_IsBuffer := t1_7745c5c3_W.(*bytes.Buffer)
		if !t1_7745c5c3_IsBuffer {
			t1_7745c5c3_Buffer = t1.GetBuffer()
			defer t1.ReleaseBuffer(t1_7745c5c3_Buffer)
		}
		ctx = t1.InitializeContext(ctx)
		t1_7745c5c3_Var2 := t1.GetChildren(ctx)
		if t1_7745c5c3_Var2 == nil {
			t1_7745c5c3_Var2 = t1.NopComponent
		}
		ctx = t1.ClearChildren(ctx)
		var t1_7745c5c3_Var3 = []any{cssComponentGreen()}
		t1_7745c5c3_Err = t1.RenderCSSItems(ctx, t1_7745c5c3_Buffer, t1_7745c5c3_Var3...)
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString("<div class=\"")
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString(t1.EscapeString(t1.CSSClasses(t1_7745c5c3_Var3).String()))
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString("\">CSS components are supported</div>")
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		if !t1_7745c5c3_IsBuffer {
			_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteTo(t1_7745c5c3_W)
		}
		return t1_7745c5c3_Err
	})
}

// Both CSS components and constants are supported.
// Only string names are really required. There is no need to use t1.Class or t1.SafeClass.
func CSSComponentsAndConstantsAreSupported() t1.Component {
	return t1.ComponentFunc(func(ctx context.Context, t1_7745c5c3_W io.Writer) (t1_7745c5c3_Err error) {
		t1_7745c5c3_Buffer, t1_7745c5c3_IsBuffer := t1_7745c5c3_W.(*bytes.Buffer)
		if !t1_7745c5c3_IsBuffer {
			t1_7745c5c3_Buffer = t1.GetBuffer()
			defer t1.ReleaseBuffer(t1_7745c5c3_Buffer)
		}
		ctx = t1.InitializeContext(ctx)
		t1_7745c5c3_Var4 := t1.GetChildren(ctx)
		if t1_7745c5c3_Var4 == nil {
			t1_7745c5c3_Var4 = t1.NopComponent
		}
		ctx = t1.ClearChildren(ctx)
		var t1_7745c5c3_Var5 = []any{cssComponentGreen(), "classA", t1.Class("&&&classB"), t1.SafeClass("classC"), "d e"}
		t1_7745c5c3_Err = t1.RenderCSSItems(ctx, t1_7745c5c3_Buffer, t1_7745c5c3_Var5...)
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString("<div class=\"")
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString(t1.EscapeString(t1.CSSClasses(t1_7745c5c3_Var5).String()))
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString("\" type=\"button\">Both CSS components and constants are supported</div>")
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		var t1_7745c5c3_Var6 = []any{t1.Classes(cssComponentGreen(), "classA", t1.Class("&&&classB"), t1.SafeClass("classC")), "d e"}
		t1_7745c5c3_Err = t1.RenderCSSItems(ctx, t1_7745c5c3_Buffer, t1_7745c5c3_Var6...)
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString("<div class=\"")
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString(t1.EscapeString(t1.CSSClasses(t1_7745c5c3_Var6).String()))
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString("\" type=\"button\">Both CSS components and constants are supported</div>")
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		if !t1_7745c5c3_IsBuffer {
			_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteTo(t1_7745c5c3_W)
		}
		return t1_7745c5c3_Err
	})
}

// Maps can be used to determine if a class should be added or not.
func MapsCanBeUsedToConditionallySetClasses() t1.Component {
	return t1.ComponentFunc(func(ctx context.Context, t1_7745c5c3_W io.Writer) (t1_7745c5c3_Err error) {
		t1_7745c5c3_Buffer, t1_7745c5c3_IsBuffer := t1_7745c5c3_W.(*bytes.Buffer)
		if !t1_7745c5c3_IsBuffer {
			t1_7745c5c3_Buffer = t1.GetBuffer()
			defer t1.ReleaseBuffer(t1_7745c5c3_Buffer)
		}
		ctx = t1.InitializeContext(ctx)
		t1_7745c5c3_Var7 := t1.GetChildren(ctx)
		if t1_7745c5c3_Var7 == nil {
			t1_7745c5c3_Var7 = t1.NopComponent
		}
		ctx = t1.ClearChildren(ctx)
		var t1_7745c5c3_Var8 = []any{map[string]bool{"a": true, "b": false, "c": true}}
		t1_7745c5c3_Err = t1.RenderCSSItems(ctx, t1_7745c5c3_Buffer, t1_7745c5c3_Var8...)
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString("<div class=\"")
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString(t1.EscapeString(t1.CSSClasses(t1_7745c5c3_Var8).String()))
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString("\">Maps can be used to determine if a class should be added or not.</div>")
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		if !t1_7745c5c3_IsBuffer {
			_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteTo(t1_7745c5c3_W)
		}
		return t1_7745c5c3_Err
	})
}

// The t1.KV function can be used to add a class if a condition is true.
func d() t1.CSSClass {
	var t1_7745c5c3_CSSBuilder strings.Builder
	t1_7745c5c3_CSSBuilder.WriteString(`font-size:12pt;`)
	t1_7745c5c3_CSSID := t1.CSSID(`d`, t1_7745c5c3_CSSBuilder.String())
	return t1.ComponentCSSClass{
		ID:    t1_7745c5c3_CSSID,
		Class: t1.SafeCSS(`.` + t1_7745c5c3_CSSID + `{` + t1_7745c5c3_CSSBuilder.String() + `}`),
	}
}

func e() t1.CSSClass {
	var t1_7745c5c3_CSSBuilder strings.Builder
	t1_7745c5c3_CSSBuilder.WriteString(`font-size:14pt;`)
	t1_7745c5c3_CSSID := t1.CSSID(`e`, t1_7745c5c3_CSSBuilder.String())
	return t1.ComponentCSSClass{
		ID:    t1_7745c5c3_CSSID,
		Class: t1.SafeCSS(`.` + t1_7745c5c3_CSSID + `{` + t1_7745c5c3_CSSBuilder.String() + `}`),
	}
}

func KVCanBeUsedToConditionallySetClasses() t1.Component {
	return t1.ComponentFunc(func(ctx context.Context, t1_7745c5c3_W io.Writer) (t1_7745c5c3_Err error) {
		t1_7745c5c3_Buffer, t1_7745c5c3_IsBuffer := t1_7745c5c3_W.(*bytes.Buffer)
		if !t1_7745c5c3_IsBuffer {
			t1_7745c5c3_Buffer = t1.GetBuffer()
			defer t1.ReleaseBuffer(t1_7745c5c3_Buffer)
		}
		ctx = t1.InitializeContext(ctx)
		t1_7745c5c3_Var9 := t1.GetChildren(ctx)
		if t1_7745c5c3_Var9 == nil {
			t1_7745c5c3_Var9 = t1.NopComponent
		}
		ctx = t1.ClearChildren(ctx)
		var t1_7745c5c3_Var10 = []any{"a", t1.KV("b", false), "c", t1.KV(d(), false), t1.KV(e(), true)}
		t1_7745c5c3_Err = t1.RenderCSSItems(ctx, t1_7745c5c3_Buffer, t1_7745c5c3_Var10...)
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString("<div class=\"")
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString(t1.EscapeString(t1.CSSClasses(t1_7745c5c3_Var10).String()))
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString("\">KV can be used to conditionally set classes.</div>")
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		if !t1_7745c5c3_IsBuffer {
			_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteTo(t1_7745c5c3_W)
		}
		return t1_7745c5c3_Err
	})
}

// Pseudo attributes can be used without any special syntax.
func PsuedoAttributesAndComplexClassNamesAreSupported() t1.Component {
	return t1.ComponentFunc(func(ctx context.Context, t1_7745c5c3_W io.Writer) (t1_7745c5c3_Err error) {
		t1_7745c5c3_Buffer, t1_7745c5c3_IsBuffer := t1_7745c5c3_W.(*bytes.Buffer)
		if !t1_7745c5c3_IsBuffer {
			t1_7745c5c3_Buffer = t1.GetBuffer()
			defer t1.ReleaseBuffer(t1_7745c5c3_Buffer)
		}
		ctx = t1.InitializeContext(ctx)
		t1_7745c5c3_Var11 := t1.GetChildren(ctx)
		if t1_7745c5c3_Var11 == nil {
			t1_7745c5c3_Var11 = t1.NopComponent
		}
		ctx = t1.ClearChildren(ctx)
		var t1_7745c5c3_Var12 = []any{"bg-violet-500", "hover:bg-red-600", "hover:bg-sky-700", "text-[#50d71e]", "w-[calc(100%-4rem)"}
		t1_7745c5c3_Err = t1.RenderCSSItems(ctx, t1_7745c5c3_Buffer, t1_7745c5c3_Var12...)
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString("<div class=\"")
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString(t1.EscapeString(t1.CSSClasses(t1_7745c5c3_Var12).String()))
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString("\">Psuedo attributes and complex class names are supported.</div>")
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		if !t1_7745c5c3_IsBuffer {
			_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteTo(t1_7745c5c3_W)
		}
		return t1_7745c5c3_Err
	})
}

// Class names are HTML escaped.
func ClassNamesAreHTMLEscaped() t1.Component {
	return t1.ComponentFunc(func(ctx context.Context, t1_7745c5c3_W io.Writer) (t1_7745c5c3_Err error) {
		t1_7745c5c3_Buffer, t1_7745c5c3_IsBuffer := t1_7745c5c3_W.(*bytes.Buffer)
		if !t1_7745c5c3_IsBuffer {
			t1_7745c5c3_Buffer = t1.GetBuffer()
			defer t1.ReleaseBuffer(t1_7745c5c3_Buffer)
		}
		ctx = t1.InitializeContext(ctx)
		t1_7745c5c3_Var13 := t1.GetChildren(ctx)
		if t1_7745c5c3_Var13 == nil {
			t1_7745c5c3_Var13 = t1.NopComponent
		}
		ctx = t1.ClearChildren(ctx)
		var t1_7745c5c3_Var14 = []any{"a\" onClick=\"alert('hello')\""}
		t1_7745c5c3_Err = t1.RenderCSSItems(ctx, t1_7745c5c3_Buffer, t1_7745c5c3_Var14...)
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString("<div class=\"")
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString(t1.EscapeString(t1.CSSClasses(t1_7745c5c3_Var14).String()))
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString("\">Class names are HTML escaped.</div>")
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		if !t1_7745c5c3_IsBuffer {
			_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteTo(t1_7745c5c3_W)
		}
		return t1_7745c5c3_Err
	})
}

// CSS components can be used with arguments.
func loading(percent int) t1.CSSClass {
	var t1_7745c5c3_CSSBuilder strings.Builder
	t1_7745c5c3_CSSBuilder.WriteString(string(t1.SanitizeCSS(`width`, fmt.Sprintf("%d%%", percent))))
	t1_7745c5c3_CSSID := t1.CSSID(`loading`, t1_7745c5c3_CSSBuilder.String())
	return t1.ComponentCSSClass{
		ID:    t1_7745c5c3_CSSID,
		Class: t1.SafeCSS(`.` + t1_7745c5c3_CSSID + `{` + t1_7745c5c3_CSSBuilder.String() + `}`),
	}
}

func CSSComponentsCanBeUsedWithArguments() t1.Component {
	return t1.ComponentFunc(func(ctx context.Context, t1_7745c5c3_W io.Writer) (t1_7745c5c3_Err error) {
		t1_7745c5c3_Buffer, t1_7745c5c3_IsBuffer := t1_7745c5c3_W.(*bytes.Buffer)
		if !t1_7745c5c3_IsBuffer {
			t1_7745c5c3_Buffer = t1.GetBuffer()
			defer t1.ReleaseBuffer(t1_7745c5c3_Buffer)
		}
		ctx = t1.InitializeContext(ctx)
		t1_7745c5c3_Var15 := t1.GetChildren(ctx)
		if t1_7745c5c3_Var15 == nil {
			t1_7745c5c3_Var15 = t1.NopComponent
		}
		ctx = t1.ClearChildren(ctx)
		var t1_7745c5c3_Var16 = []any{loading(50)}
		t1_7745c5c3_Err = t1.RenderCSSItems(ctx, t1_7745c5c3_Buffer, t1_7745c5c3_Var16...)
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString("<div class=\"")
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString(t1.EscapeString(t1.CSSClasses(t1_7745c5c3_Var16).String()))
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString("\">CSS components can be used with arguments.</div>")
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		var t1_7745c5c3_Var17 = []any{loading(100)}
		t1_7745c5c3_Err = t1.RenderCSSItems(ctx, t1_7745c5c3_Buffer, t1_7745c5c3_Var17...)
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString("<div class=\"")
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString(t1.EscapeString(t1.CSSClasses(t1_7745c5c3_Var17).String()))
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteString("\">CSS components can be used with arguments.</div>")
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		if !t1_7745c5c3_IsBuffer {
			_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteTo(t1_7745c5c3_W)
		}
		return t1_7745c5c3_Err
	})
}

// Combine all tests.
func TestComponent() t1.Component {
	return t1.ComponentFunc(func(ctx context.Context, t1_7745c5c3_W io.Writer) (t1_7745c5c3_Err error) {
		t1_7745c5c3_Buffer, t1_7745c5c3_IsBuffer := t1_7745c5c3_W.(*bytes.Buffer)
		if !t1_7745c5c3_IsBuffer {
			t1_7745c5c3_Buffer = t1.GetBuffer()
			defer t1.ReleaseBuffer(t1_7745c5c3_Buffer)
		}
		ctx = t1.InitializeContext(ctx)
		t1_7745c5c3_Var18 := t1.GetChildren(ctx)
		if t1_7745c5c3_Var18 == nil {
			t1_7745c5c3_Var18 = t1.NopComponent
		}
		ctx = t1.ClearChildren(ctx)
		t1_7745c5c3_Err = StyleTagsAreSupported().Render(ctx, t1_7745c5c3_Buffer)
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		t1_7745c5c3_Err = CSSComponentsAreSupported().Render(ctx, t1_7745c5c3_Buffer)
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		t1_7745c5c3_Err = CSSComponentsAndConstantsAreSupported().Render(ctx, t1_7745c5c3_Buffer)
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		t1_7745c5c3_Err = MapsCanBeUsedToConditionallySetClasses().Render(ctx, t1_7745c5c3_Buffer)
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		t1_7745c5c3_Err = KVCanBeUsedToConditionallySetClasses().Render(ctx, t1_7745c5c3_Buffer)
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		t1_7745c5c3_Err = PsuedoAttributesAndComplexClassNamesAreSupported().Render(ctx, t1_7745c5c3_Buffer)
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		t1_7745c5c3_Err = ClassNamesAreHTMLEscaped().Render(ctx, t1_7745c5c3_Buffer)
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		t1_7745c5c3_Err = CSSComponentsCanBeUsedWithArguments().Render(ctx, t1_7745c5c3_Buffer)
		if t1_7745c5c3_Err != nil {
			return t1_7745c5c3_Err
		}
		if !t1_7745c5c3_IsBuffer {
			_, t1_7745c5c3_Err = t1_7745c5c3_Buffer.WriteTo(t1_7745c5c3_W)
		}
		return t1_7745c5c3_Err
	})
}
