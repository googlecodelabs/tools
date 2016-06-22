package fs

import "testing"

func TestCleanPath(t *testing.T) {
	table := []struct {
		root     string
		parts    []string
		expected string
	}{
		{"root", []string{"one", "two"}, "root/one/two"},
		{"/root", []string{"one", "two"}, "/root/one/two"},
		{"root", []string{"../one", "../two"}, "root/one/two"},
		{"root", []string{"../../one/", "../../two", "/"}, "root/one/two"},
	}
	for i, test := range table {
		v := cleanPath(test.root, test.parts...)
		if v != test.expected {
			t.Errorf("%d: cleanPath(%q, %q) = %q; want %q", i, test.root, test.parts, v, test.expected)
		}
	}
}
