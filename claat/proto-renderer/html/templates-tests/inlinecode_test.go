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
			InProto: tutorial.InlineCode{},
			Out:     "<code></code>",
			Ok:      true,
		},
		{
			InProto: testingUtils.NewInlineCode("< less-than"),
			Out:     "<code>&lt; less-than</code>",
			Ok:      true,
		},
		{
			InProto: testingUtils.NewInlineCode("> greater-than"),
			Out:     "<code>&gt; greater-than</code>",
			Ok:      true,
		},
		{
			InProto: testingUtils.NewInlineCode("/ backslash"),
			Out:     "<code>/ backslash</code>",
			Ok:      true,
		},
		{
			InProto: testingUtils.NewInlineCode(`\ forwardslash`),
			Out:     `<code>\\ forwardslash</code>`,
			Ok:      true,
		},
		{
			InProto: testingUtils.NewInlineCode("& ampersand"),
			Out:     "<code>&amp; ampersand</code>",
			Ok:      true,
		},
		{
			InProto: testingUtils.NewInlineCode(`" quotation`),
			Out:     "<code>&#34; quotation</code>",
			Ok:      true,
		},
		{
			InProto: testingUtils.NewInlineCode("' apostrophe"),
			Out:     "<code>&#39; apostrophe</code>",
			Ok:      true,
		},
		{
			InProto: testingUtils.NewInlineCode("{ Αα Ββ Γγ Δδ Εε Ϝϝ Ζζ Ηη Θθ Ιι Κκ Λλ Μμ Νν Ξξ Οο Ππ Ρρ Σσς Ττ Υυ Φφ Χχ Ψψ Ωω }"),
			Out:     "<code>{ Αα Ββ Γγ Δδ Εε Ϝϝ Ζζ Ηη Θθ Ιι Κκ Λλ Μμ Νν Ξξ Οο Ππ Ρρ Σσς Ττ Υυ Φφ Χχ Ψψ Ωω }</code>",
			Ok:      true,
		},
	}
	testingUtils.CanonicalRenderingTestBatch(html.Render, tests, t)
}
