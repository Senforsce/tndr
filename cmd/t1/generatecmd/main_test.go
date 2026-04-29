package generatecmd

import (
	"bytes"
	"context"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/senforsce/tndr/cmd/t1/testproject"
	"github.com/senforsce/tndr/runtime"
	"golang.org/x/sync/errgroup"
)

func TestGenerate(t *testing.T) {
	t.Run("can print help", func(t *testing.T) {
		// t1 generate -help
		stdout := &bytes.Buffer{}
		err := Run(context.Background(), stdout, io.Discard, []string{"-help"})
		if err != nil {
			t.Fatalf("failed to run generate command: %v", err)
		}
		if !strings.Contains(stdout.String(), "usage: t1 generate") {
			t.Fatalf("expected help output, got: %s", stdout.String())
		}
	})
	t.Run("can generate a file in place", func(t *testing.T) {
		// t1 generate -f templates.t1
		dir, err := testproject.Create("github.com/senforsce/tndr/cmd/t1/testproject")
		if err != nil {
			t.Fatalf("failed to create test project: %v", err)
		}
		defer func() {
			if err := os.RemoveAll(dir); err != nil {
				t.Errorf("failed to remove test project directory: %v", err)
			}
		}()

		// Delete the templates_t1.go file to ensure it is generated.
		err = os.Remove(path.Join(dir, "templates_t1.go"))
		if err != nil {
			t.Fatalf("failed to remove templates_t1.go: %v", err)
		}

		// Run the generate command.
		err = Run(context.Background(), io.Discard, io.Discard, []string{"-f", path.Join(dir, "templates.t1")})
		if err != nil {
			t.Fatalf("failed to run generate command: %v", err)
		}

		// Check the templates_t1.go file was created.
		_, err = os.Stat(path.Join(dir, "templates_t1.go"))
		if err != nil {
			t.Fatalf("templates_t1.go was not created: %v", err)
		}
	})
	t.Run("can generate a file in watch mode", func(t *testing.T) {
		// tndr generate -f templates.t1
		dir, err := testproject.Create("github.com/senforsce/tndr/cmd/t1/testproject")
		if err != nil {
			t.Fatalf("failed to create test project: %v", err)
		}
		defer func() {
			if err := os.RemoveAll(dir); err != nil {
				t.Errorf("failed to remove test project directory: %v", err)
			}
		}()

		// Delete the templates_t1.go file to ensure it is generated.
		err = os.Remove(path.Join(dir, "templates_t1.go"))
		if err != nil {
			t.Fatalf("failed to remove templates_t1.go: %v", err)
		}
		ctx, cancel := context.WithCancel(context.Background())

		var eg errgroup.Group
		eg.Go(func() error {
			return Run(ctx, io.Discard, io.Discard, []string{"-path", dir, "-watch"})
		})

		// Check the templates_t1.go file was created, with backoff.
		devModeTextFileName := runtime.GetDevModeTextFileName(path.Join(dir, "templates_t1.go"))
		for i := range 5 {
			time.Sleep(time.Second * time.Duration(i))
			_, err = os.Stat(path.Join(dir, "templates_t1.go"))
			if err != nil {
				continue
			}
			_, err = os.Stat(devModeTextFileName)
			if err != nil {
				continue
			}
			break
		}
		if err != nil {
			t.Errorf("template files were not created: %v", err)
		}

		cancel()
		if err := eg.Wait(); err != nil {
			t.Errorf("generate command failed: %v", err)
		}

		// Check the templates_t1.txt file was removed.
		_, err = os.Stat(path.Join(dir, devModeTextFileName))
		if err == nil {
			t.Error("templates_t1.txt was not removed")
		}
	})
}

func TestDefaultWatchPattern(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		matches bool
	}{
		{
			name:    "empty file names do not match",
			input:   "",
			matches: false,
		},
		{
			name:    "*_t1.txt is no longer matched, Windows",
			input:   `C:\Users\asy\Senforsce\src\github.com\senforsce\tndr\cmd\t1\testproject\strings_t1.txt`,
			matches: false,
		},
		{
			name:    "*_t1.txt is no longer matched, Unix",
			input:   "/home/asy/Senforsce/src/github.com/senforsce/tndr/cmd/t1/testproject/strings_t1.txt",
			matches: false,
		},
		{
			name:    "*.t1 files match, Windows",
			input:   `C:\Users\asy\Senforsce\src\github.com\senforsce\tndr\cmd\t1\testproject\templates.t1`,
			matches: true,
		},
		{
			name:    "*.t1 files match, Unix",
			input:   "/home/asy/Senforsce/src/github.com/senforsce/tndr/cmd/t1/testproject/templates.t1",
			matches: true,
		},
		{
			name:    "*_t1.go files match, Windows",
			input:   `C:\Users\asy\Senforsce\src\github.com\senforsce\tndr\cmd\t1\testproject\templates_t1.go`,
			matches: true,
		},
		{
			name:    "*_t1.go files match, Unix",
			input:   "/home/asy/Senforsce/src/github.com/senforsce/tndr/cmd/t1/testproject/templates_t1.go",
			matches: true,
		},
		{
			name:    "*.go files match, Windows",
			input:   `C:\Users\asy\Senforsce\src\github.com\senforsce\tndr\cmd\t1\testproject\templates.go`,
			matches: true,
		},
		{
			name:    "*.go files match, Unix",
			input:   "/home/asy/Senforsce/src/github.com/senforsce/tndr/cmd/t1/testproject/templates.go",
			matches: true,
		},
		{
			name:    "*.css files do not match",
			input:   "/home/asy/Senforsce/src/github.com/senforsce/tndr/cmd/t1/testproject/templates.css",
			matches: false,
		},
	}
	wpRegexp, err := regexp.Compile(defaultWatchPattern)
	if err != nil {
		t.Fatalf("failed to compile default watch pattern: %v", err)
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if wpRegexp.MatchString(test.input) != test.matches {
				t.Fatalf("expected match of %q to be %v", test.input, test.matches)
			}
		})
	}
}

func TestArgs(t *testing.T) {
	t.Run("Help is true if the help flag is set", func(t *testing.T) {
		_, help, err := NewArguments(io.Discard, io.Discard, []string{"-help"})
		if err != nil {
			t.Fatal(err)
		}
		if !help {
			t.Fatal("expected help to be true")
		}
	})
	t.Run("Help is false if the help flag is not set", func(t *testing.T) {
		_, help, err := NewArguments(io.Discard, io.Discard, []string{})
		if err != nil {
			t.Fatal(err)
		}
		if help {
			t.Fatal("expected help to be false")
		}
	})
	t.Run("The worker count is set to the number of CPUs if not specified", func(t *testing.T) {
		args, _, err := NewArguments(io.Discard, io.Discard, []string{})
		if err != nil {
			t.Fatal(err)
		}
		if args.WorkerCount == 0 {
			t.Fatal("expected worker count to be set to the number of CPUs")
		}
	})
	t.Run("If toStdout is true, the file name must be specified", func(t *testing.T) {
		_, _, err := NewArguments(io.Discard, io.Discard, []string{"-stdout"})
		if err == nil {
			t.Fatal("expected error when toStdout is true but no file name is specified")
		}
	})
	t.Run("If toStdout is true, and the file name is specified, it writes to stdout", func(t *testing.T) {
		args, _, err := NewArguments(io.Discard, io.Discard, []string{"-stdout", "-f", "output.t1"})
		if err != nil {
			t.Fatal(err)
		}
		if args.FileName != "output.t1" {
			t.Fatalf("expected file name to be 'output.t1', got '%s'", args.FileName)
		}
		if args.FileWriter == nil {
			t.Fatal("expected FileWriter to be set when toStdout is true")
		}
	})
	t.Run("If the watchPattern is empty, it defaults to the default pattern", func(t *testing.T) {
		args, _, err := NewArguments(io.Discard, io.Discard, []string{})
		if err != nil {
			t.Fatal(err)
		}
		if args.WatchPattern.String() != defaultWatchPattern {
			t.Fatalf("expected watch pattern to be %q, got %q", defaultWatchPattern, args.WatchPattern.String())
		}
	})
	t.Run("If the watchPattern is set, it is checked for validity", func(t *testing.T) {
		_, _, err := NewArguments(io.Discard, io.Discard, []string{"-watch-pattern", "invalid[pattern"})
		if err == nil {
			t.Fatal("expected error when watch pattern is invalid")
		}
	})
	t.Run("If the watch flag is set, watch is set to true", func(t *testing.T) {
		args, _, err := NewArguments(io.Discard, io.Discard, []string{"-watch"})
		if err != nil {
			t.Fatal(err)
		}
		if !args.Watch {
			t.Fatal("expected watch to be true when the watch flag is set")
		}
	})
	t.Run("If the watch flag is not set, watch is false", func(t *testing.T) {
		args, _, err := NewArguments(io.Discard, io.Discard, []string{})
		if err != nil {
			t.Fatal(err)
		}
		if args.Watch {
			t.Fatal("expected watch to be false when the watch flag is not set")
		}
	})
	t.Run("The cmd flag can be set to specify a command to run after generating", func(t *testing.T) {
		args, _, err := NewArguments(io.Discard, io.Discard, []string{"-cmd", "echo hello"})
		if err != nil {
			t.Fatal(err)
		}
		if args.Command != "echo hello" {
			t.Fatalf("expected command to be 'echo hello', got '%s'", args.Command)
		}
	})
}
