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

type sampleProtoTemplate struct {
	Value interface{}
}

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
		{newSampleProtoTemplate(3), "3", false},
		{newSampleProtoTemplate(nil), "", true},
		{newSampleProtoTemplate("not-valid"), "", true},
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
			t.Errorf("\nFor:\n\t%#v\nPanic occured:\n\t%s\t(false negative)", test, err)
		}

		if err == nil && !test.ok {
			t.Errorf("\nPanic did not occur for:\n\t%#v\n\t(false positive)", test)
		}
	}(test)

	tmplOut := ExecuteTemplate(test.in, tmpl)
	if test.out != tmplOut {
		t.Errorf("Expecting:\n\t'%s'\nBut got:\n\t'%s'", test.out, tmplOut)
	}
}
