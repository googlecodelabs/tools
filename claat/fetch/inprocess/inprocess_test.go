package inprocess

import (
	"strings"
	"testing"
)

func TestNewInProcessFetcher(t *testing.T) {
	r := strings.NewReader("this is a string")
	ipf := NewInProcessFetcher(r)

	if ipf.source != r {
		t.Errorf("NewInProcessFetcher(%v).source = %v, want %v", r, ipf.source, r)
	}
}

func TestFetch(t *testing.T) {
	r := strings.NewReader("this is also a string")
	ipf := NewInProcessFetcher(r)

	out, err := ipf.Fetch()
	if err != nil {
		t.Errorf("Fetch() got err %v, want nil", err)
	}
	if out != r {
		t.Errorf("Fetch() = %v, want %v", out, r)
	}
}
