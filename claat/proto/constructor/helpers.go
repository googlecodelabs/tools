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

func NewLink(href string, contentSlice ...*tutorial.StylizedText) *tutorial.Link {
	return &tutorial.Link{
		Href:    href,
		Content: contentSlice,
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

func NewInlineImage(
	source string, alt string, height int32, width int32) *tutorial.InlineImage {
	return &tutorial.InlineImage{
		Source: source,
		Alt:    alt,
		Height: height,
		Width:  width,
	}
}

func NewInlineContentTextPlain(txt *tutorial.StylizedText) *tutorial.InlineContent {
	return &tutorial.InlineContent{
		Content: &tutorial.InlineContent_Text{
			Text: txt,
		},
	}
}

func NewInlineContentTextStrong(txt *tutorial.StylizedText) *tutorial.InlineContent {
	return &tutorial.InlineContent{
		Content: &tutorial.InlineContent_Text{
			Text: txt,
		},
	}
}

func NewInlineContentTextEmphasized(txt *tutorial.StylizedText) *tutorial.InlineContent {
	return &tutorial.InlineContent{
		Content: &tutorial.InlineContent_Text{
			Text: txt,
		},
	}
}

func NewInlineContentTextStrongAndEmphasized(txt *tutorial.StylizedText) *tutorial.InlineContent {
	return &tutorial.InlineContent{
		Content: &tutorial.InlineContent_Text{
			Text: txt
		},
	}
}

func NewInlineContentCode(code *tutorial.InlineCode) *tutorial.InlineContent {
	return &tutorial.InlineContent{
		Content: &tutorial.InlineContent_Code{
			Code: code,
		},
	}
}

func NewInlineContentLink(link *tutorial.Link) *tutorial.InlineContent {
	return &tutorial.InlineContent{
		Content: &tutorial.InlineContent_Link{
			Link: link,
		},
	}
}

func NewInlineContentButton(button *tutorial.Button) *tutorial.InlineContent {
	return &tutorial.InlineContent{
		Content: &tutorial.InlineContent_Button{
			Button: button,
		},
	}
}

func NewInlineContentInlineImage(image *tutorial.InlineImage) *tutorial.InlineContent {
	return &tutorial.InlineContent{
		Content: &tutorial.InlineContent_Image{
			Image: image,
		},
	}
}

func NewHeading(level int32, contentSlice ...*tutorial.InlineContent) *tutorial.Heading {
	return &tutorial.Heading{
		Level:   level,
		Content: contentSlice,
	}
}

func NewImageBlock(image *tutorial.InlineImage) *tutorial.ImageBlock {
	return &tutorial.ImageBlock{
		Image: image,
	}
}

func NewParagraph(contentSlice ...*tutorial.InlineContent) *tutorial.Paragraph {
	return &tutorial.Paragraph{
		Content: contentSlice,
	}
}

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

func NewList(
	variety tutorial.List_ListVariety, style tutorial.List_ListStyle,
	contentSlice ...*tutorial.Paragraph) *tutorial.List {
	return &tutorial.List{
		Variety: variety,
		Style:   style,
		Content: contentSlice,
	}
}
