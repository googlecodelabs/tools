# {{.Meta.Title}}

{{range .Steps}}{{if matchEnv .Tags $.Env}}
## {{.Title}}
{{if .Duration}}Duration: {{durationStr .Duration}}{{end}}
{{.Content | renderQwiklabsMD $.Context}}
{{end}}{{end}}
