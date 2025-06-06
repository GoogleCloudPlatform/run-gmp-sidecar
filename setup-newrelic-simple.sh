#!/bin/bash

# Simple setup script for New Relic deployment (no volume mounts!)
set -e

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
PROJECT_ID=${GCP_PROJECT:-$(gcloud config get-value project)}
REGION=${REGION:-us-east1}
SERVICE_NAME=${SERVICE_NAME:-my-cloud-run-service-newrelic}
SECRET_NAME=${SECRET_NAME:-newrelic-secrets}

echo -e "${GREEN}🚀 Setting up New Relic integration (simple version - no volumes!)${NC}"
echo "Project: $PROJECT_ID"
echo "Region: $REGION"
echo "Service: $SERVICE_NAME"
echo ""

# Check required environment variables
if [ -z "$NEW_RELIC_API_KEY" ]; then
    echo -e "${RED}❌ Error: NEW_RELIC_API_KEY environment variable is required${NC}"
    echo "Please set your New Relic Insights Insert API key:"
    echo "export NEW_RELIC_API_KEY='your-api-key-here'"
    exit 1
fi

# Enable required APIs
echo -e "${YELLOW}📋 Enabling required APIs...${NC}"
gcloud services enable run.googleapis.com --quiet
gcloud services enable artifactregistry.googleapis.com --quiet
gcloud services enable secretmanager.googleapis.com --quiet

# Create New Relic API key secret
echo -e "${YELLOW}🔐 Creating New Relic API key secret...${NC}"
if gcloud secrets describe $SECRET_NAME --project=$PROJECT_ID &>/dev/null; then
    echo "Secret $SECRET_NAME already exists, updating..."
    echo -n "$NEW_RELIC_API_KEY" | gcloud secrets versions add $SECRET_NAME --data-file=-
else
    echo "Creating new secret $SECRET_NAME..."
    echo -n "$NEW_RELIC_API_KEY" | gcloud secrets create $SECRET_NAME --data-file=-
fi

# Build and push images
echo -e "${YELLOW}🔨 Building and pushing images...${NC}"

# Create Artifact Registry repository if it doesn't exist
if ! gcloud artifacts repositories describe run-gmp --location=$REGION &>/dev/null; then
    echo "Creating Artifact Registry repository..."
    gcloud artifacts repositories create run-gmp \
        --repository-format=docker \
        --location=$REGION \
        --project=$PROJECT_ID
fi

# Authenticate Docker
gcloud auth configure-docker ${REGION}-docker.pkg.dev

# Build sample app
echo "Building sample app..."
pushd sample-apps/simple-app
docker build -t ${REGION}-docker.pkg.dev/$PROJECT_ID/run-gmp/sample-app .
docker push ${REGION}-docker.pkg.dev/$PROJECT_ID/run-gmp/sample-app
popd

# Build New Relic collector
echo "Building New Relic collector..."
docker build -f Dockerfile.newrelic -t ${REGION}-docker.pkg.dev/$PROJECT_ID/run-gmp/newrelic-collector .
docker push ${REGION}-docker.pkg.dev/$PROJECT_ID/run-gmp/newrelic-collector

# Update service configuration (simple version!)
echo -e "${YELLOW}📝 Updating service configuration...${NC}"
cp run-service-newrelic-simple.yaml run-service-deploy.yaml

sed -i "s@%SAMPLE_APP_IMAGE%@${REGION}-docker.pkg.dev/${PROJECT_ID}/run-gmp/sample-app@g" run-service-deploy.yaml
sed -i "s@%OTELCOL_IMAGE%@${REGION}-docker.pkg.dev/${PROJECT_ID}/run-gmp/newrelic-collector@g" run-service-deploy.yaml

# Deploy service
echo -e "${YELLOW}🚀 Deploying Cloud Run service...${NC}"
gcloud run services replace run-service-deploy.yaml --region=$REGION

# Make service publicly accessible (optional)
echo -e "${YELLOW}🌐 Making service publicly accessible...${NC}"
gcloud run services set-iam-policy $SERVICE_NAME policy.yaml --region=$REGION 2>/dev/null || echo "Warning: Could not set public access policy"

# Get service URL
SERVICE_URL=$(gcloud run services describe $SERVICE_NAME --region=$REGION --format="value(status.url)")

echo ""
echo -e "${GREEN}✅ Deployment complete!${NC}"
echo -e "${GREEN}🔗 Service URL: ${SERVICE_URL}${NC}"
echo ""
echo -e "${YELLOW}🧪 Test your deployment:${NC}"
echo "curl $SERVICE_URL"
echo ""
echo -e "${GREEN}📊 Default monitoring config:${NC}"
echo "  - Scrapes: 0.0.0.0:8080/metrics"
echo "  - Interval: 30 seconds"
echo "  - No custom config needed!"
echo ""
echo -e "${YELLOW}📊 Check your metrics in New Relic:${NC}"
echo "- Log into your New Relic account"
echo "- Navigate to Metrics & Events"
echo "- Look for metrics from your service"
echo ""
echo -e "${YELLOW}🔍 Troubleshooting:${NC}"
echo "- Check logs: gcloud run services logs tail $SERVICE_NAME --region=$REGION"
echo "- Verify API key: gcloud secrets versions access latest --secret=$SECRET_NAME"

# Clean up temporary file
rm -f run-service-deploy.yaml

echo -e "${GREEN}🎉 Setup complete! Super simple - no volumes needed!${NC}"
