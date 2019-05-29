package codelab_renderer

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
		{&h2.Md(), "#### &lt;script&gt;some _very_ bad code!;&lt;/script&gt;"},
		{&h3.Md(), "##### D@ ?òü ǝ$çâpæ? {&#43;_~}! ^&lt;^ |*_*| {&amp;]"},
		{&h4.Md(), "###### **__Extra Markdown not ![pro](cessed)__**"},
	}
	testBatch(tests, t)
}