package md

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/googlecodelabs/tools/claat/types"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestParse(t *testing.T) {
	type output struct {
		codelab  *types.Codelab
		metadata map[string]string
	}
	tests := []struct {
		name string
		in   string
		out  output
		ok   bool
	}{
		{
			name: "Simple",
			in: `this is some text before the title

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
baz`,
			out: output{
				codelab: &types.Codelab{
					Meta: types.Meta{
						Title: "Title Of Codelab",
						Extra: map[string]string{},
					},
					Steps: []*types.Step{
						&types.Step{
							Title: "First Step",
						},
						&types.Step{
							Title: "Second Step",
						},
					},
				},
				metadata: map[string]string{"original publication date": "2019-09-23", "events": "IO2019"},
			},
			ok: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			codelab, metadata, err := Parse(strings.NewReader(tc.in))

			if err == nil && !tc.ok {
				t.Errorf("Parse(string) = %+v, %+v, want err", codelab, metadata)
				return
			}
			if err != nil && tc.ok {
				t.Errorf("Parse(string) = %+v, want %+v, %+v", err, tc.out.codelab, tc.out.metadata)
				return
			}

			// Steps will never be equal because DeepEqual compares addresses.
			var gotSteps, wantSteps []*types.Step
			gotSteps, codelab.Steps = codelab.Steps, []*types.Step{}
			wantSteps, tc.out.codelab.Steps = tc.out.codelab.Steps, []*types.Step{}

			if tc.ok {
				if !reflect.DeepEqual(codelab, tc.out.codelab) || !reflect.DeepEqual(metadata, tc.out.metadata) {
					t.Errorf("Parse(string) = %#v, %+v, want %#v, %+v", codelab, metadata, tc.out.codelab, tc.out.metadata)
					return
				}

				if len(gotSteps) != len(wantSteps) {
					t.Errorf("Parse(string) got %d steps, want %d", len(gotSteps), len(wantSteps))
					return
				}
				for i, v := range gotSteps {
					// We don't care about content yet.
					// TODO pay attention to content equality
					v.Content = nil
					if !reflect.DeepEqual(*v, *wantSteps[i]) {
						t.Errorf("Parse(string).codelab.steps[%d] = %+v, want %+v", i, *v, *wantSteps[i])
					}
				}
			}
		})
	}
}

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
