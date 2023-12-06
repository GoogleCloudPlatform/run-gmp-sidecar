#!/bin/bash
# Copyright 2023 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PROJECT_ID=$(gcloud config get-value project)
SA_NAME="run-gmp-sa"
REGION="${REGION:-us-east1}"
SERVICE_NAME="my-cloud-run-service"
SECRET="run-gmp-config"
REPO="run-gmp"
SA="${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"

# Delete the service
if gcloud run services list --project=${PROJECT_ID} --region=${REGION} | grep ${SERVICE_NAME}
then
  gcloud run services delete ${SERVICE_NAME} --region ${REGION} --quiet
fi

# Delete secret if we created it before.
if gcloud secrets list --filter="name ~ .*${SECRET}.*" | grep ${SECRET}
then
  gcloud secrets delete ${SECRET}
fi

# Delete AR repo if it exists.
if gcloud artifacts repositories list --project=${PROJECT_ID} --location=${REGION} | grep ${REPO}
then
  gcloud artifacts repositories delete ${REPO} \
    --location=${REGION} \
    --quiet
fi

# Delete SA if exists.
if gcloud iam service-accounts list --project=${PROJECT_ID} | grep ${SA}
then
  gcloud iam service-accounts delete ${SA}
fi
