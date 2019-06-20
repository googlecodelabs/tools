package genrenderer

import (
	"go/build"
	"path/filepath"
	"testing"
	"text/template"
)

type encapsulatedTest struct {
	in  sampleProtoTemplate
	out string
	ok  bool
}

// TODO: update to be proto dependent on next PR
type sampleProtoTemplate struct {
	Value interface{}
}

// TODO: update to be proto dependent on next PR
func newSampleProtoTemplate(el interface{}) sampleProtoTemplate {
	return sampleProtoTemplate{Value: el}
}

func TestExecuteTemplate(t *testing.T) {
	tmplsRltvDir := "src/github.com/googlecodelabs/tools/claat/proto-renderer/testdata/*"
	tmplsAbsDir := filepath.Join(build.Default.GOPATH, tmplsRltvDir)
	funcMap := template.FuncMap{
		"returnInt": func(i int) int { return i },
	}
	tmpl := template.Must(template.New("dummy").Funcs(funcMap).ParseGlob(tmplsAbsDir))

	tests := []encapsulatedTest{
		{newSampleProtoTemplate(3), "3", true},
		{newSampleProtoTemplate(nil), "", false},
		{newSampleProtoTemplate("not-valid"), "", false},
	}

	for _, tc := range tests {
		runEncapsulatedTest(tc, tmpl, t)
	}
}

func runEncapsulatedTest(test encapsulatedTest, tmpl *template.Template, t *testing.T) {
	// Check wheather template failed to render by checking for panic
	defer func(test encapsulatedTest) {
		err := recover()
		if err != nil && test.ok {
			t.Errorf("\nExecuteTemplate(\n\t%#v,\n\t%v,\n) = %#v\nPanic occured:\n\t%#v\n(false negative)", test, tmpl, test.out, err)
		}

		if err == nil && !test.ok {
			t.Errorf("\nExecuteTemplate(\n\t%#v,\n\t%v,\n) = %#v\nWant panic\n(false positive)", test, tmpl, test.out)
		}
	}(test)

	tmplOut := ExecuteTemplate(test.in, tmpl)
	// never gets here if above panicked
	if test.out != tmplOut {
		t.Errorf("Expecting:\n\t'%s'\nBut got:\n\t'%s'", test.out, tmplOut)
	}
}
