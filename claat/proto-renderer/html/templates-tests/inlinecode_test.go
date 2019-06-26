package htmltests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto-renderer/html"
	"github.com/googlecodelabs/tools/claat/proto-renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

func TestRenderInlineCodeTemplateEscaping(t *testing.T) {
	tests := []*testingUtils.RendererTestingBatch{
		{
			tutorial.InlineCode{},
			"<code></code>",
			true,
		},
		{
			testingUtils.NewInlineCode("< less-than"),
			"<code>&lt; less-than</code>",
			true,
		},
		{
			testingUtils.NewInlineCode("> greater-than"),
			"<code>&gt; greater-than</code>",
			true,
		},
		{
			testingUtils.NewInlineCode("/ backslash"),
			"<code>/ backslash</code>",
			true,
		},
		{
			testingUtils.NewInlineCode(`\ forwardslash`),
			`<code>\\ forwardslash</code>`,
			true,
		},
		{
			testingUtils.NewInlineCode("& ampersand"),
			"<code>&amp; ampersand</code>",
			true,
		},
		{
			testingUtils.NewInlineCode(`" quotation`),
			"<code>&#34; quotation</code>",
			true,
		},
		{
			testingUtils.NewInlineCode("' apostrophe"),
			"<code>&#39; apostrophe</code>",
			true,
		},
		{
			testingUtils.NewInlineCode("{ Αα Ββ Γγ Δδ Εε Ϝϝ Ζζ Ηη Θθ Ιι Κκ Λλ Μμ Νν Ξξ Οο Ππ Ρρ Σσς Ττ Υυ Φφ Χχ Ψψ Ωω }"),
			"<code>{ Αα Ββ Γγ Δδ Εε Ϝϝ Ζζ Ηη Θθ Ιι Κκ Λλ Μμ Νν Ξξ Οο Ππ Ρρ Σσς Ττ Υυ Φφ Χχ Ψψ Ωω }</code>",
			true,
		},
	}
	testingUtils.CanonicalRenderingTestBatch(html.Render, tests, t)
}
