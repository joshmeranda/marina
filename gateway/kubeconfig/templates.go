package kubeconfig

import "html/template"

const (
	tokenTemplateText = `apiVersion: v1
kind: Config
clusters:
{{- range .Nodes}}
- name: "{{.ClusterName}}"
  cluster:
    server: "{{.Server}}"
{{- if ne .Cert "" }}
    certificate-authority-data: "{{.Cert}}"
{{- end }}
{{- end}}
users:
- name: "{{.User}}"
  user:
    token: "{{.Token}}"
contexts:
{{- range .Nodes}}
- name: "{{.ClusterName}}"
  context:
    user: "{{.User}}"
    cluster: "{{.ClusterName}}"
{{- end}}

current-context: "{{.ClusterName}}"
`
)

var (
	tokenTemplate = template.Must(template.New("tokenTemplate").Parse(tokenTemplateText))
)
