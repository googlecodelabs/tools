package types

import (
	"testing"
	"time"
	_ "time/tzdata"

	"github.com/google/go-cmp/cmp"
)

func TestContextTimeMarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		in   ContextTime
		out  []byte
	}{
		{
			name: "Epoch",
			in:   ContextTime(time.Unix(0, 0).UTC()),
			out:  []byte(`"1970-01-01T00:00:00Z"`),
		},
		{
			name: "WithTimeZone",
			in:   ContextTime(time.Unix(1629497889, 0).In(time.FixedZone("San Francisco (DST)", -7*60*60))),
			out:  []byte(`"2021-08-20T15:18:09-07:00"`),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, err := tc.in.MarshalJSON()
			if err != nil {
				t.Errorf("ContextTime.MarshalJSON() = %+v , want %+q", err, tc.out)
				return
			}
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("ContextTime.MarshalJSON got diff (-want +got): %s", diff)
				return
			}
		})
	}
}

func TestContextTimeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		// ContextTime sets internal fields we don't care about -- easier to compare Unix times.
		out int64
		ok  bool
	}{
		{
			name: "RFC3339",
			in:   []byte(`"2021-08-20T22:35:40Z"`),
			out:  1629498940,
			ok:   true,
		},
		{
			name: "YYYY-MM-DD",
			in:   []byte(`"2021-08-20"`),
			out:  1629417600,
			ok:   true,
		},
		// TODO should wrong quotes be accepted?
		{
			name: "RFC3339NoQuotes",
			in:   []byte(`2021-08-20T22:35:40Z`),
			out:  1629498940,
			ok:   true,
		},
		{
			name: "YYYY-MM-DDNoQuotes",
			in:   []byte(`2021-08-20`),
			out:  1629417600,
			ok:   true,
		},
		{
			name: "RFC3339OnlyOpeningQuote",
			in:   []byte(`"2021-08-20T22:35:40Z`),
			out:  1629498940,
			ok:   true,
		},
		{
			name: "YYYY-MM-DDOnlyOpeningQuote",
			in:   []byte(`"2021-08-20`),
			out:  1629417600,
			ok:   true,
		},
		{
			name: "RFC3339OnlyClosingQuote",
			in:   []byte(`2021-08-20T22:35:40Z"`),
			out:  1629498940,
			ok:   true,
		},
		{
			name: "YYYY-MM-DDOnlyClosingQuote",
			in:   []byte(`2021-08-20"`),
			out:  1629417600,
			ok:   true,
		},
		{
			name: "Invalid",
			in:   []byte("foobar"),
		},
		{
			name: "BrokenRFC3339",
			in:   []byte(`"2021-08-2022:35:40Z"`),
		},
		{
			name: "BrokenYYYY-MM-DD",
			in:   []byte(`"2021-13-20"`),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ct := &ContextTime{}
			err := ct.UnmarshalJSON(tc.in)
			if tc.ok && err != nil {
				t.Errorf("ContextTime.UnmarshalJSON(%+v) got err %+v, want %+v", tc.in, err, tc.out)
				return
			}
			if !tc.ok && err == nil {
				t.Errorf("ContextTime.UnmarshalJSON(%+v) got %+v, want err", tc.in, ct)
				return
			}
			// ContextTime sets internal fields that we don't care about that makes cmp.Diff undesirable for comparison here.
			gotUnixTime := time.Time(*ct).Unix()
			if tc.ok && gotUnixTime != tc.out {
				t.Errorf("ContextTime.UnmarshalJSON(%+v) got time %d, want %d", tc.in, gotUnixTime, tc.out)
				return
			}
		})
	}
}
