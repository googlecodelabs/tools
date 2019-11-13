package md

import "testing"

type testCase struct {
	s  string
	v  string
	sp string
}

func TestSplitSpaceLeft(t *testing.T) {
	testCases := []testCase{
		testCase{"basic test case", "basic test case", ""},
		testCase{"  some space  ", "some space  ", "  "},
		testCase{"\t\n\rsome less common spaces", "some less common spaces", "\t\n\r"},
		testCase{" 🙂 a smiley space ", "🙂 a smiley space ", " "},
		testCase{"       ", "", "       "},
		testCase{"", "", ""},
		testCase{"       x", "x", "       "},
	}

	for _, tc := range testCases {
		v, sp := splitSpaceLeft(tc.s)
		if v != tc.v {
			t.Errorf("v=%#v, expected %#v, (s=%#v)", v, tc.v, tc.s)
		}
		if sp != tc.sp {
			t.Errorf("sp=%#v, expected %#v, (s=%#v)", sp, tc.sp, tc.s)
		}
	}
}

func TestSplitSpaceRight(t *testing.T) {
	testCases := []testCase{
		testCase{"basic test case", "basic test case", ""},
		testCase{"  some space  ", "  some space", "  "},
		testCase{"some less common spaces\t\n\r", "some less common spaces", "\t\n\r"},
		testCase{" 🙂 a smiley space ", " 🙂 a smiley space", " "},
		testCase{"       ", "", "       "},
		testCase{"", "", ""},
		testCase{"x       ", "x", "       "},
	}

	for _, tc := range testCases {
		v, sp := splitSpaceRight(tc.s)
		if v != tc.v {
			t.Errorf("v=%#v, expected %#v, (s=%#v)", v, tc.v, tc.s)
		}
		if sp != tc.sp {
			t.Errorf("sp=%#v, expected %#v, (s=%#v)", sp, tc.sp, tc.s)
		}
	}
}
