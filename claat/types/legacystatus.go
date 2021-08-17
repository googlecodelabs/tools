package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// LegacyStatus supports legacy status values which are strings
// as opposed to an array, e.g. "['one', u'two', ...]".
type LegacyStatus []string

// MarshalJSON implements Marshaler interface.
func (s LegacyStatus) MarshalJSON() ([]byte, error) {
	if len(s) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal([]string(s))
}

// UnmarshalJSON implements Unmarshaler interface.
func (s *LegacyStatus) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	if b[0] == '"' {
		// legacy status: "['s1', u's2', ...]"
		// assume no status value contains single quotes
		b = bytes.Trim(b, `"`)
		b = bytes.Replace(b, []byte("u'"), []byte(`"`), -1)
		b = bytes.Replace(b, []byte("'"), []byte(`"`), -1)
	}
	var v []string
	if err := json.Unmarshal(b, &v); err != nil {
		return fmt.Errorf("%v: %s", err, b)
	}
	*s = LegacyStatus(v)
	return nil
}

// String turns a status into a string
func (s LegacyStatus) String() string {
	ss := []string(s)
	if len(ss) == 0 {
		return ""
	}

	return "[" + strings.Join(ss, ",") + "]"
}
