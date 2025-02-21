package visualize

import (
	"context"
	"fmt"
	"html"
	"io"
	"strconv"
	"strings"

	t1 "github.com/senforsce/tndr"
	"github.com/senforsce/toolbelt/sourcemap"
)

func HTML(t1FileName string, t1Contents, goContents string, sourceMap *sourcemap.SourceMap) t1.Component {
	tl := t1Lines{contents: string(t1Contents), sourceMap: sourceMap}
	gl := goLines{contents: string(goContents), sourceMap: sourceMap}
	return combine(t1FileName, tl, gl)
}

type t1Lines struct {
	contents  string
	sourceMap *sourcemap.SourceMap
}

func (tl t1Lines) Render(ctx context.Context, w io.Writer) (err error) {
	t1Lines := strings.Split(tl.contents, "\n")
	for lineIndex, line := range t1Lines {
		if _, err = w.Write([]byte("<span>" + strconv.Itoa(lineIndex) + "&nbsp;</span>\n")); err != nil {
			return
		}
		for colIndex, c := range line {
			if tgt, ok := tl.sourceMap.TargetPositionFromSource(int(lineIndex), int(colIndex)); ok {
				sourceID := fmt.Sprintf("src_%d_%d", lineIndex, colIndex)
				targetID := fmt.Sprintf("tgt_%d_%d", tgt.Line, tgt.Col)
				if err := mappedCharacter(string(c), sourceID, targetID).Render(ctx, w); err != nil {
					return err
				}
			} else {
				s := html.EscapeString(string(c))
				s = strings.ReplaceAll(s, "\t", "&nbsp;")
				s = strings.ReplaceAll(s, " ", "&nbsp;")
				if _, err := w.Write([]byte(s)); err != nil {
					return err
				}
			}
		}
		if _, err = w.Write([]byte("\n<br/>\n")); err != nil {
			return
		}
	}
	return nil
}

type goLines struct {
	contents  string
	sourceMap *sourcemap.SourceMap
}

func (gl goLines) Render(ctx context.Context, w io.Writer) (err error) {
	t1Lines := strings.Split(gl.contents, "\n")
	for lineIndex, line := range t1Lines {
		if _, err = w.Write([]byte("<span>" + strconv.Itoa(lineIndex) + "&nbsp;</span>\n")); err != nil {
			return
		}
		for colIndex, c := range line {
			if src, ok := gl.sourceMap.SourcePositionFromTarget(int(lineIndex), int(colIndex)); ok {
				sourceID := fmt.Sprintf("src_%d_%d", src.Line, src.Col)
				targetID := fmt.Sprintf("tgt_%d_%d", lineIndex, colIndex)
				if err := mappedCharacter(string(c), sourceID, targetID).Render(ctx, w); err != nil {
					return err
				}
			} else {
				s := html.EscapeString(string(c))
				s = strings.ReplaceAll(s, "\t", "&nbsp;")
				s = strings.ReplaceAll(s, " ", "&nbsp;")
				if _, err := w.Write([]byte(s)); err != nil {
					return err
				}
			}
		}
		if _, err = w.Write([]byte("\n<br/>\n")); err != nil {
			return
		}
	}
	return nil
}
