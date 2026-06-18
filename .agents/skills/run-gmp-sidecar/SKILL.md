---
name: run-gmp-sidecar-dev
description: Developer guide, code generation workflow, and testing practices for run-gmp-sidecar.
---

# Developer Guide & Agent Skill: Run-GMP-Sidecar

This guide documents the design, building blocks, and developer workflow for the `run-gmp-sidecar` project. It serves as context for developers and AI assistants to quickly ramp up and make correct contributions.

## Overview
`run-gmp-sidecar` is an OpenTelemetry Collector distribution packaged as a sidecar container for Google Cloud Run. It collects application metrics in Prometheus format and exports them to Google Cloud Managed Service for Prometheus, where they are mapped to the [prometheus_target](https://cloud.google.com/monitoring/api/resources#tag_prometheus_target) monitored resource.

It is built on top of `distrogen` (see the [opentelemetry-operations-collector/cmd/distrogen](https://github.com/GoogleCloudPlatform/opentelemetry-operations-collector/tree/master/cmd/distrogen) command directory), an OpenTelemetry collector distribution generator. It supports running on both:
- **Cloud Run Services**
- **Cloud Run Worker Pools**

See the [Cloud Run Container Contract](https://docs.cloud.google.com/run/docs/container-contract) for runtime environment requirements.

---

## Architectural Components

### 1. `spec.yaml`
This file defines the collector distribution recipe. It lists all OpenTelemetry receivers, processors, exporters, and extensions that are compiled into this custom collector binary (`rungmpcol`).
- If you need to upgrade dependencies or add/remove collector components, edit `spec.yaml` and regenerate the collector code.

### 2. `confgenerator`
The sidecar is designed to scrape metrics using dynamically generated configurations.
- It parses a unified configuration schema (like `default-config.yaml` containing the `RunMonitoring` custom resource spec).
- It generates the standard Prometheus scrape configs and relabeling rules dynamically at runtime.
- It outputs the generated OTel configuration to `/run/rungmp/otel.yaml`.

### 3. `entrypoint.go`
The entrypoint acts as the wrapper process (PID 1) in the container:
1. Calls `confgenerator` to parse the `RunMonitoring` config (either default or custom config mounted from Secret Manager).
2. Generates the OTel config file `/run/rungmp/otel.yaml`.
3. Starts the custom OTel collector (`rungmpcol`) process pointing to the generated config file.

---

## Code Generation & Build Commands

Always use the following `make` targets to build and test:

*   **`make gen`** / **`make regen`**: Regenerates the collector code (`collector/` subdirectory) using `distrogen` based on `spec.yaml`. Run this after editing `spec.yaml`.
*   **`make build`**: Compiles both the wrapper entrypoint (`bin/run-gmp-entrypoint`) and the custom collector binary (`bin/rungmpcol`).
*   **`make test`**: Runs unit tests.
*   **`make test_update`**: Updates the unit test golden files inside `confgenerator/testdata/` (e.g., when adding relabeling features or modifying generated configurations).
*   **`make docker-build-image`**: Compiles the container locally using `./run-gmp-sidecar/Dockerfile.build`.
*   **`make docker-push-image`**: Tags and pushes the local sidecar image to the registry.
*   **`make docker-build-and-push`**: Combined target to build and push the sidecar image.

---

## Testing Workflows & Gotchas

### Golden Files Verification
Unit tests in `confgenerator` verify the generated configuration against golden output files in `confgenerator/testdata/<test-case>/golden/otel.yaml`.
- When adding features or changing relabeling logic, make sure to run `make test_update` to regenerate these files, and review the diff to ensure correctness before committing.

### End-to-End (E2E) Testing in Google Cloud
To perform integration or End-to-End (E2E) testing on actual Cloud Run environments, follow the deployment and validation guides in [README.md](../../../README.md).

#### E2E Testing Guidelines:
1. **Target Both Environments:** Cloud Run Services and Cloud Run Worker Pools are separate, distinct deployment types. You **must** deploy and test both environments separately to verify E2E correctness.
2. **Environment Variable Prompts:** Before running E2E tests, prompt the user for the required environment variables:
   - `GCP_PROJECT`: The Google Cloud project ID.
   - `REGION`: The region to deploy to (e.g., `us-east1`).
3. **Optional Secret Config:** Prompt the user to confirm if they want to deploy with a custom sidecar configuration (which requires mounting a secret from GCP Secret Manager). If they wish to skip it, proceed with the default configuration using `default-config.yaml`.

After deploying the Service or Worker Pool, wait 1–2 minutes for the metrics ingestion loop to initialize. You can query the ingested metrics using the helper script:

```bash
# Query the custom app gauge metric
go run scripts/query_metrics.go -project=my-project-id -metric=prometheus.googleapis.com/foo_metric/gauge

# Query self observability uptime metric
go run scripts/query_metrics.go -project=my-project-id -metric=agent/uptime
```

### Resource and Metric Label Processing (groupbyattrs)
In the generated OpenTelemetry configuration (`otel.yaml`), the scrape jobs initially capture all target metadata (like `namespace` and `cluster`) as **metric labels** (data point attributes) during Prometheus scrape relabeling.

To ensure they are correctly mapped to Google Cloud Monitoring's `prometheus_target` monitored resource, the configuration uses the `groupbyattrs` processor:
- **`groupbyattrs`**: Promotes specific metric labels (e.g., `namespace` and `cluster`) to **resource attributes**.
- **`transform`**: Using OTTL statements, it copies/renames attributes (like mapping `faas.instance` to `instanceId` and resetting `gcp.project.id` from `project_id`) to cleanly separate what ends up as a Monitored Resource Label vs. what remains as a Metric Label in GCM.

This pipeline behavior can be reviewed in the test golden files:
- [confgenerator/testdata/builtin/golden/otel.yaml](../../../confgenerator/testdata/builtin/golden/otel.yaml) (Services)
- [confgenerator/testdata/builtin-workerpool/golden/otel.yaml](../../../confgenerator/testdata/builtin-workerpool/golden/otel.yaml) (Worker Pools)

### Secrets and Custom Configurations
The collector can read a custom configuration file mounted from GCP Secret Manager (mapped to `/etc/rungmp/config.yaml`).
- To test custom configs in local/GCP deployments, create the secret first:
  ```bash
  gcloud secrets create run-gmp-config --data-file=default-config.yaml
  ```
- Reference the secret in your YAML using:
  ```yaml
  run.googleapis.com/secrets: "run-gmp-config:projects/<PROJECT_ID>/secrets/run-gmp-config"
  ```

### CPU Throttling in Worker Pools vs. Services
- In standard **Cloud Run Services**, CPU throttling is enabled by default. Background tasks (like OTel scrapers) require setting `run.googleapis.com/cpu-throttling: "false"` so that CPU is allocated even when no requests are being processed.
- In **Cloud Run Worker Pools**, CPU throttling does not exist since there are no incoming requests. The platform always allocates CPU. Setting `run.googleapis.com/cpu-throttling: "true"` or `"false"` is ignored on Worker Pools, but keeping it out of `run-workerpool.yaml` keeps configuration clean.

### Monitored Resource Fallbacks & Worker Pool Label Omissions
- **Required Monitored Resource Labels:** GCM's `prometheus_target` and underlying `cloud_run_revision` resource types define `configuration_name` as a required field.
- **Service Mapping:** In Cloud Run Services, `configuration_name` maps to the Knative configuration name.
- **Worker Pool Fallbacks & Omissions:**
  - Because Cloud Run Worker Pools do not have service or Knative configuration resources under the hood, the `service` (or `service_name`) and `configuration_name` metric labels are **omitted entirely** from the target scraping configurations.
  - The `namespace` label automatically falls back to using the worker pool name.
  - To prevent schema write errors on GCM backend, the Google Cloud exporter automatically resolves and maps any missing `configuration_name` to the worker pool name (`faas.name`).

### Final Exported Metric & Resource Labels

Telemetry sent to Google Cloud Managed Service for Prometheus is mapped to the `prometheus_target` monitored resource type. The table below represents the final structure of Resource Labels and Metric Labels recorded in Google Cloud Monitoring for each environment.

#### 1. Cloud Run Services

```yaml
Resource Type: prometheus_target
Resource Labels:
  project_id: "my-project-id"
  location: "us-central1"
  cluster: "__run__"
  namespace: "my-service" # Service name acts as the namespace
  job: "run-gmp-sidecar"
  instance: "abc123xyz:8080" # Instance ID + Scrape Port
Metric Labels:
  instanceId: "abc123xyz"
  revision_name: "my-service-00001-abc"
  service_name: "my-service"
  configuration_name: "my-service"
```

#### 2. Cloud Run Worker Pools

```yaml
Resource Type: prometheus_target
Resource Labels:
  project_id: "my-project-id"
  location: "us-central1"
  cluster: "__run__"
  namespace: "my-worker-pool" # Worker Pool name acts as the namespace
  job: "run-gmp-sidecar"
  instance: "abc123xyz:8080" # Instance ID + Scrape Port
Metric Labels:
  instanceId: "abc123xyz"
  revision_name: "my-worker-pool-00001-xyz"
  worker_pool: "my-worker-pool"
  # service_name and configuration_name are omitted entirely
```
