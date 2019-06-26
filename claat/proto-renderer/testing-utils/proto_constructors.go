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

func NewStylizedTextEmphazied(txt string) *tutorial.StylizedText {
	return &tutorial.StylizedText{
		Text:         txt,
		IsEmphasized: true,
	}
}

func NewStylizedTextStrongAndEmphazied(txt string) *tutorial.StylizedText {
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
			&tutorial.InlineCode{
				Code: txt,
			},
		},
	}
}

func NewInlineContentTextPlain(txt string) *tutorial.InlineContent {
	return &tutorial.InlineContent{
		Content: &tutorial.InlineContent_Text{
			&tutorial.StylizedText{
				Text: txt,
			},
		},
	}
}

func NewInlineContentTextStrong(txt string) *tutorial.InlineContent {
	return &tutorial.InlineContent{
		Content: &tutorial.InlineContent_Text{
			&tutorial.StylizedText{
				Text:     txt,
				IsStrong: true,
			},
		},
	}
}

func NewInlineContentTextEmphazied(txt string) *tutorial.InlineContent {
	return &tutorial.InlineContent{
		Content: &tutorial.InlineContent_Text{
			&tutorial.StylizedText{
				Text:         txt,
				IsEmphasized: true,
			},
		},
	}
}

func NewInlineContentTextStrongAndEmphazied(txt string) *tutorial.InlineContent {
	return &tutorial.InlineContent{
		Content: &tutorial.InlineContent_Text{
			&tutorial.StylizedText{
				Text:         txt,
				IsStrong:     true,
				IsEmphasized: true,
			},
		},
	}
}

func NewTextBlock(contentSlice ...*tutorial.InlineContent) *tutorial.TextBlock {
	return &tutorial.TextBlock{
		Content: contentSlice,
	}
}
