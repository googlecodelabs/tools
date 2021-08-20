package types

import (
	"bytes"
	"time"
)

// Context is an export context.
// It is defined in this package so that it can be used by both cli and a server.
type Context struct {
	Env     string       `json:"environment"`       // Current export environment
	Format  string       `json:"format"`            // Output format, e.g. "html"
	Prefix  string       `json:"prefix,omitempty"`  // Assets URL prefix for HTML-based formats
	MainGA  string       `json:"mainga,omitempty"`  // Global Google Analytics ID
	Updated *ContextTime `json:"updated,omitempty"` // Last update timestamp
}

// ContextMeta is a composition of export context and meta data.
type ContextMeta struct {
	Context
	Meta
}

// ContextTime is a wrapper around time.Time so we can implement JSON marshalling.
type ContextTime time.Time

// MarshalJSON implements Marshaler interface.
// The output format is RFC3339.
func (ct ContextTime) MarshalJSON() ([]byte, error) {
	v := time.Time(ct).Format(time.RFC3339)
	b := make([]byte, len(v)+2)
	b[0] = '"'
	b[len(b)-1] = '"'
	copy(b[1:], v)
	return b, nil
}

// UnmarshalJSON implements Unmarshaler interface.
// Accepted formats:
// - RFC3339
// - YYYY-MM-DD
func (ct *ContextTime) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(b, `"`)
	t, err := time.Parse(time.RFC3339, string(b))
	if err != nil {
		t, err = time.Parse("2006-01-02", string(b))
	}
	if err != nil {
		return err
	}
	*ct = ContextTime(t)
	return nil
}
