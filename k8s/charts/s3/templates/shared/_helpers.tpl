{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to
this (by the DNS naming spec). If release name contains chart name it will
be used as a full name.
*/}}
{{- define "hanzo.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create a truncated component name.
Usage: {{ include "hanzo.componentName" (list . "component-suffix") }}
*/}}
{{- define "hanzo.componentName" -}}
{{- $context := index . 0 -}}
{{- $suffix := index . 1 -}}
{{- if gt (len $suffix) 61 -}}
{{-   fail (printf "Suffix '%s' is too long for componentName helper. Max length is 61." $suffix) -}}
{{- end -}}
{{- $fullname := include "hanzo.fullname" $context -}}
{{- $maxLen := sub 62 (len $suffix) | int -}}
{{- $truncatedFullname := trunc $maxLen $fullname | trimSuffix "-" -}}
{{- printf "%s-%s" $truncatedFullname $suffix -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "hanzo.chart" -}}
{{- printf "%s-helm" .Chart.Name | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Expand the name of the chart.
*/}}
{{- define "hanzo.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Inject extra environment vars in the format key:value, if populated
*/}}
{{- define "hanzo.extraEnvironmentVars" -}}
{{- if .extraEnvironmentVars -}}
{{- range $key, $value := .extraEnvironmentVars }}
- name: {{ $key }}
  value: {{ $value | quote }}
{{- end -}}
{{- end -}}
{{- end -}}

{{- define "hanzo.mergeExtraEnvironmentVars" -}}
{{- $global := ((.global | default dict).extraEnvironmentVars | default dict) -}}
{{- $component := ((.component | default dict).extraEnvironmentVars | default dict) -}}
{{- $target := .target -}}
{{- range $key, $value := $global }}
{{- $_ := set $target $key $value }}
{{- end }}
{{- range $key, $value := $component }}
{{- $_ := set $target $key $value }}
{{- end }}
{{- end -}}

{{/* Return the proper filer image */}}
{{- define "hanzo.filer.image" -}}
{{- if .Values.filer.imageOverride -}}
{{- $imageOverride := .Values.filer.imageOverride -}}
{{- printf "%s" $imageOverride -}}
{{- else -}}
{{- include "hanzo.image" . }}
{{- end -}}
{{- end -}}

{{/* Return the proper master image */}}
{{- define "hanzo.master.image" -}}
{{- if .Values.master.imageOverride -}}
{{- $imageOverride := .Values.master.imageOverride -}}
{{- printf "%s" $imageOverride -}}
{{- else -}}
{{- include "hanzo.image" . }}
{{- end -}}
{{- end -}}

{{/* Return the proper s3 image */}}
{{- define "hanzo.s3.image" -}}
{{- if .Values.s3.imageOverride -}}
{{- $imageOverride := .Values.s3.imageOverride -}}
{{- printf "%s" $imageOverride -}}
{{- else -}}
{{- include "hanzo.image" . }}
{{- end -}}
{{- end -}}

{{/* Return the proper sftp image */}}
{{- define "hanzo.sftp.image" -}}
{{- if .Values.sftp.imageOverride -}}
{{- $imageOverride := .Values.sftp.imageOverride -}}
{{- printf "%s" $imageOverride -}}
{{- else -}}
{{- include "hanzo.image" . }}
{{- end -}}
{{- end -}}

{{/* Return the proper admin image */}}
{{- define "hanzo.admin.image" -}}
{{- if .Values.admin.imageOverride -}}
{{- $imageOverride := .Values.admin.imageOverride -}}
{{- printf "%s" $imageOverride -}}
{{- else -}}
{{- include "hanzo.image" . }}
{{- end -}}
{{- end -}}

{{/* Return the proper worker image */}}
{{- define "hanzo.worker.image" -}}
{{- if .Values.worker.imageOverride -}}
{{- $imageOverride := .Values.worker.imageOverride -}}
{{- printf "%s" $imageOverride -}}
{{- else -}}
{{- include "hanzo.image" . }}
{{- end -}}
{{- end -}}

{{/* Return the proper volume image */}}
{{- define "hanzo.volume.image" -}}
{{- if .Values.volume.imageOverride -}}
{{- $imageOverride := .Values.volume.imageOverride -}}
{{- printf "%s" $imageOverride -}}
{{- else -}}
{{- include "hanzo.image" . }}
{{- end -}}
{{- end -}}

{{/* Computes the container image name for all components (if they are not overridden) */}}
{{- define "hanzo.image" -}}
{{- $registryName := default .Values.image.registry .Values.global.imageRegistry | toString -}}
{{- $repositoryName := default .Values.image.repository .Values.global.s3.image.repository | toString -}}
{{- $name := .Values.global.s3.image.name | toString -}}
{{- $tag := default .Chart.AppVersion .Values.image.tag  | toString -}}
{{- if .Values.image.repository -}}
{{-   $name = $repositoryName -}}
{{- else if $repositoryName -}}
{{-   $name = printf "%s/%s" (trimSuffix "/" $repositoryName) (base $name) -}}
{{- end -}}
{{- if $registryName -}}
{{-   printf "%s/%s:%s" $registryName $name $tag -}}
{{- else -}}
{{-   printf "%s:%s" $name $tag -}}
{{- end -}}
{{- end -}}

{{/* check if any Volume PVC exists */}}
{{- define "hanzo.volume.pvc_exists" -}}
{{- if or (or (eq .Values.volume.data.type "persistentVolumeClaim") (and (eq .Values.volume.idx.type "persistentVolumeClaim") .Values.volume.dir_idx )) (eq .Values.volume.logs.type "persistentVolumeClaim") -}}
{{- printf "true" -}}
{{- else -}}
{{- printf "" -}}
{{- end -}}
{{- end -}}

{{/* check if any Filer PVC exists */}}
{{- define "hanzo.filer.pvc_exists" -}}
{{- if or (eq .Values.filer.data.type "persistentVolumeClaim") (eq .Values.filer.logs.type "persistentVolumeClaim") -}}
{{- printf "true" -}}
{{- else -}}
{{- printf "" -}}
{{- end -}}
{{- end -}}

{{/* check if any Master PVC exists */}}
{{- define "hanzo.master.pvc_exists" -}}
{{- if or (eq .Values.master.data.type "persistentVolumeClaim") (eq .Values.master.logs.type "persistentVolumeClaim") -}}
{{- printf "true" -}}
{{- else -}}
{{- printf "" -}}
{{- end -}}
{{- end -}}

{{/* check if any Admin PVC exists */}}
{{- define "hanzo.admin.pvc_exists" -}}
{{- if or (eq .Values.admin.data.type "persistentVolumeClaim") (eq .Values.admin.logs.type "persistentVolumeClaim") -}}
{{- printf "true" -}}
{{- else -}}
{{- printf "" -}}
{{- end -}}
{{- end -}}

{{/* check if any InitContainers exist for Volumes */}}
{{- define "hanzo.volume.initContainers_exists" -}}
{{- if or (not (empty .Values.volume.idx )) (not (empty .Values.volume.initContainers )) -}}
{{- printf "true" -}}
{{- else -}}
{{- printf "" -}}
{{- end -}}
{{- end -}}

{{/* Return the proper imagePullSecrets */}}
{{- define "hanzo.imagePullSecrets" -}}
{{- with .Values.global.imagePullSecrets }}
imagePullSecrets:
{{- if kindIs "string" . }}
  - name: {{ . }}
{{- else }}
{{- range . }}
  {{- if kindIs "string" . }}
  - name: {{ . }}
  {{- else }}
  - {{ toYaml . }}
  {{- end}}
{{- end }}
{{- end }}
{{- end }}
{{- end -}}

{{/*
Renders a value that contains template perhaps with scope if the scope is present.
Usage:
{{ include "hanzo.tplvalues.render" ( dict "value" .Values.path.to.the.Value "context" $ ) }}
{{ include "hanzo.tplvalues.render" ( dict "value" .Values.path.to.the.Value "context" $ "scope" $app ) }}
*/}}
{{- define "hanzo.tplvalues.render" -}}
{{- $value := typeIs "string" .value | ternary .value (.value | toYaml) }}
{{- if contains "{{" (toJson .value) }}
  {{- if .scope }}
      {{- tpl (cat "{{- with $.RelativeScope -}}" $value "{{- end }}") (merge (dict "RelativeScope" .scope) .context) }}
  {{- else }}
    {{- tpl $value .context }}
  {{- end }}
{{- else }}
    {{- $value }}
{{- end }}
{{- end -}}

{{/*
Converts a Kubernetes quantity like "256Mi" or "2G" to a float64 in base units,
handling both binary (Ki, Mi, Gi) and decimal (m, k, M) suffixes; numeric inputs
Usage:
{{ include "hanzo.resource-quantity" "10Gi" }}
*/}}
{{- define "hanzo.resource-quantity" -}}
    {{- $value := . -}}
    {{- $unit := 1.0 -}}
    {{- if typeIs "string" . -}}
        {{- $base2 := dict "Ki" 0x1p10 "Mi" 0x1p20 "Gi" 0x1p30 "Ti" 0x1p40 "Pi" 0x1p50 "Ei" 0x1p60 -}}
        {{- $base10 := dict "m" 1e-3 "k" 1e3 "M" 1e6 "G" 1e9 "T" 1e12 "P" 1e15 "E" 1e18 -}}
        {{- range $k, $v := merge $base2 $base10 -}}
            {{- if hasSuffix $k $ -}}
                {{- $value = trimSuffix $k $ -}}
                {{- $unit = $v -}}
            {{- end -}}
        {{- end -}}
    {{- end -}}
    {{- mulf (float64 $value) $unit -}}
{{- end -}}

{{/*
getOrGeneratePassword will check if a password exists in a secret and return it,
or generate a new random password if it doesn't exist.
*/}}
{{- define "hanzo.getOrGeneratePassword" -}}
{{- $params := . -}}
{{- $namespace := $params.namespace -}}
{{- $secretName := $params.secretName -}}
{{- $key := $params.key -}}
{{- $length := default 16 $params.length -}}

{{- $existingSecret := default (lookup "v1" "Secret" $namespace $secretName) $params.existingSecret -}}
{{- if and $existingSecret (index $existingSecret.data $key) -}}
  {{- index $existingSecret.data $key | b64dec -}}
{{- else -}}
  {{- randAlphaNum $length -}}
{{- end -}}
{{- end -}}

{{/*
Compute the master service address to be used in cluster env vars.
If allInOne is enabled, point to the all-in-one service; otherwise, point to the master service.
*/}}
{{- define "hanzo.cluster.masterAddress" -}}
{{- $component := ternary "all-in-one" "master" .Values.allInOne.enabled -}}
{{- printf "%s.%s:%d" (include "hanzo.componentName" (list . $component)) .Release.Namespace (int .Values.master.port) -}}
{{- end -}}

{{/*
Compute the filer service address to be used in cluster env vars.
If allInOne is enabled, point to the all-in-one service; otherwise, point to the filer-client service.
*/}}
{{- define "hanzo.cluster.filerAddress" -}}
{{- $component := ternary "all-in-one" "filer-client" .Values.allInOne.enabled -}}
{{- printf "%s.%s:%d" (include "hanzo.componentName" (list . $component)) .Release.Namespace (int .Values.filer.port) -}}
{{- end -}}

{{/*
Generate comma-separated list of master server addresses.
Usage: {{ include "hanzo.masterServers" . }}
Output example: my-release-master-0.my-release-master.namespace:9333,my-release-master-1...
*/}}
{{- define "hanzo.masterServers" -}}
{{- $masterName := include "hanzo.componentName" (list . "master") -}}
{{- range $index := until (.Values.master.replicas | int) -}}
{{- if $index }},{{ end -}}
{{ $masterName }}-{{ $index }}.{{ $masterName }}.{{ $.Release.Namespace }}:{{ $.Values.master.port }}
{{- end -}}
{{- end -}}

{{/*
Generate master server argument value, using global.masterServer if set, otherwise the generated list.
Usage: {{ include "hanzo.masterServerArg" . }}
*/}}
{{- define "hanzo.masterServerArg" -}}
{{- if .Values.global.s3.masterServer -}}
{{- .Values.global.s3.masterServer -}}
{{- else -}}
{{- include "hanzo.masterServers" . -}}
{{- end -}}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "hanzo.serviceAccountName" -}}
{{- .Values.global.s3.serviceAccountName | default "s3" -}}
{{- end -}}

{{/* True when security.toml should be rendered and mounted. volumeWrite is
     excluded unless its non-default expiration is configured. */}}
{{- define "hanzo.securityConfigEnabled" -}}
{{- $sec := (.Values.global.s3).securityConfig | default dict -}}
{{- $jwt := $sec.jwtSigning | default dict -}}
{{- $expiresAfterSeconds := $jwt.expiresAfterSeconds | default dict -}}
{{- $volumeWriteExpirationConfigured := and $jwt.volumeWrite (gt (int $expiresAfterSeconds.volumeWrite) 0) -}}
{{- if or .Values.global.s3.enableSecurity $volumeWriteExpirationConfigured $jwt.volumeRead $jwt.filerWrite $jwt.filerRead -}}
true
{{- end -}}
{{- end -}}

{{/* S3 TLS cert/key arguments, using custom secret if s3.tlsSecret is set */}}
{{- define "hanzo.s3.tlsArgs" -}}
{{- $prefix := .prefix -}}
{{- $root := .root -}}
{{- if $root.Values.s3.tlsSecret -}}
-{{ $prefix }}cert.file=/usr/local/share/ca-certificates/s3/tls.crt \
-{{ $prefix }}key.file=/usr/local/share/ca-certificates/s3/tls.key \
{{- else -}}
-{{ $prefix }}cert.file=/usr/local/share/ca-certificates/client/tls.crt \
-{{ $prefix }}key.file=/usr/local/share/ca-certificates/client/tls.key \
{{- end -}}
{{- end -}}

{{/* S3 custom TLS volume mount */}}
{{- define "hanzo.s3.tlsVolumeMount" -}}
{{- if .Values.s3.tlsSecret }}
- name: s3-tls-cert
  readOnly: true
  mountPath: /usr/local/share/ca-certificates/s3/
{{- end }}
{{- end -}}

{{/* S3 custom TLS volume */}}
{{- define "hanzo.s3.tlsVolume" -}}
{{- if .Values.s3.tlsSecret }}
- name: s3-tls-cert
  secret:
    secretName: {{ .Values.s3.tlsSecret }}
{{- end }}
{{- end -}}

{{/* Generate a compatible trafficDistribution value due to "PreferClose" fast deprecation in k8s v1.35.
     Accepts a dict with "value" (the trafficDistribution string) and "Capabilities". */}}
{{- define "hanzo.trafficDistribution" -}}
{{- if .value -}}
{{- and (eq .value "PreferClose") (semverCompare ">=1.35-0" .Capabilities.KubeVersion.GitVersion) | ternary "PreferSameZone" .value -}}
{{- end -}}
{{- end -}}
