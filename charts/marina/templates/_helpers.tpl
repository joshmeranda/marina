{{/*
Expand the name of the chart.
*/}}
{{- define "marina.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "marina.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "marina.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "marina.labels" -}}
helm.sh/chart: {{ include "marina.chart" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Gateway Labels
*/}}
{{- define "marina.gateway.labels" -}}
{{ include "marina.labels" . }}
app.kubernetes.io/component: "gateway"
{{- end }}

{{/*
Gateway Selector Labels
*/}}
{{- define "marina.gateway.selectorLabels" -}}
app.kubernetes.io/name: {{ include "marina.name" . }}-gateway
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Operator Labels
*/}}
{{- define "marina.operator.labels" -}}
{{ include "marina.labels" . }}
app.kubernetes.io/component: "operator"
{{- end }}

{{/*
Gateway Selector Labels
*/}}
{{- define "marina.operator.selectorLabels" -}}
app.kubernetes.io/name: {{ include "marina.name" . }}-operator
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "marina.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "marina.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Conform an image ref to the string format expected by the gateway.
*/}}
{{- define "toReference" -}}
{{- $ref := "" }}
{{- if .registry -}}
{{ $ref = printf "%s/" .registry }}
{{- end -}}
{{- if .repository -}}
{{ $ref = printf "%s%s" $ref .repository }}
{{- end -}}
{{- if .tag -}}
{{ $ref = printf "%s:%s" $ref .tag }}
{{- end -}}
{{- if .sha -}}
{{ $ref = printf "%s@sha256:%s" $ref .sha }}
{{- end -}}
{{- $ref | quote -}}
{{- end -}}
