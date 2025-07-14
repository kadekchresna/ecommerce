{{- define "order-service.name" -}}
{{ .Chart.Name }}
{{- end }}

{{- define "order-service.fullname" -}}
{{ printf "%s-%s" .Release.Name .Chart.Name }}
{{- end }}
