package types

import (
	"bytes"
	"time"
)

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
