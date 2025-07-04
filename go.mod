module github.com/GoogleCloudPlatform/run-gmp-sidecar

go 1.23.0

toolchain go1.24.1

require (
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector v0.49.0
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector/googlemanagedprometheus v0.49.0
	github.com/goccy/go-yaml v1.11.2
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter v0.113.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudexporter v0.113.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter v0.113.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatautil v0.113.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheus v0.113.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor v0.113.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor v0.113.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricstransformprocessor v0.113.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor v0.113.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor v0.113.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor v0.113.0
	github.com/open-telemetry/opentelemetry-collector-contrib/testbed v0.113.0
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/stretchr/testify v1.9.0
	go.opentelemetry.io/collector/component v0.113.0
	go.opentelemetry.io/collector/confmap v1.19.0
	go.opentelemetry.io/collector/confmap/provider/fileprovider v1.19.0
	go.opentelemetry.io/collector/consumer v0.113.0
	go.opentelemetry.io/collector/consumer/consumertest v0.113.0
	go.opentelemetry.io/collector/exporter v0.113.0
	go.opentelemetry.io/collector/exporter/otlpexporter v0.113.0
	go.opentelemetry.io/collector/exporter/otlphttpexporter v0.113.0
	go.opentelemetry.io/collector/extension v0.113.0
	go.opentelemetry.io/collector/extension/zpagesextension v0.113.0
	go.opentelemetry.io/collector/featuregate v1.19.0
	go.opentelemetry.io/collector/otelcol v0.113.0
	go.opentelemetry.io/collector/otelcol/otelcoltest v0.112.0
	go.opentelemetry.io/collector/pdata v1.19.0
	go.opentelemetry.io/collector/processor v0.113.0
	go.opentelemetry.io/collector/processor/batchprocessor v0.113.0
	go.opentelemetry.io/collector/processor/memorylimiterprocessor v0.113.0
	go.opentelemetry.io/collector/receiver v0.113.0
	go.opentelemetry.io/collector/receiver/otlpreceiver v0.113.0
	go.opentelemetry.io/collector/semconv v0.113.0
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.27.0
	golang.org/x/text v0.23.0
	gotest.tools/v3 v3.5.0
)

require (
	github.com/prometheus/common v0.60.1
	go.opentelemetry.io/collector/component/componentstatus v0.113.0
	go.opentelemetry.io/collector/confmap/provider/envprovider v1.19.0
	go.opentelemetry.io/collector/confmap/provider/httpprovider v1.19.0
	go.opentelemetry.io/collector/confmap/provider/httpsprovider v1.17.0
	go.opentelemetry.io/collector/confmap/provider/yamlprovider v1.19.0
	go.opentelemetry.io/collector/exporter/debugexporter v0.113.0
	go.opentelemetry.io/collector/exporter/exportertest v0.113.0
	go.opentelemetry.io/collector/receiver/receivertest v0.113.0
)

require (
	cloud.google.com/go/auth v0.10.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.5 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.13.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.7.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.10.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5 v5.7.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4 v4.3.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.2.2 // indirect
	github.com/Code-Hex/go-generics-cache v1.5.1 // indirect
	github.com/antchfx/xmlquery v1.4.2 // indirect
	github.com/antchfx/xpath v1.3.2 // indirect
	github.com/apache/thrift v0.21.0 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/ebitengine/purego v0.8.1 // indirect
	github.com/elastic/go-grok v0.3.1 // indirect
	github.com/elastic/lunes v0.1.0 // indirect
	github.com/expr-lang/expr v1.17.0 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/go-playground/validator/v10 v10.15.5 // indirect
	github.com/go-viper/mapstructure/v2 v2.3.0 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/s2a-go v0.1.8 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.22.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/hetznercloud/hcloud-go/v2 v2.10.2 // indirect
	github.com/jaegertracing/jaeger v1.62.0 // indirect
	github.com/jonboulle/clockwork v0.4.0 // indirect
	github.com/knadh/koanf/v2 v2.1.1 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/leodido/go-syslog/v4 v4.2.0 // indirect
	github.com/leodido/ragel-machinery v0.0.0-20190525184631-5f46317e436b // indirect
	github.com/lightstep/go-expohisto v1.0.0 // indirect
	github.com/magefile/mage v1.15.0 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/connector/routingconnector v0.113.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/connector/spanmetricsconnector v0.113.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/opencensusexporter v0.113.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/syslogexporter v0.113.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/zipkinexporter v0.113.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/filter v0.113.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig v0.113.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/pdatautil v0.113.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/golden v0.113.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/resourcetotelemetry v0.113.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza v0.113.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/jaeger v0.113.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/opencensus v0.113.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheusremotewrite v0.113.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/zipkin v0.113.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jaegerreceiver v0.113.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/opencensusreceiver v0.113.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/syslogreceiver v0.113.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zipkinreceiver v0.113.0 // indirect
	github.com/openshift/api v3.9.0+incompatible // indirect
	github.com/openshift/client-go v0.0.0-20210521082421-73d9475a9142 // indirect
	github.com/openzipkin/zipkin-go v0.4.3 // indirect
	github.com/ovh/go-ovh v1.6.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/shirou/gopsutil/v4 v4.24.10 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/tinylru v1.2.1 // indirect
	github.com/tidwall/wal v1.1.7 // indirect
	github.com/ua-parser/uap-go v0.0.0-20240611065828-3a4781585db6 // indirect
	github.com/valyala/fastjson v1.6.4 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	go.opentelemetry.io/collector v0.113.0 // indirect
	go.opentelemetry.io/collector/client v1.19.0 // indirect
	go.opentelemetry.io/collector/config/configauth v0.113.0 // indirect
	go.opentelemetry.io/collector/config/configcompression v1.19.0 // indirect
	go.opentelemetry.io/collector/config/configgrpc v0.113.0 // indirect
	go.opentelemetry.io/collector/config/confighttp v0.113.0 // indirect
	go.opentelemetry.io/collector/config/confignet v1.19.0 // indirect
	go.opentelemetry.io/collector/config/configopaque v1.19.0 // indirect
	go.opentelemetry.io/collector/config/configretry v1.19.0 // indirect
	go.opentelemetry.io/collector/config/configtelemetry v0.113.0 // indirect
	go.opentelemetry.io/collector/config/configtls v1.19.0 // indirect
	go.opentelemetry.io/collector/config/internal v0.113.0 // indirect
	go.opentelemetry.io/collector/connector v0.113.0 // indirect
	go.opentelemetry.io/collector/connector/connectorprofiles v0.113.0 // indirect
	go.opentelemetry.io/collector/connector/connectortest v0.113.0 // indirect
	go.opentelemetry.io/collector/consumer/consumererror v0.113.0 // indirect
	go.opentelemetry.io/collector/consumer/consumererror/consumererrorprofiles v0.113.0 // indirect
	go.opentelemetry.io/collector/consumer/consumerprofiles v0.113.0 // indirect
	go.opentelemetry.io/collector/exporter/exporterhelper/exporterhelperprofiles v0.113.0 // indirect
	go.opentelemetry.io/collector/exporter/exporterprofiles v0.113.0 // indirect
	go.opentelemetry.io/collector/extension/auth v0.113.0 // indirect
	go.opentelemetry.io/collector/extension/experimental/storage v0.113.0 // indirect
	go.opentelemetry.io/collector/extension/extensioncapabilities v0.113.0 // indirect
	go.opentelemetry.io/collector/internal/fanoutconsumer v0.113.0 // indirect
	go.opentelemetry.io/collector/internal/memorylimiter v0.113.0 // indirect
	go.opentelemetry.io/collector/internal/sharedcomponent v0.113.0 // indirect
	go.opentelemetry.io/collector/pdata/pprofile v0.113.0 // indirect
	go.opentelemetry.io/collector/pdata/testdata v0.113.0 // indirect
	go.opentelemetry.io/collector/pipeline v0.113.0 // indirect
	go.opentelemetry.io/collector/pipeline/pipelineprofiles v0.113.0 // indirect
	go.opentelemetry.io/collector/processor/processorprofiles v0.113.0 // indirect
	go.opentelemetry.io/collector/processor/processortest v0.113.0 // indirect
	go.opentelemetry.io/collector/receiver/receiverprofiles v0.113.0 // indirect
	go.opentelemetry.io/collector/service v0.113.0 // indirect
	go.opentelemetry.io/contrib/bridges/otelzap v0.6.0 // indirect
	go.opentelemetry.io/contrib/config v0.10.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.31.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.31.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.31.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.31.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.31.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.7.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.31.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.31.0 // indirect
	go.opentelemetry.io/otel/log v0.7.0 // indirect
	go.opentelemetry.io/otel/sdk/log v0.7.0 // indirect
	go.opentelemetry.io/proto/otlp v1.3.1 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241021214115-324edc3d5d38 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241021214115-324edc3d5d38 // indirect
	google.golang.org/grpc/stats/opentelemetry v0.0.0-20241028142157-ada6787961b3 // indirect
)

require (
	cloud.google.com/go/compute/metadata v0.5.2 // indirect
	cloud.google.com/go/longrunning v0.6.2 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl v0.113.0 // indirect; indir6.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/klog/v2 v2.130.1 // indirect
)

require (
	cloud.google.com/go v0.116.0 // indirect
	cloud.google.com/go/logging v1.12.0 // indirect
	cloud.google.com/go/monitoring v1.21.2 // indirect
	cloud.google.com/go/trace v1.11.2 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.25.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.25.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.49.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/Showmax/go-fqdn v1.0.0 // indirect
	github.com/alecthomas/participle/v2 v2.1.1 // indirect
	github.com/alecthomas/units v0.0.0-20240626203959-61d1e3462e30 // indirect
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/aws/aws-sdk-go v1.55.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cncf/xds/go v0.0.0-20240905190251-b4127c9b8d78 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dennwc/varint v1.0.0 // indirect
	github.com/digitalocean/godo v1.118.0 // indirect
	github.com/docker/docker v27.3.1+incompatible // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/emicklei/go-restful/v3 v3.11.0 // indirect
	github.com/envoyproxy/go-control-plane v0.13.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.1.0 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-kit/log v0.2.1
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-openapi/jsonpointer v0.20.2 // indirect
	github.com/go-openapi/jsonreference v0.20.4 // indirect
	github.com/go-openapi/swag v0.22.9 // indirect
	github.com/go-resty/resty/v2 v2.13.1 // indirect
	github.com/go-zookeeper/zk v1.0.3 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gogo/protobuf v1.3.2
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.4 // indirect
	github.com/googleapis/gax-go/v2 v2.13.0 // indirect
	github.com/gophercloud/gophercloud v1.13.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/grafana/regexp v0.0.0-20240518133315-a468a5bfb3bc // indirect
	github.com/hashicorp/consul/api v1.30.0 // indirect
	github.com/hashicorp/cronexpr v1.1.2 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/nomad/api v0.0.0-20240717122358-3d93bd3778f3 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/ionos-cloud/sdk-go/v6 v6.1.11 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/knadh/koanf v1.5.0 // indirect
	github.com/kolo/xmlrpc v0.0.0-20220921171641-a4b6fa1dd06b // indirect
	github.com/linode/linodego v1.37.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20220913051719-115f729f3c8c // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/miekg/dns v1.1.61 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.2
	github.com/mitchellh/mapstructure v1.5.1-0.20231216201459-8508981c8b6c
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mostynb/go-grpc-compression v1.2.3 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/ecsutil v0.113.0 // indirect; indir6.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/common v0.113.0 // indirect; indir6.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.113.0 // indirect; indir6.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/metadataproviders v0.113.0 // indirect; indir6.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/sharedcomponent v0.113.0 // indirect; indir6.0
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20220216144756-c35f1ee13d7c // indirect
	github.com/prometheus/client_golang v1.20.5
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common/sigv4 v0.1.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/prometheus/prometheus v1.8.2-0.20211119115433-692a54649ed7
	github.com/rs/cors v1.11.1 // indirect
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.29 // indirect
	github.com/spf13/cobra v1.8.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/vultr/govultr/v2 v2.17.2 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.56.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.56.0 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.31.0 // indirect
	go.opentelemetry.io/contrib/zpages v0.56.0 // indirect
	go.opentelemetry.io/otel v1.31.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.53.0 // indirect
	go.opentelemetry.io/otel/metric v1.31.0 // indirect
	go.opentelemetry.io/otel/sdk v1.31.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.31.0 // indirect
	go.opentelemetry.io/otel/trace v1.31.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/goleak v1.3.0
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842 // indirect
	golang.org/x/mod v0.20.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/oauth2 v0.23.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/term v0.30.0 // indirect
	golang.org/x/time v0.7.0 // indirect
	golang.org/x/tools v0.24.0 // indirect
	gonum.org/v1/gonum v0.15.1 // indirect
	google.golang.org/api v0.204.0 // indirect
	google.golang.org/genproto v0.0.0-20241021214115-324edc3d5d38 // indirect
	google.golang.org/grpc v1.67.1 // indirect
	google.golang.org/protobuf v1.35.1
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.31.2 // indirect
	k8s.io/apimachinery v0.31.2
	k8s.io/client-go v0.31.2 // indirect; indirect4	k8s.io/klog/v2 v2.70.1 // indirect
	k8s.io/kube-openapi v0.0.0-20240228011516-70dd3763d340 // indirect
	k8s.io/utils v0.0.0-20240711033017-18e509b52bc8 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

// Currently causes build issues on windows. Downgrading to previous version.
replace github.com/mattn/go-ieproxy v0.0.9 => github.com/mattn/go-ieproxy v0.0.1

// Point to our fork of Prometheus at branch: feature/shutdownscrape_v0.54.1
replace github.com/prometheus/prometheus => github.com/googleCloudPlatform/prometheus v0.0.0-20241010204127-522fcab662ec
