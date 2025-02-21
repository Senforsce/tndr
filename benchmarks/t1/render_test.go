package testhtml

import (
	"context"
	"html/template"
	"io"
	"strings"
	"testing"

	_ "embed"

	tf "github.com/senforsce/level0/templatefile"
	"github.com/senforsce/toolbelt/templatefile"
)

func BenchmarkTemplRender(b *testing.B) {
	b.ReportAllocs()
	t := Render(Person{
		Name:  "Abdoul Sy",
		Email: "example@senforsce.com",
	})

	w := new(strings.Builder)
	for i := 0; i < b.N; i++ {
		err := t.Render(context.Background(), w)
		if err != nil {
			b.Errorf("failed to render: %v", err)
		}
		w.Reset()
	}
}

//go:embed template.t1
var parserBenchmarkTemplate string

func BenchmarkTemplParser(b *testing.B) {
	tfp := tf.NewTemplateFileParser("benchmark")
	for i := 0; i < b.N; i++ {
		tf, err := templatefile.ParseString(parserBenchmarkTemplate, tfp, tf.TemplateNodeParserList)
		if err != nil {
			b.Fatal(err)
		}
		if tf.Package.Expression.Value == "" {
			b.Fatal("unexpected nil template")
		}
	}
}

var goTemplate = template.Must(template.New("example").Parse(`<div>
	<h1>{{.Name}}</h1>
	<div style="font-family: &#39;sans-serif&#39;" id="test" data-contents="something with &#34;quotes&#34; and a &lt;tag&gt;">
		<div>
			email:<a href="mailto: {{.Email}}">{{.Email}}</a></div>
		</div>
	</div>
	<hr noshade>
	<hr optionA optionB optionC="other">
	<hr noshade>
`))

func BenchmarkGoTemplateRender(b *testing.B) {
	w := new(strings.Builder)
	person := Person{
		Name:  "Abdoul Sy",
		Email: "example@senforsce.com",
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		err := goTemplate.Execute(w, person)
		if err != nil {
			b.Errorf("failed to render: %v", err)
		}
		w.Reset()
	}
}

const html = `<div><h1>Abdoul Sy</h1><div style="font-family: &#39;sans-serif&#39;" id="test" data-contents="something with &#34;quotes&#34; and a &lt;tag&gt;"><div>email:<a href="mailto: example@senforsce.com">example@senforsce.com</a></div></div></div><hr noshade><hr optionA optionB optionC="other"><hr noshade>`

func BenchmarkIOWriteString(b *testing.B) {
	b.ReportAllocs()
	w := new(strings.Builder)
	for i := 0; i < b.N; i++ {
		_, err := io.WriteString(w, html)
		if err != nil {
			b.Errorf("failed to render: %v", err)
		}
		w.Reset()
	}
}
