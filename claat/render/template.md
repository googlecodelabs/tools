---
{{metaHeaderYaml .Meta}}
---

# {{.Meta.Title}}

{{if .Meta.Feedback}}[Codelab Feedback]({{.Meta.Feedback}}){{end}}

{{range .Steps}}{{if matchEnv .Tags $.Env}}
## {{.Title}}
{{if .Duration}}Duration: {{durationStr .Duration}}{{end}}
{{.Content | renderMD $.Context}}
{{end}}{{end}}
