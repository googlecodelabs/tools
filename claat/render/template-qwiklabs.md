---
{{metaHeaderYaml .Meta}}
---

# {{.Meta.Title}}

My **Local** MD Template!


{{if .Meta.Feedback}}[Codelab Feedback]({{.Meta.Feedback}}){{end}}

{{range .Steps}}{{if matchEnv .Tags $.Env}}
## {{.Title}}
{{if .Duration}}Duration: {{durationStr .Duration}}{{end}}
{{.Content | renderQwiklabsMD $.Context}}
{{end}}{{end}}
