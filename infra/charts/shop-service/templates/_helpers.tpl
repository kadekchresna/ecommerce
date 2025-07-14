{{- define "shop-service.name" -}}
{{ .Chart.Name }}
{{- end }}

{{- define "shop-service.fullname" -}}
{{ printf "%s-%s" .Release.Name .Chart.Name }}
{{- end }}
