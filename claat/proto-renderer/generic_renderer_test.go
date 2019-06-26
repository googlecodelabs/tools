package genrenderer

import (
	"go/build"
	"path/filepath"
	"testing"
	"text/template"
)

type encapsulatedTest struct {
	in  *SampleProtoTemplate
	out string
	ok  bool
}

func TestExecuteTemplate(t *testing.T) {
	tmplsRltvDir := "src/github.com/googlecodelabs/tools/claat/proto-renderer/testdata/*"
	tmplsAbsDir := filepath.Join(build.Default.GOPATH, tmplsRltvDir)
	funcMap := template.FuncMap{
		"returnInt": func(i int) int { return i },
	}
	tmpl := template.Must(template.New("dummy").Funcs(funcMap).ParseGlob(tmplsAbsDir))

	tests := []encapsulatedTest{
		{NewSampleProtoTemplate(3), "3", true},
		{NewSampleProtoTemplate(nil), "", false},
		{NewSampleProtoTemplate("not-valid"), "", false},
	}

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
