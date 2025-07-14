{{- define "user-service.name" -}}
{{ .Chart.Name }}
{{- end }}

{{- define "user-service.fullname" -}}
{{ printf "%s-%s" .Release.Name .Chart.Name }}
{{- end }}
