{{- define "chartName" -}}
"{{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}"
{{- end -}}

{{- define "fullName" -}}
"{{ .Chart.Name }}-server"
{{- end -}}