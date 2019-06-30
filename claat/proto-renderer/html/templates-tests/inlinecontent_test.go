package htmltests

import (
	"go/build"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/googlecodelabs/tools/claat/proto-renderer/html"
	"github.com/googlecodelabs/tools/claat/proto-renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

func TestRenderInlineContentTemplateLinkFromFile(t *testing.T) {
	linkFileRelDir := "src/github.com/googlecodelabs/tools/claat/proto-renderer/html/templates-tests/testdata/InlineContent/google_weather.html"
	linkFileAbsDir := filepath.Join(build.Default.GOPATH, linkFileRelDir)
	weatherLinkBytes, err := ioutil.ReadFile(linkFileAbsDir)
	if err != nil {
		t.Errorf("Reading %#v outputted %#v", linkFileAbsDir, err)
		return
	}
	weatherLinkOutput := string(weatherLinkBytes[:])

	linkProto := testingutils.NewLink(
		"https://www.google.com/search?q=weather+in+nyc",
		testingutils.NewStylizedTextPlain("hey google,"),
		testingutils.NewStylizedTextStrong(" how's the"),
		testingutils.NewStylizedTextEmphasized(" weather in "),
		testingutils.NewStylizedTextStrongAndEmphasized("NYC today?"),
	)

	canonicalTests := []*testingutils.CanonicalRenderingBatch{
		{
			InProto: linkProto,
			Out:     weatherLinkOutput,
			Ok:      true,
		},
	}
	testingutils.CanonicalRenderTestBatch(html.Render, canonicalTests, t)

	identityTests := []*testingutils.RendererIdendityBatch{
		{
			InProto:  testingutils.NewInlineContentLink(linkProto),
			OutProto: linkProto,
			Out:      weatherLinkOutput,
			Ok:       true,
		},
	}
	testingutils.RenderingIdendityTestBatch(html.Render, identityTests, t)
}

func TestRenderInlineContentTemplateIdentity(t *testing.T) {
	tests := []*testingutils.RendererIdendityBatch{
		{
			InProto:  testingutils.NewInlineContentTextPlain(`<script>alert("you've been hacked!");</script>!`),
			OutProto: testingutils.NewStylizedTextPlain(`<script>alert("you've been hacked!");</script>!`),
			Out:      `&lt;script&gt;alert(&#34;you&#39;ve been hacked!&#34;);&lt;/script&gt;!`,
			Ok:       true,
		},
		{
			InProto:  testingutils.NewInlineContentCode(`<script>alert("you've been hacked!");</script>!`),
			OutProto: testingutils.NewInlineCode(`<script>alert("you've been hacked!");</script>!`),
			Out:      `<code>&lt;script&gt;alert(&#34;you&#39;ve been hacked!&#34;);&lt;/script&gt;!</code>`,
			Ok:       true,
		},
	}
	testingutils.RenderingIdendityTestBatch(html.Render, tests, t)
}

func TestRenderInlineContentTemplate(t *testing.T) {
	tests := []*testingutils.CanonicalRenderingBatch{
		{
			InProto: &tutorial.InlineContent{},
			Out:     "",
			Ok:      false,
		},
		{
			InProto: testingutils.NewInlineContentTextPlain(""),
			Out:     "",
			Ok:      true,
		},
		{
			InProto: testingutils.NewInlineContentTextStrong("strong"),
			Out:     "<strong>strong</strong>",
			Ok:      true,
		},
		{
			InProto: testingutils.NewInlineContentTextEmphasized("emphasized"),
			Out:     "<em>emphasized</em>",
			Ok:      true,
		},
		{
			InProto: testingutils.NewInlineContentTextStrongAndEmphasized("strong & emphasized"),
			Out:     "<strong><em>strong &amp; emphasized</em></strong>",
			Ok:      true,
		},
		{
			InProto: testingutils.NewInlineContentCode(`~!@#$%^&*()_+-=[]{}\|'";:/?.><,`),
			Out:     `<code>~!@#$%^&amp;*()_+-=[]{}\\|&#39;&#34;;:/?.&gt;&lt;,</code>`,
			Ok:      true,
		},
	}
	testingutils.CanonicalRenderTestBatch(html.Render, tests, t)
}
