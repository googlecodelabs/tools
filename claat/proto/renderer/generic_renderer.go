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
package genrenderer

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/googlecodelabs/tools/third_party"
)

// ExecuteTemplate returns the evaluated template per passed templating
// namespace, based on the passed tutorial proto type string name
func ExecuteTemplate(el interface{}, t *template.Template) string {
	var w bytes.Buffer
	e := t.ExecuteTemplate(&w, outputFormatTemplateName(el, t), el)
	if e != nil {
		// This method outputs directly to templates. Panicking to surfance errors
		// since we should not handle multiple returns in templates.
		// Errors will be more gracefully handled in output-format packages
		panic(fmt.Sprintf("Templating panic: %s\n", e))
	}
	return w.String()
}

// outputFormatTemplateName concatenates the template name mapping of the passed proto
// with its output package extension
func outputFormatTemplateName(el interface{}, t *template.Template) string {
	return templateName(el) + "." + t.Name()
}

// templateName Maps protos to their type string name
func templateName(el interface{}) string {
	switch el.(type) {
	case *tutorial.StylizedText, tutorial.StylizedText:
		return "StylizedText"
	case *tutorial.InlineCode, tutorial.InlineCode:
		return "InlineCode"
	case *tutorial.InlineContent, tutorial.InlineContent:
		return "InlineContent"
	case *tutorial.Heading, tutorial.Heading:
		return "Heading"
	case *tutorial.Paragraph, tutorial.Paragraph:
		return "Paragraph"
	case *tutorial.Link, tutorial.Link:
		return "Link"
	case *tutorial.Button, tutorial.Button:
		return "Button"
	case *tutorial.List, tutorial.List:
		return "List"
	case *tutorial.InlineImage, tutorial.InlineImage:
		return "InlineImage"
	case *tutorial.ImageBlock, tutorial.ImageBlock:
		return "ImageBlock"
	case *tutorial.YoutubeVideo, tutorial.YoutubeVideo:
		return "YoutubeVideo"
	case *tutorial.CodeBlock, tutorial.CodeBlock:
		return "CodeBlock"
	case *tutorial.SurveyQuestion, tutorial.SurveyQuestion:
		return "SurveyQuestion"
	case *tutorial.Survey, tutorial.Survey:
		return "Survey"
	}
	// This will cause a debug-friendly panic
	return TypeNotSupported("genrenderer.templateName", el)
}
