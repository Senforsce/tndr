package runtime

import (
	"testing"
)

func TestWatchMode(t *testing.T) {
	t.Setenv("TNDR_DEV_MODE_ROOT", "/tmp")

	t.Run("GetDevModeTextFileName respects the TNDR_DEV_MODE_ROOT environment variable", func(t *testing.T) {
		expected := "/tmp/t1_d8e3d8e2a7fd6a65f6fc55c304c18c144deed9e01c6f41505b81b7d20c8adad4.txt"
		actual := GetDevModeTextFileName("test.t1")
		if actual != expected {
			t.Errorf("got %q, want %q", actual, expected)
		}
	})
	t.Run("GetDevModeTextFileName replaces _tndr.go with .t1", func(t *testing.T) {
		expected := "/tmp/t1_7dde161b5c5e1382c2285c1deee7799935ee7d0f9b7b30f85a05d6d511a554dd.txt"
		actual := GetDevModeTextFileName("test_tndr.go")
		if actual != expected {
			t.Errorf("got %q, want %q", actual, expected)
		}
	})
	t.Run("GetDevModeTextFileName accepts absolute Linux paths", func(t *testing.T) {
		expected := "/tmp/t1_a6bb7103e12402b472a90cb9b244b5d61bb93208c9523be78ce27e11b6134d3d.txt"
		actual := GetDevModeTextFileName("/home/user/test.t1")
		if actual != expected {
			t.Errorf("got %q, want %q", actual, expected)
		}
	})
	t.Run("GetDevModeTextFileName accepts absolute Windows paths, which are normalized to Unix style before hashing", func(t *testing.T) {
		expected := "/tmp/t1_585a40d058029305fb5dadb8e6a9da8a7de57d148150992f4789f842390cf553.txt"
		actual := GetDevModeTextFileName(`C:\Windows\System32\test.t1`)
		if actual != expected {
			t.Errorf("got %q, want %q", actual, expected)
		}
	})
}
