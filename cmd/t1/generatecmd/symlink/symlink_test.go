package symlink

import (
	"context"
	"io"
	"os"
	"path"
	"testing"

	"github.com/senforsce/tndr/cmd/t1/generatecmd"
	"github.com/senforsce/tndr/cmd/t1/testproject"
)

func TestSymlink(t *testing.T) {
	t.Run("can generate if root is symlink", func(t *testing.T) {
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

		symlinkPath := dir + "-symlink"
		err = os.Symlink(dir, symlinkPath)
		if err != nil {
			t.Fatalf("failed to create dir symlink: %v", err)
		}
		defer func() {
			if err = os.Remove(symlinkPath); err != nil {
				t.Errorf("failed to remove symlink directory: %v", err)
			}
		}()

		// Delete the templates_t1.go file to ensure it is generated.
		err = os.Remove(path.Join(symlinkPath, "templates_t1.go"))
		if err != nil {
			t.Fatalf("failed to remove templates_t1.go: %v", err)
		}

		// Run the generate command.
		err = generatecmd.Run(context.Background(), io.Discard, io.Discard, []string{"-path", symlinkPath})
		if err != nil {
			t.Fatalf("failed to run generate command: %v", err)
		}

		// Check the templates_t1.go file was created.
		_, err = os.Stat(path.Join(symlinkPath, "templates_t1.go"))
		if err != nil {
			t.Fatalf("templates_t1.go was not created: %v", err)
		}
	})
}
