module github.com/GoogleCloudPlatform/run-gmp-sidecar

go 1.25.0

require (
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector v0.54.0
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector/googlemanagedprometheus v0.54.0
	github.com/go-kit/log v0.2.1
	github.com/goccy/go-yaml v1.19.2
	github.com/gogo/protobuf v1.3.2
	github.com/golang/snappy v1.0.0
	github.com/mitchellh/hashstructure/v2 v2.0.2
	github.com/mitchellh/mapstructure v1.5.1-0.20231216201459-8508981c8b6c
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter v0.143.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudexporter v0.143.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter v0.143.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatautil v0.143.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheus v0.143.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor v0.143.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor v0.143.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricstransformprocessor v0.143.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor v0.143.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor v0.143.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor v0.143.0
	github.com/open-telemetry/opentelemetry-collector-contrib/testbed v0.143.0
	github.com/prometheus/client_golang v1.23.2
	github.com/prometheus/common v0.67.5
	github.com/prometheus/prometheus v0.309.1
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/stretchr/testify v1.11.1
	go.opentelemetry.io/collector/component v1.49.0
	go.opentelemetry.io/collector/component/componentstatus v0.143.0
	go.opentelemetry.io/collector/component/componenttest v0.143.0
	go.opentelemetry.io/collector/confmap v1.49.0
	go.opentelemetry.io/collector/confmap/provider/envprovider v1.49.0
	go.opentelemetry.io/collector/confmap/provider/fileprovider v1.49.0
	go.opentelemetry.io/collector/confmap/provider/httpprovider v1.49.0
	go.opentelemetry.io/collector/confmap/provider/httpsprovider v1.49.0
	go.opentelemetry.io/collector/confmap/provider/yamlprovider v1.49.0
	go.opentelemetry.io/collector/consumer v1.49.0
	go.opentelemetry.io/collector/consumer/consumertest v0.143.0
	go.opentelemetry.io/collector/exporter v1.49.0
	go.opentelemetry.io/collector/exporter/debugexporter v0.143.0
	go.opentelemetry.io/collector/exporter/exporterhelper v0.143.0
	go.opentelemetry.io/collector/exporter/exportertest v0.143.0
	go.opentelemetry.io/collector/exporter/otlpexporter v0.143.0
	go.opentelemetry.io/collector/exporter/otlphttpexporter v0.143.0
	go.opentelemetry.io/collector/extension v1.49.0
	go.opentelemetry.io/collector/extension/zpagesextension v0.143.0
	go.opentelemetry.io/collector/featuregate v1.49.0
	go.opentelemetry.io/collector/otelcol v0.143.0
	go.opentelemetry.io/collector/otelcol/otelcoltest v0.143.0
	go.opentelemetry.io/collector/pdata v1.49.0
	go.opentelemetry.io/collector/processor v1.49.0
	go.opentelemetry.io/collector/processor/batchprocessor v0.143.0
	go.opentelemetry.io/collector/processor/memorylimiterprocessor v0.143.0
	go.opentelemetry.io/collector/receiver v1.49.0
	go.opentelemetry.io/collector/receiver/otlpreceiver v0.143.0
	go.opentelemetry.io/collector/receiver/receiverhelper v0.143.0
	go.opentelemetry.io/collector/receiver/receivertest v0.143.0
	go.opentelemetry.io/collector/semconv v0.128.1-0.20250610090210-188191247685
	go.uber.org/goleak v1.3.0
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.27.1
	golang.org/x/text v0.33.0
	google.golang.org/protobuf v1.36.11
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools/v3 v3.5.2
	k8s.io/apimachinery v0.35.0
)

require (
	cloud.google.com/go v0.123.0 // indirect
	cloud.google.com/go/auth v0.17.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	cloud.google.com/go/logging v1.13.0 // indirect
	cloud.google.com/go/longrunning v0.6.7 // indirect
	cloud.google.com/go/monitoring v1.24.2 // indirect
	cloud.google.com/go/trace v1.11.6 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.20.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.13.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.11.2 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5 v5.7.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4 v4.3.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.6.0 // indirect
	github.com/Code-Hex/go-generics-cache v1.5.1 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.30.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.30.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.54.0 // indirect
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/Masterminds/semver/v3 v3.4.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/Showmax/go-fqdn v1.0.0 // indirect
	github.com/alecthomas/participle/v2 v2.1.4 // indirect
	github.com/alecthomas/units v0.0.0-20240927000941-0f3dac36c52b // indirect
	github.com/antchfx/xmlquery v1.5.0 // indirect
	github.com/antchfx/xpath v1.3.5 // indirect
	github.com/apache/arrow-go/v18 v18.2.0 // indirect
	github.com/apache/thrift v0.22.0 // indirect
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/aws/aws-sdk-go-v2 v1.41.0 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.32.6 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.19.6 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.279.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ecs v1.69.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/lightsail v1.50.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.41.5 // indirect
	github.com/aws/smithy-go v1.24.0 // indirect
	github.com/axiomhq/hyperloglog v0.2.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cncf/xds/go v0.0.0-20251022180443-0feb69152e9f // indirect
	github.com/containerd/errdefs v1.0.0 // indirect
	github.com/containerd/errdefs/pkg v0.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dennwc/varint v1.0.0 // indirect
	github.com/dgryski/go-metro v0.0.0-20180109044635-280f6062b5bc // indirect
	github.com/digitalocean/go-metadata v0.0.0-20250129100319-e3650a3df44b // indirect
	github.com/digitalocean/godo v1.171.0 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/docker v28.5.2+incompatible // indirect
	github.com/docker/go-connections v0.6.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/ebitengine/purego v0.9.1 // indirect
	github.com/elastic/go-grok v0.3.1 // indirect
	github.com/elastic/lunes v0.2.0 // indirect
	github.com/emicklei/go-restful/v3 v3.12.2 // indirect
	github.com/envoyproxy/go-control-plane/envoy v1.36.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.3.0 // indirect
	github.com/expr-lang/expr v1.17.7 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/foxboron/go-tpm-keyfiles v0.0.0-20250903184740-5d135037bd4d // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/fxamacker/cbor/v2 v2.9.0 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-openapi/jsonpointer v0.22.1 // indirect
	github.com/go-openapi/jsonreference v0.21.3 // indirect
	github.com/go-openapi/swag v0.25.4 // indirect
	github.com/go-openapi/swag/cmdutils v0.25.4 // indirect
	github.com/go-openapi/swag/conv v0.25.4 // indirect
	github.com/go-openapi/swag/fileutils v0.25.4 // indirect
	github.com/go-openapi/swag/jsonname v0.25.4 // indirect
	github.com/go-openapi/swag/jsonutils v0.25.4 // indirect
	github.com/go-openapi/swag/loading v0.25.4 // indirect
	github.com/go-openapi/swag/mangling v0.25.4 // indirect
	github.com/go-openapi/swag/netutils v0.25.4 // indirect
	github.com/go-openapi/swag/stringutils v0.25.4 // indirect
	github.com/go-openapi/swag/typeutils v0.25.4 // indirect
	github.com/go-openapi/swag/yamlutils v0.25.4 // indirect
	github.com/go-resty/resty/v2 v2.17.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/go-zookeeper/zk v1.0.4 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/flatbuffers v25.2.10+incompatible // indirect
	github.com/google/gnostic-models v0.7.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/go-tpm v0.9.8 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.7 // indirect
	github.com/googleapis/gax-go/v2 v2.15.0 // indirect
	github.com/gophercloud/gophercloud/v2 v2.9.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/websocket v1.5.4-0.20250319132907-e064f32e3674 // indirect
	github.com/grafana/regexp v0.0.0-20250905093917-f7b3be9d1853 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3 // indirect
	github.com/hashicorp/consul/api v1.32.1 // indirect
	github.com/hashicorp/cronexpr v1.1.3 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-version v1.8.0 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/hashicorp/nomad/api v0.0.0-20251216171439-1dee0671280e // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/hetznercloud/hcloud-go/v2 v2.33.0 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/ionos-cloud/sdk-go/v6 v6.3.5 // indirect
	github.com/jaegertracing/jaeger-idl v0.6.0 // indirect
	github.com/jonboulle/clockwork v0.5.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kamstrup/intmap v0.5.1 // indirect
	github.com/klauspost/compress v1.18.2 // indirect
	github.com/klauspost/cpuid/v2 v2.2.11 // indirect
	github.com/knadh/koanf/maps v0.1.2 // indirect
	github.com/knadh/koanf/providers/confmap v1.0.0 // indirect
	github.com/knadh/koanf/v2 v2.3.0 // indirect
	github.com/kolo/xmlrpc v0.0.0-20220921171641-a4b6fa1dd06b // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/leodido/go-syslog/v4 v4.3.0 // indirect
	github.com/leodido/ragel-machinery v0.0.0-20190525184631-5f46317e436b // indirect
	github.com/lightstep/go-expohisto v1.0.0 // indirect
	github.com/linode/go-metadata v0.2.3 // indirect
	github.com/linode/linodego v1.63.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20250317134145-8bc96cf8fc35 // indirect
	github.com/magefile/mage v1.15.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/miekg/dns v1.1.69 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.3-0.20250322232337-35a7c28c31ee // indirect
	github.com/mostynb/go-grpc-compression v1.2.3 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/connector/routingconnector v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/connector/spanmetricsconnector v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/stefexporter v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/syslogexporter v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/zipkinexporter v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/ecsutil v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/common v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/filter v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/grpcutil v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/metadataproviders v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/otelarrow v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/pdatautil v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/sharedcomponent v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/core/xidutils v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/golden v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/resourcetotelemetry v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/jaeger v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheusremotewrite v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/zipkin v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jaegerreceiver v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/otelarrowreceiver v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/stefreceiver v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/syslogreceiver v0.143.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zipkinreceiver v0.143.0 // indirect
	github.com/open-telemetry/otel-arrow/go v0.45.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/openshift/api v0.0.0-20251015095338-264e80a2b6e7 // indirect
	github.com/openshift/client-go v0.0.0-20251015124057-db0dee36e235 // indirect
	github.com/openzipkin/zipkin-go v0.4.3 // indirect
	github.com/ovh/go-ovh v1.9.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.23 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/prometheus/client_golang/exp v0.0.0-20251212205219-7ba246a648ca // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/otlptranslator v1.0.0 // indirect
	github.com/prometheus/procfs v0.19.2 // indirect
	github.com/prometheus/sigv4 v0.3.0 // indirect
	github.com/rs/cors v1.11.1 // indirect
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.36 // indirect
	github.com/shirou/gopsutil/v4 v4.25.11 // indirect
	github.com/spf13/cobra v1.10.2 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/splunk/stef/go/grpc v0.0.9 // indirect
	github.com/splunk/stef/go/otel v0.0.9 // indirect
	github.com/splunk/stef/go/pdata v0.0.9 // indirect
	github.com/splunk/stef/go/pkg v0.0.9 // indirect
	github.com/stackitcloud/stackit-sdk-go/core v0.20.1 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/match v1.2.0 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/tinylru v1.2.1 // indirect
	github.com/tidwall/wal v1.2.1 // indirect
	github.com/tklauser/go-sysconf v0.3.16 // indirect
	github.com/tklauser/numcpus v0.11.0 // indirect
	github.com/twmb/murmur3 v1.1.8 // indirect
	github.com/ua-parser/uap-go v0.0.0-20240611065828-3a4781585db6 // indirect
	github.com/valyala/fastjson v1.6.7 // indirect
	github.com/vultr/govultr/v2 v2.17.2 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/collector v0.143.0 // indirect
	go.opentelemetry.io/collector/client v1.49.0 // indirect
	go.opentelemetry.io/collector/config/configauth v1.49.0 // indirect
	go.opentelemetry.io/collector/config/configcompression v1.49.0 // indirect
	go.opentelemetry.io/collector/config/configgrpc v0.143.0 // indirect
	go.opentelemetry.io/collector/config/confighttp v0.143.0 // indirect
	go.opentelemetry.io/collector/config/configmiddleware v1.49.0 // indirect
	go.opentelemetry.io/collector/config/confignet v1.49.0 // indirect
	go.opentelemetry.io/collector/config/configopaque v1.49.0 // indirect
	go.opentelemetry.io/collector/config/configoptional v1.49.0 // indirect
	go.opentelemetry.io/collector/config/configretry v1.49.0 // indirect
	go.opentelemetry.io/collector/config/configtelemetry v0.143.0 // indirect
	go.opentelemetry.io/collector/config/configtls v1.49.0 // indirect
	go.opentelemetry.io/collector/confmap/xconfmap v0.143.0 // indirect
	go.opentelemetry.io/collector/connector v0.143.0 // indirect
	go.opentelemetry.io/collector/connector/connectortest v0.143.0 // indirect
	go.opentelemetry.io/collector/connector/xconnector v0.143.0 // indirect
	go.opentelemetry.io/collector/consumer/consumererror v0.143.0 // indirect
	go.opentelemetry.io/collector/consumer/consumererror/xconsumererror v0.143.0 // indirect
	go.opentelemetry.io/collector/consumer/xconsumer v0.143.0 // indirect
	go.opentelemetry.io/collector/exporter/exporterhelper/xexporterhelper v0.143.0 // indirect
	go.opentelemetry.io/collector/exporter/xexporter v0.143.0 // indirect
	go.opentelemetry.io/collector/extension/extensionauth v1.49.0 // indirect
	go.opentelemetry.io/collector/extension/extensioncapabilities v0.143.0 // indirect
	go.opentelemetry.io/collector/extension/extensionmiddleware v0.143.0 // indirect
	go.opentelemetry.io/collector/extension/extensiontest v0.143.0 // indirect
	go.opentelemetry.io/collector/extension/xextension v0.143.0 // indirect
	go.opentelemetry.io/collector/internal/fanoutconsumer v0.143.0 // indirect
	go.opentelemetry.io/collector/internal/memorylimiter v0.143.0 // indirect
	go.opentelemetry.io/collector/internal/sharedcomponent v0.143.0 // indirect
	go.opentelemetry.io/collector/internal/telemetry v0.143.0 // indirect
	go.opentelemetry.io/collector/pdata/pprofile v0.143.0 // indirect
	go.opentelemetry.io/collector/pdata/testdata v0.143.0 // indirect
	go.opentelemetry.io/collector/pdata/xpdata v0.143.0 // indirect
	go.opentelemetry.io/collector/pipeline v1.49.0 // indirect
	go.opentelemetry.io/collector/pipeline/xpipeline v0.143.0 // indirect
	go.opentelemetry.io/collector/processor/processorhelper v0.143.0 // indirect
	go.opentelemetry.io/collector/processor/processorhelper/xprocessorhelper v0.143.0 // indirect
	go.opentelemetry.io/collector/processor/processortest v0.143.0 // indirect
	go.opentelemetry.io/collector/processor/xprocessor v0.143.0 // indirect
	go.opentelemetry.io/collector/receiver/xreceiver v0.143.0 // indirect
	go.opentelemetry.io/collector/service v0.143.0 // indirect
	go.opentelemetry.io/collector/service/hostcapabilities v0.143.0 // indirect
	go.opentelemetry.io/contrib/bridges/otelzap v0.13.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.63.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.64.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.64.0 // indirect
	go.opentelemetry.io/contrib/otelconf v0.18.0 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.38.0 // indirect
	go.opentelemetry.io/contrib/zpages v0.63.0 // indirect
	go.opentelemetry.io/otel v1.39.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc v0.14.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.14.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.39.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.39.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.39.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.39.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.39.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.60.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.14.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.39.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.39.0 // indirect
	go.opentelemetry.io/otel/log v0.15.0 // indirect
	go.opentelemetry.io/otel/metric v1.39.0 // indirect
	go.opentelemetry.io/otel/sdk v1.39.0 // indirect
	go.opentelemetry.io/otel/sdk/log v0.14.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.39.0 // indirect
	go.opentelemetry.io/otel/trace v1.39.0 // indirect
	go.opentelemetry.io/proto/otlp v1.9.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.yaml.in/yaml/v2 v2.4.3 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.46.0 // indirect
	golang.org/x/exp v0.0.0-20251209150349-8475f28825e9 // indirect
	golang.org/x/mod v0.31.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/oauth2 v0.34.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/telemetry v0.0.0-20251203150158-8fff8a5912fc // indirect
	golang.org/x/term v0.38.0 // indirect
	golang.org/x/time v0.14.0 // indirect
	golang.org/x/tools v0.40.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	gonum.org/v1/gonum v0.16.0 // indirect
	google.golang.org/api v0.257.0 // indirect
	google.golang.org/genproto v0.0.0-20250922171735-9219d122eba9 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251213004720-97cd9d5aeac2 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/grpc v1.78.0 // indirect
	gopkg.in/evanphx/json-patch.v4 v4.13.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.34.3 // indirect
	k8s.io/client-go v0.34.3 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/kube-openapi v0.0.0-20250910181357-589584f1c912 // indirect
	k8s.io/utils v0.0.0-20251002143259-bc988d571ff4 // indirect
	modernc.org/b/v2 v2.1.10 // indirect
	sigs.k8s.io/json v0.0.0-20250730193827-2d320260d730 // indirect
	sigs.k8s.io/randfill v1.0.0 // indirect
	sigs.k8s.io/structured-merge-diff/v6 v6.3.0 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)
