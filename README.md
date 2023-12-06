# Run GMP Sidecar

This sidecar is Google's recommended way to get GMP/Prometheus styled monitoring
for Cloud Run services. It is powered by GMP on the server side and OpenTelemetry
on the client side. This uses the **Cloud Run
multicontainer (sidecar) feature** to run the Collector as a sidecar container
alongside your workload container.

[Learn more about sidecars in Cloud Run (currently a preview feature)
here.](https://cloud.google.com/run/docs/deploying#multicontainer)

## Getting started

The following steps walk you through setting up a sample app on Cloud Run that
exports your application's Prometheus metrics to GMP.

### Prerequisites

* Enable the Cloud Run API in your Google Cloud project.
* Enable the Artifact Registry API in your Google Cloud project.
* Authenticate with your Google Cloud project in a local terminal.

To enable the depending service APIs with `gcloud` command, you can the following commands.

```console
gcloud services enable run.googleapis.com --quiet
gcloud services enable artifactregistry.googleapis.com --quiet
gcloud services enable monitoring.googleapis.com --quiet
```

To run the sample app, you will also need to make sure your [Cloud Run Service
Account](https://cloud.google.com/run/docs/configuring/service-accounts) has, at
minimum, the following IAM roles:

* `roles/monitoring.metricWriter`
* `roles/logging.logWriter`

The default Compute Engine Service Account has these roles already.

Export several environment variables to control the project, region and secret
name to use.
```
export GCP_PROJECT=<project-id>
export REGION=us-east1
export RUN_GMP_CONFIG=run-gmp-config
```

### Run sample (automated)

Because this sample requires `docker` or similar container build system for Linux runtime, you can use Cloud Build when you are trying without local Docker support. To enable Cloud Build, you need to enable Cloud Build API in your Google Cloud project.

```console
gcloud services enable cloudbuild.googleapis.com --quiet
```

The bundled configuration file for Cloud Build (`cloudbuild-simple.yaml`) requires a new service account with the following roles or stronger:

* `roles/iam.serviceAccountUser`
* `roles/storage.objectViewer`
* `roles/monitoring.metricWriter`
* `roles/logging.logWriter`
* `roles/artifactregistry.createOnPushWriter`
* `roles/run.admin`
* `roles/secretmanager.admin` (Needed for custom configs only)
* `roles/secretmanager.secretAccessor`(Needed for custom configs only)

Running `create-sa-and-ar.sh` creates a new service account `run-gmp-sa@<project-id>.iam.gserviceaccount.com` for you, and an Artifact Registry repo for the images. Then launch a Cloud Build task with `gcloud` command.

```console
./create-sa-and-ar.sh
gcloud builds submit . --config=cloudbuild-simple.yaml --region=${REGION}
```

> **_NOTE:_**  If you have an Org policy that prevents unauthenticated access, then you might see a failure in the final step. You can safely ignore this failure.

After the build, run the following command to check the endpoint URL.

```console
gcloud run services describe my-cloud-run-service --region=${REGION} --format="value(status.url)"
```

### Run sample (manual steps)

#### Build the sample app

The `app` directory contains a sample app written in Go. This app generates some
simple Prometheus metrics (a gauge and a counter).

Create an Artifact Registry container image repository with the following
commands:

```
gcloud artifacts repositories create run-gmp \
    --repository-format=docker \
    --location=${REGION}
```

Authenticate your Docker client with `gcloud`:

```
gcloud auth configure-docker \
    ${REGION}-docker.pkg.dev
```

Build and push the app with the following commands:

```
pushd sample-apps/simple-app
docker build -t ${REGION}-docker.pkg.dev/$GCP_PROJECT/run-gmp/sample-app .
docker push ${REGION}-docker.pkg.dev/$GCP_PROJECT/run-gmp/sample-app
popd
```

#### Build the Collector image

The `collector` directory contains a Dockerfile and OpenTelemetry Collector
config file. The Dockerfile builds a Collector image that bundles the local
config file with it.

Build the Collector image with the following commands:

```
docker build -t ${REGION}-docker.pkg.dev/$GCP_PROJECT/run-gmp/collector .
docker push ${REGION}-docker.pkg.dev/$GCP_PROJECT/run-gmp/collector
```

#### Create the Cloud Run Service (default config)

The `run-service-simple.yaml` file defines a multicontainer Cloud Run Service with the
sample app and Collector images built above. This will run with the default config, which scrapes an application emitting metrics at port `8080` at the path `/metrics`.

Replace the `%SAMPLE_APP_IMAGE%` and `%OTELCOL_IMAGE%` placeholders in
`run-service-simple.yaml` with the images you built above, ie:

```
sed -i s@%OTELCOL_IMAGE%@${REGION}-docker.pkg.dev/${GCP_PROJECT}/run-gmp/collector@g run-service-simple.yaml
sed -i s@%SAMPLE_APP_IMAGE%@${REGION}-docker.pkg.dev/${GCP_PROJECT}/run-gmp/sample-app@g run-service-simple.yaml
```

Create the Service with the following command:

```
gcloud run services replace run-service-simple.yaml --region=${REGION}
```

This command will return an external URL for your Service’s endpoint. Save this
and use it in the next section to trigger the sample app so you can see the
telemetry collected by OpenTelemetry.

#### Create the Cloud Run Service (custom config)

##### Create RunMonitoring config and store as a secret

Create a `RunMonitoring` config and store it in secret manager. In this example,
we use `run-gmp-config` as the secret name. The file we're using is
`default-config.yaml` and it scrapes the main container at port `8080` using the
path `/metrics` every 30s. You can replace this with any `RunMonitoring` config
file that you want the sidecar to use.

```
gcloud secrets create ${RUN_GMP_CONFIG}  --data-file=default-config.yaml
```

##### Deploy the service

The `run-service.yaml` file defines a multicontainer Cloud Run Service with the
sample app and Collector images built above, using the config you placed in secret manager.

Replace the `%SAMPLE_APP_IMAGE%`, `%OTELCOL_IMAGE%`, `%PROJECT%` and `%SECRET%`
placeholders in `run-service.yaml` with the images you built above, ie:

```
sed -i s@%OTELCOL_IMAGE%@${REGION}-docker.pkg.dev/${GCP_PROJECT}/run-gmp/collector@g run-service.yaml
sed -i s@%SAMPLE_APP_IMAGE%@${REGION}-docker.pkg.dev/${GCP_PROJECT}/run-gmp/sample-app@g run-service.yaml
sed -i s@%PROJECT%@${GCP_PROJECT}@g run-service.yaml
sed -i s@%SECRET%@${RUN_GMP_CONFIG}@g run-service.yaml
```

Create the Service with the following command:

```
gcloud run services replace run-service.yaml --region=${REGION}
```

This command will return an external URL for your Service’s endpoint. Save this
and use it in the next section to trigger the sample app so you can see the
telemetry collected by OpenTelemetry.

#### Allow unauthenticated HTTP access

Finally before you make make the request to the URL, you need to change
the Cloud Run service policy to accept unauthenticated HTTP access.

```
gcloud run services set-iam-policy my-cloud-run-service policy.yaml --region=${REGION}
```

> **_NOTE:_**  If you have an Org policy that prevents unauthenticated access, then this step will fail. But fear not, you can simply curl the endpoint using `curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" <ENDPOINT>` instead.

### View telemetry in Google Cloud

Use `curl` to make a request to your Cloud Run Service’s endpoint URL:

```
export SERVICE_URL=<service-url>
curl $SERVICE_URL
```

This should return the following output on success:

```
User request received!
```

> **_NOTE:_**  If you get permission errors because of unauthenticated access, then this will fail. But fear not, you can simply curl the endpoint using `curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" $SERVICE_URL` instead.

You should now be able to use Cloud Monitoring to find metrics from the application. The `app` container emits the following metrics:
- `foo_metric`: A `gauge` metric that emits the current time as a float
- `bar_metric`: A `counter` metric that emits the current time as a float

#### Troubleshooting the sidecar

The sidecar reports self metrics and self logs to Cloud Monitoring and Cloud Logging respectively.

##### Self observability metrics
You should also check out the sidecar's self metrics:
- `agent_uptime`: Uptime of the sidecar collector
- `agent_memory_usage`: Memory in use by the sidecar collector
- `agent_api_request_count`: Count of API requests from the sidecar collector
- `agent_monitoring_point_count`: Count of metric points written by the agent to Cloud Monitoring by the sidecar collector

Querying these metrics using the Google Cloud Monitoring UI is left as an
exercise for the reader. Be sure to check out the resource and metric labels for
added homework.

##### Self observability logs
Logs from the sidecar are written against the `Cloud Run Revision` [monitored resource](https://cloud.google.com/monitoring/api/resources#tag_cloud_run_revision) in Cloud Logging.

### Clean up

After running the demo, please make sure to clean up your project so that you don't consume unexpected resources and get charged.

```console
./clean-up-cloud-run.sh
```
