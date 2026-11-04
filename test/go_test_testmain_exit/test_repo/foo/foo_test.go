package foo

import "testing"

// TestMain deliberately doesn't call os.Exit, like Go 1.15+ allows.
func TestMain(m *testing.M) {
	m.Run()
}

func TestFails(t *testing.T) {
	t.Fatal("failing on purpose")
}
