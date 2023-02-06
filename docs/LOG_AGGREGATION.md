# Log Aggregation

## Introduction
The log aggregation layer is based on Loki and OpenTelemetry (OTLP) collector. Logs are fetched, processed and exported to Loki using the OTLP collector.
Ingested logs can be further analyzed using Grafana.

## OpenTelemetry Collector
The OpenTelemetry Collector offers a vendor-agnostic implementation which supports the collection of Kubernetes container logs and their export to the Loki backend.
The collector is deployed as a daemon set (agent mode) which is the most suitable for collecting container logs. It has been deployed using Helm Chart
rather than the Operator since the Helm chart support presets for `log collection` and `kubernetes attributes` while these features are still
missing from the OpenTelemetryCollector CRD used by the Operator.

Loki backend (exporter) is supported through [Contrib repository](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/lokiexporter) where components that are not suitable for the core repository are placed.
In the processor stage, the Loki exporter is used to convert OTLP resources and log attributes into Loki labels, which are indexed.
In that way, labels like `name_space`, `pod_name`, `container_name`, and `log_iostream`, can be used for issuing fast log queries from Grafana.

## Loki
Loki is a log aggregation system which is simple and easy to operate. It does not index the contents of the logs, but rather a set of labels for each log stream.
It has been deployed in `single binary` mode with the logs stored in the file system. Parameters like `parallelise_shardable_queries=false`,
`split_queries_by_interval=0` and `max_outstanding_requests_per_tenant=10000` are used to tune performance for mentioned deployment mode.

## Grafana
Grafana is deployed independently of Loki to be used by both Loki and Kube Prometheus applications.
It can be accessed on `https://grafana.kube.local` through MetalLB load balancer, but first resolving needs to be configured in /etc/hosts file:
```bash
172.18.255.200  grafana.kube.local
```
At this moment logging in to Grafana is done by using default credentials defined in [Grafana values](https://github.com/nearform/k8s-kurated-addons/blob/main/addons/grafana/values.yaml#L6)
In the future default password will be replaced with generated one.
Grafana is preconfigured to use Loki data source and logs are available shortly upon stack deployment. Custom dashboards for Loki are not included with Grafana, rather idea is to use Grafana’s `explore` feature.

Log filtering is based on LogQL, a Grafana Loki’s PromQL-inspired query language. It uses labels and operators for filtering. A few examples are shown below:
```bash
{exporter="OTLP"} |= "error"  # uses exporter=OTLP label to select all logs and further filter for errors
```
```bash
{container_name =~ "istio.+"} |= "timeout"  # uses container_name label and regular expression to filter for all istio containers and timeout keyword
```

More details can be found in Grafana [Log queries documentation](https://grafana.com/docs/loki/latest/logql/log_queries/)
