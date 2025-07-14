{{- define "product-service.name" -}}
{{ .Chart.Name }}
{{- end }}

{{- define "product-service.fullname" -}}
{{ printf "%s-%s" .Release.Name .Chart.Name }}
{{- end }}
