//go:build !solution

package ciletters

import (
	"strings"
	"text/template"
)

const (
	tpl = `Your pipeline #{{ .Pipeline.ID}} {{ if ne .Pipeline.Status "ok" }}has failed{{ else }}passed{{ end }}!
    Project:      {{ .Project.GroupID }}/{{ .Project.ID }}
    Branch:       🌿 {{ .Branch }}
    Commit:       {{ slice .Commit.Hash 0 8 }} {{ .Commit.Message }}
    CommitAuthor: {{ .Commit.Author }}{{ range $job := .Pipeline.FailedJobs }}
        Stage: {{ $job.Stage }}, Job {{ $job.Name }}{{ range cmdLog $job.RunnerLog }}
            {{ . }}{{ end }}
{{ end }}`
)

func cmdLog(s string) []string {
	logLines := strings.Split(s, "\n")
	if len(logLines) > 9 {
		return logLines[9:]
	}

	return logLines
}

func MakeLetter(n *Notification) (string, error) {
	var sb strings.Builder
	instance, err := template.
		New("email").
		Funcs(
			template.FuncMap{
				"cmdLog": cmdLog,
			}).
		Parse(tpl)

	if err != nil {
		return "", err
	}

	if err = instance.Execute(&sb, n); err != nil {
		return "", err
	}

	return sb.String(), nil
}
