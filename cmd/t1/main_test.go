package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/senforsce/tndr"
)

func TestMain(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedStdout string
		expectedStderr string
		expectedCode   int
	}{
		{
			name:           "no args prints usage",
			args:           []string{},
			expectedStderr: usageText,
			expectedCode:   64, // EX_USAGE
		},
		{
			name:           `"tndr help" prints help`,
			args:           []string{"tndr", "help"},
			expectedStdout: usageText,
			expectedCode:   0,
		},
		{
			name:           `"tndr --help" prints help`,
			args:           []string{"tndr", "--help"},
			expectedStdout: usageText,
			expectedCode:   0,
		},
		{
			name:           `"tndr version" prints version`,
			args:           []string{"tndr", "version"},
			expectedStdout: tndr.Version() + "\n",
			expectedCode:   0,
		},
		{
			name:           `"tndr --version" prints version`,
			args:           []string{"tndr", "--version"},
			expectedStdout: tndr.Version() + "\n",
			expectedCode:   0,
		},
		{
			name:           `"tndr fmt --help" prints usage`,
			args:           []string{"tndr", "fmt", "--help"},
			expectedStdout: fmtUsageText,
			expectedCode:   0,
		},
		{
			name:           `"tndr lsp --help" prints usage`,
			args:           []string{"tndr", "lsp", "--help"},
			expectedStdout: lspUsageText,
			expectedCode:   0,
		},
		{
			name:           `"tndr info --help" prints usage`,
			args:           []string{"tndr", "info", "--help"},
			expectedStdout: infoUsageText,
			expectedCode:   0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stdin := strings.NewReader("")
			stdout := bytes.NewBuffer(nil)
			stderr := bytes.NewBuffer(nil)
			actualCode := run(stdin, stdout, stderr, test.args)

			if actualCode != test.expectedCode {
				t.Errorf("expected code %v, got %v", test.expectedCode, actualCode)
			}
			if diff := cmp.Diff(test.expectedStdout, stdout.String()); diff != "" {
				t.Error(diff)
				t.Error("expected stdout:")
				t.Error(test.expectedStdout)
				t.Error("actual stdout:")
				t.Error(stdout.String())
			}
			if diff := cmp.Diff(test.expectedStderr, stderr.String()); diff != "" {
				t.Error(diff)
				t.Error("expected stderr:")
				t.Error(test.expectedStderr)
				t.Error("actual stderr:")
				t.Error(stderr.String())
			}
		})
	}
}
