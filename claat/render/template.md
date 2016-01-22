# {{.Meta.Title}}

[//]: # (This is the default template for 'md' output format of the tool)
[//]: # (This file is distributed under Apache-2 license; see LICENSE file at the root of this repo)

{{if .Meta.Feedback}}[Codelab Feedback]({{.Meta.Feedback}}){{end}}

{{range .Steps}}
## {{.Title}}

{{if .Duration}}Duration is {{.Duration}}{{end}}

{{range .Content.Nodes}}
{{.}}
{{end}}
{{end}}
