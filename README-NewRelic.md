# New Relic Integration for run-gmp-sidecar

This repository has been adapted to send metrics to New Relic instead of Google Cloud Monitoring, while maintaining all the existing functionality and patterns of the original Google GMP Sidecar.

## What Changed

### Key Modifications

1. **New Relic Exporter**: Added a custom OpenTelemetry exporter (`collector/exporter/newrelicexporter`) that formats metrics for New Relic's Metric API
2. **Configuration Update**: Modified the config generator to use the New Relic exporter instead of Google Managed Prometheus
3. **Environment Variables**: Added support for New Relic API key and service metadata via environment variables
4. **Deployment Scripts**: Created new deployment scripts and configurations for New Relic setup

### Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Your App     │    │  Sidecar         │    │   New Relic     │
│                │    │                  │    │                 │
│  Port 8080     │───▶│  Prometheus      │───▶│   Metric API    │
│  /metrics      │    │  Receiver        │    │                 │
│                │    │                  │    │                 │
└─────────────────┘    │  New Relic       │    └─────────────────┘
                       │  Exporter        │
                       └──────────────────┘
```

## Quick Start

### Prerequisites

1. **New Relic Account**: You need a New Relic account and an Insights Insert API key
2. **Google Cloud**: GCP project with Cloud Run and Artifact Registry APIs enabled
3. **Docker**: For building container images

### 1. Get Your New Relic API Key

1. Log into your New Relic account
2. Go to **Settings** → **API Keys**
3. Create or copy an **Insights Insert API Key**
4. Set it as an environment variable:
   ```bash
   export NEW_RELIC_API_KEY="your-api-key-here"
   ```

### 2. Deploy to Cloud Run

```bash
# Set your GCP project
export GCP_PROJECT="your-project-id"
export REGION="us-east1"

# Run the setup script
./setup-newrelic-simple.sh
```

The setup script will:
- Create necessary secrets in Google Secret Manager
- Build and push container images
- Deploy your service to Cloud Run
- Configure the New Relic integration

### 3. Test Your Deployment

```bash
# Get your service URL
SERVICE_URL=$(gcloud run services describe my-cloud-run-service-newrelic --region=$REGION --format="value(status.url)")

# Test the endpoint
curl $SERVICE_URL
```

### 4. View Metrics in New Relic

1. Log into your New Relic account
2. Navigate to **Metrics & Events**
3. Look for metrics from your service (they should appear within a few minutes)

## Configuration

### Environment Variables

The New Relic exporter supports these environment variables:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `NEW_RELIC_API_KEY` | ✅ | - | Your New Relic Insights Insert API key |
| `SERVICE_NAME` | ❌ | `run-gmp-sidecar` | Name of your service |
| `SERVICE_VERSION` | ❌ | `1.0.0` | Version of your service |
| `ENVIRONMENT` | ❌ | `production` | Deployment environment |

### Volume Mount NOT Required!

Unlike many monitoring solutions, this New Relic integration works with **no volume mounts needed**! The sidecar uses a default configuration that scrapes `0.0.0.0:8080/metrics` every 30 seconds. Perfect for simple deployments.

## Metrics Format

### Prometheus to New Relic Mapping

The exporter automatically converts Prometheus metrics to New Relic format:

| Prometheus Type | New Relic Type | Notes |
|----------------|----------------|-------|
| Gauge | `gauge` | Direct mapping |
| Counter | `count` | Cumulative counters |
| Histogram | `distribution` | With bucket data |
| Summary | `summary` | With quantile data |

### Example Metrics

Your application metrics will appear in New Relic with these attributes:

```json
{
  "name": "foo_metric",
  "type": "gauge", 
  "value": 123.45,
  "timestamp": 1640995200,
  "attributes": {
    "service.name": "my-cloud-run-service",
    "service.version": "1.0.0",
    "deployment.environment": "production",
    "job": "my-app",
    "instance": "localhost:8080"
  }
}
```

## Monitoring and Troubleshooting

### Check Logs

```bash
# View sidecar logs
gcloud run services logs tail my-cloud-run-service-newrelic --region=$REGION

# Filter for New Relic exporter logs
gcloud run services logs tail my-cloud-run-service-newrelic --region=$REGION | grep newrelic
```

### Verify Configuration

```bash
# Check if API key secret exists
gcloud secrets describe newrelic-secrets

# View the secret value (be careful!)
gcloud secrets versions access latest --secret=newrelic-secrets
```

### Common Issues

1. **No metrics appearing in New Relic**
   - Verify API key is correct
   - Check service logs for authentication errors
   - Ensure your app is exposing metrics on `/metrics`

2. **High memory usage**
   - Increase memory limits in `run-service-newrelic-simple.yaml`
   - Reduce scraping frequency in your config

3. **Authentication errors**
   - Verify the API key has Insights Insert permissions
   - Check that the secret is properly mounted

## Migration from GMP

If you're migrating from the original Google Managed Prometheus setup:

1. **No changes needed** to your application code or metrics format
2. **Keep existing** RunMonitoring configurations
3. **Update environment variables** to use New Relic instead of GCP
4. **Change dashboards** from Cloud Monitoring to New Relic

### Side-by-Side Comparison

| Feature | Original GMP | New Relic Version |
|---------|-------------|-------------------|
| Prometheus scraping | ✅ | ✅ (unchanged) |
| Custom configs | ✅ | ✅ (unchanged) |
| Sidecar pattern | ✅ | ✅ (unchanged) |
| Cloud Run integration | ✅ | ✅ (unchanged) |
| Volume mounts | Required | **NOT Required!** |
| Destination | Google Cloud Monitoring | New Relic |
| API Key | GCP Service Account | New Relic Insert Key |

## License

This project maintains the same Apache 2.0 license as the original Google run-gmp-sidecar project.
