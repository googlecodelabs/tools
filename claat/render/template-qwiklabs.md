# {{.Meta.Title}}

{{if .Meta.Feedback}}[Codelab Feedback]({{.Meta.Feedback}}){{end}}

{{range .Steps}}{{if matchEnv .Tags $.Env}}
## {{.Title}}

{{if .Duration}}*Duration is {{.Duration.Minutes}} min*{{end}}
{{.Content | renderQwiklabs $.Env}}
{{end}}{{end}}
