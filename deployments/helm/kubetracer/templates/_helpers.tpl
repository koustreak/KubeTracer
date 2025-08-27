{{- define "kubetracer.name" -}}
kubetracer
{{- end -}}

{{- define "kubetracer.fullname" -}}
kubetracer
{{- end -}}

{{- define "kubetracer.chart" -}}
kubetracer-0.1.0
{{- end -}}

{{- define "kubetracer.labels" -}}
app.kubernetes.io/name: kubetracer
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: "1.0.0"
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{- define "kubetracer.selectorLabels" -}}
app.kubernetes.io/name: kubetracer
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "kubetracer.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
kubetracer
{{- else -}}
default
{{- end -}}
{{- end }}
