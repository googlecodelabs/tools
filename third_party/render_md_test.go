package devrel_tutorial

import (
	"testing"
	"github.com/googlecodelabs/tools/claat/util"
)

func TestRenderHeadingMd(t *testing.T) {
	h2 := &Heading{
    Level: 2,
    Text: "<script>some _very_ bad code!;</script>",
	}
	h3 := &Heading{
    Level: 3,
    Text: "D@ ?òü ǝ$çâpæ? {+_~}! ^<^ |*_*| {&]",
	}
	h4 := &Heading{
    Level: 4,
    Text: "**__Extra Markdown not ![pro](cessed)__**",
	}
	tests := []*util.TestingBatch {
		{h2.Md(), "#### <script>some _very_ bad code!;</script>"},
		{h3.Md(), "##### D@ ?òü ǝ$çâpæ? {+_~}! ^<^ |*_*| {&]"},
		{h4.Md(), "###### **__Extra Markdown not ![pro](cessed)__**"},
	}
	util.TestBatch(tests, t)
}

func TestRenderHtmlStylizedText(t *testing.T) {
	p := &StylizedText{
		Text: "hello!",
	}
	b := &StylizedText{
		Text:   "hello!",
		IsBold: true,
	}
	i := &StylizedText{
		Text:         "hello!",
		IsEmphasized: true,
	}
	b_i := &StylizedText{
		Text:         "hello!",
		IsBold:       true,
		IsEmphasized: true,
	}
	tests := []*util.TestingBatch{
		{p.Md(), "hello!"},
		{b.Md(), "**hello!**"},
		{i.Md(), "__hello!__"},
		{b_i.Md(), "**__hello!__**"},
	}
	util.TestBatch(tests, t)
}
