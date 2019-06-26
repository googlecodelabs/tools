package testingUtils

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
