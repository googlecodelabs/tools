---
id: {{.Meta.ID}}
{{if .Meta.Status}}status: {{.Meta.Status}}{{end}}
{{if .Meta.Summary}}summary: {{.Meta.Summary}}{{end}}
{{if .Meta.Author}}author: {{.Meta.Author}}{{end}}
{{if .Meta.Categories}}categories: {{commaSep .Meta.Categories}}{{end}}
{{if .Meta.Tags}}tags: {{commaSep .Meta.Tags}}{{end}}
{{if .Meta.Feedback}}feedback link: {{.Meta.Feedback}}{{end}}
{{if .Meta.GA}}analytics account: {{.Meta.GA}}{{end}}

---

# {{.Meta.Title}}

{{if .Meta.Feedback}}[Codelab Feedback]({{.Meta.Feedback}}){{end}}

{{range .Steps}}{{if matchEnv .Tags $.Env}}
## {{.Title}}
{{if .Duration}}Duration: {{.Duration.Minutes}}{{end}}
{{.Content | renderMD $.Env}}
{{end}}{{end}}
