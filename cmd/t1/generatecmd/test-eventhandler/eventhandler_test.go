package testeventhandler

import (
	"context"
	"errors"
	"fmt"
	"go/scanner"
	"go/token"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/google/go-cmp/cmp"
	"github.com/senforsce/generator"
	"github.com/senforsce/tndr/cmd/t1/generatecmd"
)

// extractErrorList unwraps errors until it finds a scanner.ErrorList
func extractErrorList(err error) (scanner.ErrorList, bool) {
	if err == nil {
		return nil, false
	}

	if list, ok := err.(scanner.ErrorList); ok {
		return list, true
	}

	return extractErrorList(errors.Unwrap(err))
}

type MessageAndPosition struct {
	Position token.Position
	Msg      string
}

func TestErrorLocationMapping(t *testing.T) {
	tests := []struct {
		name           string
		rawFileName    string
		errorPositions []MessageAndPosition
	}{
		{
			name:        "single error outputs location in srcFile",
			rawFileName: "single_error.t1.error",
			errorPositions: []MessageAndPosition{
				{
					Position: token.Position{Offset: 43, Line: 3, Column: 17},
					Msg:      "missing ',' in parameter list",
				},
			},
		},
		{
			name:        "multiple errors all output locations in srcFile",
			rawFileName: "multiple_errors.t1.error",
			errorPositions: []MessageAndPosition{
				{
					Position: token.Position{Offset: 41, Line: 3, Column: 15},
					Msg:      "missing ',' in parameter list",
				},
				{
					Position: token.Position{Offset: 98, Line: 7, Column: 19},
					Msg:      "missing ',' in parameter list",
				},
				{
					Position: token.Position{Offset: 122, Line: 10, Column: 0},
					Msg:      "illegal character U+00A7 '§'",
				},
			},
		},
	}

	slog := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	var fw generatecmd.FileWriterFunc
	fseh := generatecmd.NewFSEventHandler(slog, ".", false, []generator.GenerateOpt{}, false, false, fw, false)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// The raw files cannot end in .t1 because they will cause the generator to fail. Instead,
			// we create a tmp file that ends in .t1 only for the duration of the test.
			rawFile, err := os.Open(test.rawFileName)
			if err != nil {
				t.Fatalf("Failed to open file %s: %v", test.rawFileName, err)
			}
			defer func() {
				if err = rawFile.Close(); err != nil {
					t.Fatalf("Failed to close raw file %s: %v", test.rawFileName, err)
				}
			}()

			file, err := os.CreateTemp("", fmt.Sprintf("*%s.t1", test.rawFileName))
			if err != nil {
				t.Fatalf("Failed to create a tmp file at %s: %v", file.Name(), err)
			}
			tempFileName := file.Name()
			defer func() {
				_ = file.Close()
				if err := os.Remove(tempFileName); err != nil {
					t.Logf("Warning: Failed to remove tmp file %s: %v", tempFileName, err)
				}
			}()

			if _, err = io.Copy(file, rawFile); err != nil {
				t.Fatalf("Failed to copy contents from raw file %s to tmp %s: %v", test.rawFileName, tempFileName, err)
			}

			// Ensure file is synced to disk and file pointer is at the beginning
			if err = file.Sync(); err != nil {
				t.Fatalf("Failed to sync file: %v", err)
			}

			event := fsnotify.Event{Name: tempFileName, Op: fsnotify.Write}
			_, err = fseh.HandleEvent(context.Background(), event)
			if err == nil {
				t.Fatal("Expected an error but none was thrown")
			}

			list, ok := extractErrorList(err)
			if !ok {
				t.Fatal("Failed to extract ErrorList from error")
			}

			if len(list) != len(test.errorPositions) {
				for i, err := range list {
					expected := test.errorPositions[i]
					expected.Position.Filename = tempFileName
					t.Errorf("Expected Filename=%s", tempFileName)
					t.Errorf("sen:Error: %s ::Offset=%d -> Line=%d -> Column=%d", err.Msg, err.Pos.Offset, err.Pos.Line, err.Pos.Column)

					if diff := cmp.Diff(expected.Position, err.Pos); diff != "" {
						t.Errorf("Error position mismatch (-expected +actual):\n%s", diff)
					}
				}
				for _, msgAndPosition := range test.errorPositions {
					t.Errorf("sen:Error: Offset=%d -> Line=%d -> Column=%d", msgAndPosition.Position.Offset, msgAndPosition.Position.Line, msgAndPosition.Position.Column)

				}
				t.Fatalf("Expected %d errors but got %d", len(test.errorPositions), len(list))
			}

			for i, err := range list {
				expected := test.errorPositions[i]
				expected.Position.Filename = tempFileName

				if diff := cmp.Diff(expected.Position, err.Pos); diff != "" {
					t.Errorf("Error position mismatch (-expected +actual):\n%s", diff)
				}
			}
		})
	}
}
