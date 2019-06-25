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
