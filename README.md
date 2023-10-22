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
exports your applciations prometheus metrics to GMP.

### Prerequisites

* Enable the Cloud Run API in your Google Cloud project.
* Enable the Artifact Registry API in your Google Cloud project.
* Authenticate with your Google Cloud project in a local terminal.

To enable the depending service APIs with `gcloud` command, you can the following commands.

```console
gcloud services enable run.googleapis.com --quiet
gcloud services enable artifactregistry.googleapis.com --quiet
gcloud services enable cloudtrace.googleapis.com --quiet
gcloud services enable monitoring.googleapis.com --quiet
```

To run the sample app, you will also need to make sure your [Cloud Run Service
Account](https://cloud.google.com/run/docs/configuring/service-accounts) has, at
minimum, the following IAM roles:

* `roles/monitoring.metricWriter`
* `roles/cloudtrace.agent`
* `roles/logging.logWriter`

The default Compute Engine Service Account has these roles already.

### Run sample

#### Cloud Build

Because this sample requires `docker` or similar container build system for Linux runtime, you can use Cloud Build when you are trying without local Docker support. To enable Cloud Build, you need to enable Cloud Build API in your Google Cloud project.

```console
gcloud services enable cloudbuild.googleapis.com --quiet
```

The bundled configuration file for Cloud Build (`cloudbuild.yaml`) requires a new servcie account with the following roles or stronger:

* `roles/iam.serviceAccountUser`
* `roles/storage.objectViewer`
* `roles/logging.logWriter`
* `roles/artifactregistry.createOnPushWriter`
* `roles/run.admin`

Running `create-service-account.sh` creates a new service account `run-gmp-sa@<project-id>.iam.gserviceaccount.com` for you. Then launch a Cloud Build task with `gcloud` command.

```console
./create-service-account.sh
gcloud builds submit . --config=cloudbuild.yaml
```

After the build, run the following command to check the endpoint URL.

```console
gcloud run services describe run-gmp-sidecar-service --region=us-east1 --format="value(status.url)"
```

#### Build and Run Manually

##### Build the sample app

The `app` directory contains a sample app written in Go. This app generates some
simple prometheus metrics (a gauge and a counter).

Create an Artifact Registry container image repository with the following
commands:

```
export GCP_PROJECT=<project-id>
gcloud artifacts repositories create run-gmp \
    --repository-format=docker \
    --location=us-east1
```

Authenticate your Docker client with `gcloud`:

```
gcloud auth configure-docker \
    us-east1-docker.pkg.dev
```

Build and push the app with the following commands:

```
pushd sample-apps/simple-app
docker build -t us-east1-docker.pkg.dev/$GCP_PROJECT/run-gmp/sample-app .
docker push us-east1-docker.pkg.dev/$GCP_PROJECT/run-gmp/sample-app
popd
```

##### Build the Collector image

The `collector` directory contains a Dockerfile and OpenTelemetry Collector
config file. The Dockerfile builds a Collector image that bundles the local
config file with it.

Build the Collector image with the following commands:

```
docker build -t us-east1-docker.pkg.dev/$GCP_PROJECT/run-gmp/collector .
docker push us-east1-docker.pkg.dev/$GCP_PROJECT/run-gmp/collector
```

##### Create the Cloud Run Service

The `run-service.yaml` file defines a multicontainer Cloud Run Service with the
sample app and Collector images built above.

Replace the `%SAMPLE_APP_IMAGE%` and `%OTELCOL_IMAGE%` placeholders in
`run-service.yaml` with the images you built above, ie:

```
sed -i s@%OTELCOL_IMAGE%@us-east1-docker.pkg.dev/${GCP_PROJECT}/run-gmp/collector@g run-service.yaml
sed -i s@%SAMPLE_APP_IMAGE%@us-east1-docker.pkg.dev/${GCP_PROJECT}/run-gmp/sample-app@g run-service.yaml
```

Create the Service with the following command:

```
gcloud run services replace run-service.yaml
```

This command will return an external URL for your Service’s endpoint. Save this
and use it in the next section to trigger the sample app so you can see the
telemetry collected by OpenTelemetry.

Finally before you make make the request to the URL, you need to change
the Cloud Run service policy to accept unauthenticated HTTP access.

```
gcloud run services set-iam-policy run-gmp-sidecar-service policy.yaml
```

### View telemetry in Google Cloud

Use `curl` to make a request to your Cloud Run Service’s endpoint URL:

```
export SERVICE_URL=<service-url>
curl $SERVICE_URL/metrics
```

This should return the following output on success:

```
User request received!
```

### Clean up

After running the demo, please make sure to clean up your project so that you don't consume unexpected resources and get charged.

```console
gcloud run services delete run-gmp-sidecar-service --region us-east1 --quiet
gcloud artifacts repositories delete run-gmp \
  --location=us-east1 \
  --quiet
```
