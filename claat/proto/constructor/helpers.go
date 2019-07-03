package protoconstructors

import (
	"github.com/googlecodelabs/tools/third_party"
)

// Helper constructor functions

func NewStylizedTextPlain(txt string) *tutorial.StylizedText {
	return &tutorial.StylizedText{
		Text: txt,
	}
}

func NewStylizedTextStrong(txt string) *tutorial.StylizedText {
	return &tutorial.StylizedText{
		Text:     txt,
		IsStrong: true,
	}
}

func NewStylizedTextEmphasized(txt string) *tutorial.StylizedText {
	return &tutorial.StylizedText{
		Text:         txt,
		IsEmphasized: true,
	}
}

func NewStylizedTextStrongAndEmphasized(txt string) *tutorial.StylizedText {
	return &tutorial.StylizedText{
		Text:         txt,
		IsStrong:     true,
		IsEmphasized: true,
	}
}

func NewInlineCode(txt string) *tutorial.InlineCode {
	return &tutorial.InlineCode{
		Code: txt,
	}
}

func NewInlineContentCode(txt string) *tutorial.InlineContent {
	return &tutorial.InlineContent{
		Content: &tutorial.InlineContent_Code{
			Code: &tutorial.InlineCode{
				Code: txt,
			},
		},
	}
}

func NewLink(href string, contentSlice ...*tutorial.StylizedText) *tutorial.Link {
	return &tutorial.Link{
		Href:    href,
		Content: contentSlice,
	}
}

func NewInlineContentLink(link *tutorial.Link) *tutorial.InlineContent {
	return &tutorial.InlineContent{
		Content: &tutorial.InlineContent_Link{
			Link: link,
		},
	}
}

func NewButtonPlain(link *tutorial.Link) *tutorial.Button {
	return &tutorial.Button{
		Link: link,
	}
}

func NewButtonDownload(link *tutorial.Link) *tutorial.Button {
	return &tutorial.Button{
		Link:           link,
		IsDownloadable: true,
	}
}

// TODO: Add to InlineContent tests
func NewInlineContentButton(button *tutorial.Button) *tutorial.InlineContent {
	return &tutorial.InlineContent{
		Content: &tutorial.InlineContent_Button{
			Button: button,
		},
	}
}

// TODO: Add to InlineContent tests
func NewInlineContentImage(image *tutorial.Image) *tutorial.InlineContent {
	return &tutorial.InlineContent{
		Content: &tutorial.InlineContent_Image{
			Image: image,
		},
	}
}

func NewInlineContentTextPlain(txt string) *tutorial.InlineContent {
	return &tutorial.InlineContent{
		Content: &tutorial.InlineContent_Text{
			Text: &tutorial.StylizedText{
				Text: txt,
			},
		},
	}
}

func NewInlineContentTextStrong(txt string) *tutorial.InlineContent {
	return &tutorial.InlineContent{
		Content: &tutorial.InlineContent_Text{
			Text: &tutorial.StylizedText{
				Text:     txt,
				IsStrong: true,
			},
		},
	}
}

func NewInlineContentTextEmphasized(txt string) *tutorial.InlineContent {
	return &tutorial.InlineContent{
		Content: &tutorial.InlineContent_Text{
			Text: &tutorial.StylizedText{
				Text:         txt,
				IsEmphasized: true,
			},
		},
	}
}

func NewInlineContentTextStrongAndEmphasized(txt string) *tutorial.InlineContent {
	return &tutorial.InlineContent{
		Content: &tutorial.InlineContent_Text{
			Text: &tutorial.StylizedText{
				Text:         txt,
				IsStrong:     true,
				IsEmphasized: true,
			},
		},
	}
}

func NewParagraph(contentSlice ...*tutorial.InlineContent) *tutorial.Paragraph {
	return &tutorial.Paragraph{
		Content: contentSlice,
	}
}

// TODO: Implement NewList and its tests
// TODO: Implement NewImage and its tests
