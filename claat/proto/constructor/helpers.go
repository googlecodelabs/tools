// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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

func NewHeading(level int32, contentSlice ...*tutorial.InlineContent) *tutorial.Heading {
	return &tutorial.Heading{
		Level:   level,
		Content: contentSlice,
	}
}

func NewParagraph(contentSlice ...*tutorial.InlineContent) *tutorial.Paragraph {
	return &tutorial.Paragraph{
		Content: contentSlice,
	}
}

// TODO: Implement NewList and its tests
// TODO: Implement NewImage and its tests
// TODO: Implement NewImageBlock and its tests

func NewYTVideo(id string) *tutorial.YoutubeVideo {
	return &tutorial.YoutubeVideo{
		Id: id,
	}
}

func NewCodeBlockPlain(code string) *tutorial.CodeBlock {
	return &tutorial.CodeBlock{
		Code: code,
	}
}

func NewCodeBlockHighlighted(code string) *tutorial.CodeBlock {
	return &tutorial.CodeBlock{
		Code:          code,
		IsHighlighted: true,
	}
}

func NewSurvey(id string, contentSlice ...*tutorial.SurveyQuestion) *tutorial.Survey {
	return &tutorial.Survey{
		Id:      id,
		Content: contentSlice,
	}
}

func NewSurveyQuestion(question string, contentSlice ...string) *tutorial.SurveyQuestion {
	return &tutorial.SurveyQuestion{
		Name:    question,
		Options: contentSlice,
	}
}
