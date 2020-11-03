package collectIstags

import (
	"testing"
)

func TestCleanImages(t *testing.T) {
	expected := "Hallo Welt!"
	if ret := Hello("Welt"); ret != expected {
		t.Errorf("Hello() = %q, want %q", ret, expected)
	}
}
