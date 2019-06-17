package genrenderer

import (
	"go/build"
	"path/filepath"
	"testing"
	"text/template"
)

type encapsulatedTest struct {
	In      sampleProtoTemplate
	Out     string
	WantErr bool
}

type sampleProtoTemplate struct {
	Value interface{}
}

func newSampleProtoTemplate(el interface{}) sampleProtoTemplate {
	return sampleProtoTemplate{Value: el}
}

func runEncapsulatedTest(test encapsulatedTest, tmpl *template.Template, t *testing.T) {
	// Check wheather template failed to render by checking if there was a panic
	defer func(test encapsulatedTest) {
		r := recover()
		if r != nil && !test.WantErr {
			t.Errorf("False positive panic for: %+v", test)
		}

		if r == nil && test.WantErr {
			t.Errorf("False negative panic for: %+v", test)
		}
	}(test)

	tmplOut := ExecuteTemplate(test.In, tmpl)
	if test.Out != tmplOut && !test.WantErr {
		t.Errorf("Expecting:\n\t'%s', but got \n\t'%s'", test.Out, test.In)
	}
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
		{newSampleProtoTemplate("not-valid"), "", true},
	}

	for _, test := range tests {
		runEncapsulatedTest(test, tmpl, t)
	}
}
