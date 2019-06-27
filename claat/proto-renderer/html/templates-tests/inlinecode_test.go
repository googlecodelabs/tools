package htmltests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto-renderer/html"
	"github.com/googlecodelabs/tools/claat/proto-renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

func TestRenderInlineCodeTemplateEscaping(t *testing.T) {
	tests := []*testingutils.CanonicalRenderingBatch{
		{
			InProto: tutorial.InlineCode{},
			Out:     "<code></code>",
			Ok:      true,
		},
		{
			InProto: testingutils.NewInlineCode("< less-than"),
			Out:     "<code>&lt; less-than</code>",
			Ok:      true,
		},
		{
			InProto: testingutils.NewInlineCode("> greater-than"),
			Out:     "<code>&gt; greater-than</code>",
			Ok:      true,
		},
		{
			InProto: testingutils.NewInlineCode("/ backslash"),
			Out:     "<code>/ backslash</code>",
			Ok:      true,
		},
		{
			InProto: testingutils.NewInlineCode(`\ forwardslash`),
			Out:     `<code>\\ forwardslash</code>`,
			Ok:      true,
		},
		{
			InProto: testingutils.NewInlineCode("& ampersand"),
			Out:     "<code>&amp; ampersand</code>",
			Ok:      true,
		},
		{
			InProto: testingutils.NewInlineCode(`" quotation`),
			Out:     "<code>&#34; quotation</code>",
			Ok:      true,
		},
		{
			InProto: testingutils.NewInlineCode("' apostrophe"),
			Out:     "<code>&#39; apostrophe</code>",
			Ok:      true,
		},
		{
			InProto: testingutils.NewInlineCode("{ Αα Ββ Γγ Δδ Εε Ϝϝ Ζζ Ηη Θθ Ιι Κκ Λλ Μμ Νν Ξξ Οο Ππ Ρρ Σσς Ττ Υυ Φφ Χχ Ψψ Ωω }"),
			Out:     "<code>{ Αα Ββ Γγ Δδ Εε Ϝϝ Ζζ Ηη Θθ Ιι Κκ Λλ Μμ Νν Ξξ Οο Ππ Ρρ Σσς Ττ Υυ Φφ Χχ Ψψ Ωω }</code>",
			Ok:      true,
		},
	}
	testingutils.CanonicalRenderTestBatch(html.Render, tests, t)
}
