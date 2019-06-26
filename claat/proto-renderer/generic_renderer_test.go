package genrenderer

import (
	"go/build"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/googlecodelabs/tools/claat/proto-renderer/testing-utils"
)

type encapsulatedTest struct {
	in  interface{}
	out interface{}
	ok  bool
}

// Setup variables
var (
	tmplsRltvDir = "src/github.com/googlecodelabs/tools/claat/proto-renderer/testdata/*"
	tmplsAbsDir  = filepath.Join(build.Default.GOPATH, tmplsRltvDir)
	funcMap      = template.FuncMap{
		"returnString": func(i string) string { return i },
	}
	invalidCases = []encapsulatedTest{
		{3, nil, false},
		{nil, nil, false},
		{testingUtils.UnsupportedType{}, nil, false},
	}
)

// Demonstrates behavior of non-namespace-compliant template objects
func TestExecuteTemplateInvalidNamespace(t *testing.T) {
	tmpl := template.New("always-panics-dummy")
	runEncapsulatedTestSet(invalidCases, tmpl, t)

	// These cases are only valid for namepace-compliant templates
	validYetNonCompliantCases := []encapsulatedTest{
		{testingUtils.NewDummyProto(), "dummy", false},
	}
	runEncapsulatedTestSet(validYetNonCompliantCases, tmpl, t)
}

// Demonstrates expected behavior
func TestExecuteTemplateValidNamespace(t *testing.T) {
	tmpl := template.Must(template.New("valid-dummy").Funcs(funcMap).ParseGlob(tmplsAbsDir))
	runEncapsulatedTestSet(invalidCases, tmpl, t)

	validCases := []encapsulatedTest{
		{testingUtils.NewDummyProto(), "dummy", true},
	}
	runEncapsulatedTestSet(validCases, tmpl, t)
}

// Iterator helper for 'runEncapsulatedTest'
func runEncapsulatedTestSet(tcs []encapsulatedTest, tmpl *template.Template, t *testing.T) {
	for _, tc := range tcs {
		runEncapsulatedTest(tc, tmpl, t)
	}
}

// runEncapsulatedTest constrains the scope of panics for all tests in this file,
// otherwise we cannot iterate through consecutive panic-causing test-cases
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
	// never gets below if above panicked
	if tc.out != tmplOut {
		t.Errorf("Expecting:\n\t%#v'\nBut got:\n\t%#v", tc.out, tmplOut)
	}
	// dummy return, using for shared defer scope
	return tmplOut
}
