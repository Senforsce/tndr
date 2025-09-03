package format

import (
	"github.com/senforsce/tndr/internal/prettier"
	"github.com/senforsce/tndr/parser/v2"
)

func StyleElement(se *parser.RawElement, depth int, prettierCommand string) (err error) {
	if se.Name != "style" {
		return nil
	}

	// Skip empty style elements, as they don't need formatting.
	if len(se.Contents) == 0 {
		return nil
	}

	// Prettyify the style contents.
	se.Contents, err = prettier.Element("style", "text/css", se.Contents, depth, prettierCommand)
	if err != nil {
		return err
	}

	return nil
}
