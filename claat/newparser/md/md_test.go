package md

import (
	"bytes"
	"testing"

	"github.com/googlecodelabs/tools/claat/types"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestParse(t *testing.T) {
	tests := []struct {
		in  string
		out struct {
			codelab  *types.Codelab
			metadata map[string]string
		}
		ok bool
	}{}
}

const in = `this is some text before the title

# Title Of Codelab

original publication date: 2019-09-23
events: IO2019

some other text before the first step

## First Step

blah

blah

blah

## Second Step

foo
bar
baz`

func TestSeek(t *testing.T) {
	testData := []byte(`<body>
	<p class="pineapple">some text!</p>
	<p class="strawberry">
		<blink class="peach">some nested text!</blink>
	</p>
</body>`)

	tests := []struct {
		in  atom.Atom
		out string
		ok  bool
	}{
		// one is present
		{
			in:  atom.Blink,
			out: "peach",
			ok:  true,
		},
		// more than one is present
		{
			in:  atom.P,
			out: "pineapple",
			ok:  true,
		},
		// zero are present
		{
			in: atom.Marquee,
		},
	}
	for _, tc := range tests {
		testDataNode, err := html.Parse(bytes.NewReader(testData))
		if err != nil {
			t.Fatal("test data could not be parsed")
		}

		out, err := seek(testDataNode, tc.in)

		if err == nil && !tc.ok {
			t.Errorf("seek(testdata, %q) = %+v, want error", tc.in, out)
		}
		if err != nil && tc.ok {
			t.Errorf("seek(testdata, %q) = %+v, want %+v", tc.in, err, tc.out)
		}

		if tc.ok {
			// It's difficult to test html.Node equality exactly because of all the pointers.
			// Instead we check for classes in addition to atoms to confirm correctness.
			var outClass string
			for _, a := range out.Attr {
				if a.Key == "class" {
					outClass = a.Val
				}
			}
			if out.DataAtom != tc.in || outClass != tc.out {
				t.Errorf("seek(testdata, %q) = %+v, want %+v", tc.in, out, tc.out)
			}
		}
	}
}
