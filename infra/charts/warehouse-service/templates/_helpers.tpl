{{- define "warehouse-service.name" -}}
{{ .Chart.Name }}
{{- end }}

{{- define "warehouse-service.fullname" -}}
{{ printf "%s-%s" .Release.Name .Chart.Name }}
{{- end }}
