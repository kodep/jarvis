{{/* vim: set filetype=mustache: */}}
{{/*
Return the app name as defined in values.
*/}}
{{- define "jarvis.app.name" -}}
{{- .Values.app.name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Return the database name as defined in values.
*/}}
{{- define "jarvis.db.name" -}}
{{- .Values.database.name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "jarvis.labels" -}}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
App selector labels
*/}}
{{- define "jarvis.app.selectorLabels" -}}
app: {{ include "jarvis.app.name" . }}
{{- end }}

{{/*
Database selector labels
*/}}
{{- define "jarvis.db.selectorLabels" -}}
app: {{ include "jarvis.db.name" . }}
{{- end }}

{{/*
Service name for app
*/}}
{{- define "jarvis.app.service.name" -}}
{{- printf "%s-service" (include "jarvis.app.name" .) | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
PVC name for database
*/}}
{{- define "jarvis.db.pvc.name" -}}
{{- printf "%s-pvc" (include "jarvis.db.name" .) | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Ingress name
*/}}
{{- define "jarvis.ingress.name" -}}
{{- printf "%s-ingress" (include "jarvis.app.name" .) | trunc 63 | trimSuffix "-" }}
{{- end }}%    
