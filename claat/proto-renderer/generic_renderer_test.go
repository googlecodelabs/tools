package genrenderer

import (
	"go/build"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/googlecodelabs/tools/claat/proto-renderer"
)

type encapsulatedTest struct {
	in  interface{}
	out string
	ok  bool
}

func newDummyProto(in string) {
	return &devrel_tutorial.StylizedText{
		Text: in,
	}
}

// Setup variables
var (
	tmplsRltvDir = "src/github.com/googlecodelabs/tools/claat/proto-renderer/testdata/*"
	tmplsAbsDir  = filepath.Join(build.Default.GOPATH, tmplsRltvDir)
	funcMap      = template.FuncMap{
		"returnInt": func(i int) int { return i },
	}
)

func TestExecuteTemplate(t *testing.T) {

	tests := []encapsulatedTest{
		// invalid inputs
		{3, nil, false},
		{nil, nil, false},
		{encapsulatedTest{}, nil, false},
		// valid inputs
		{newDummyProto("hello"), "hello", true},
	}
}
func TestExecuteTemplate(t *testing.T) {

	// need bad  templatetests...
	tmpl := template.Must(template.New("valid-dummy").Funcs(funcMap).ParseGlob(tmplsAbsDir))
	for _, tc := range tests {
		runEncapsulatedTest(tc, tmpl, t)
	}
}

// runEncapsulatedTest constrains the scope of panics, else we cannot iterate
// through consecutive panic-causing test-cases
func runEncapsulatedTest(tc encapsulatedTest, tmpl *template.Template, t *testing.T) (tmplOut string) {
	// Check whether template failed to render by checking for panic
	defer func(tc encapsulatedTest) {
		err := recover()
		if err != nil && tc.ok {
			t.Errorf("\nExecuteTemplate(\n\t%#v,\n\t%v,\n)\nPanic: %v(false negative)\nWant: %#v", tc.in, tmpl, err, tc.out)
		}

		if err == nil && !tc.ok {
			t.Errorf("\nExecuteTemplate(\n\t%#v,\n\t%v,\n) = %#v\nWant Panic\n(false positive)", tc.in, tmpl, tmplOut)
		}
	}(tc)

	tmplOut = ExecuteTemplate(tc.in, tmpl)
	// never gets here if above panicked
	if tc.out != tmplOut {
		t.Errorf("Expecting:\n\t'%s'\nBut got:\n\t'%s'", tc.out, tmplOut)
	}
	// dummy return, using for shared defer scope
	return tmplOut
}
