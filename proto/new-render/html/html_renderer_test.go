package codelab_renderer

import (
	"testing"
)

func testHtmlBatch(tests []TestingBatch, t *testing.T) {
	for _, test := range tests {
		r := string(test)
		if test.o != r {
			t.Errorf("Expecting:\n\t'%s', but got \n\t'%s'", test.o, r)
			continue
		}
	}
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

	tests := []TestingBatch{
		{&p, "hello!"},
		{&b, "<strong>hello!</strong>"},
		{&i, "<em>hello!</em>"},
		{&b_i, "<strong><em>hello!</em></strong>"},
	}
	testHtmlBatch(tests, t)
}

func TestRenderHtmlHeading(t *testing.T) {
	h2 := &Heading{
		Level: 2,
		Text:  "<script>some _very_ bad code!</script>",
	}
	h3 := &Heading{
		Level: 3,
		Text:  "D@ ?òü ǝ$çâpæ? {+_~}! ^<^ |*_*| {&]",
	}
	h4 := &Heading{
		Level: 4,
		Text:  "**__Markdown not ![esca](ped)__**",
	}

	tests := []TestingBatch{
		{&h2, "<h2>&lt;script&gt;some _very_ bad code!&lt;/script&gt;</h2>"},
		{&h3, "<h3>D@ ?òü ǝ$çâpæ? {&#43;_~}! ^&lt;^ |*_*| {&amp;]</h3>"},
		{&h4, "<h4>**__Markdown not ![esca](ped)__**</h4>"},
	}
	testHtmlBatch(tests, t)
}

func TestRenderHtmlLink(t *testing.T) {
	// this needs to change...
	// either change it all to plain text templates, or figure out where we need to escape html...
	l := &Link{
		Href: "https://www.youtube.com/watch?v=lyRPyRKHO8M&list=PLOU2XLYxmsILVTiOlMJdo7RQS55jYhsMi",
		Text: &StylizedText{
			Text: "Google I/O 2019 All Sessions",
			IsBold: true,
		},
	}

	tests := []TestingBatch{
		{&l, `<a href"=https://www.youtube.com/watch?v=lyRPyRKHO8M&list=PLOU2XLYxmsILVTiOlMJdo7RQS55jYhsMi target="_blank">Google I/O 2019 All Sessions</a>`},
	}
	testHtmlBatch(tests, t)
}

func TestRenderHtmlLink(t *testing.T) {
	
	// tests := []TestingBatch{
	// 	{&l, `<a href"=https://www.youtube.com/watch?v=lyRPyRKHO8M&list=PLOU2XLYxmsILVTiOlMJdo7RQS55jYhsMi target="_blank">Google I/O 2019 All Sessions</a>`},
	// }
	// testHtmlBatch(tests, t)
}
