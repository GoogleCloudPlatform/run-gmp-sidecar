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
REGION="us-east1"

gcloud run services delete run-gmp-sidecar-service --region ${REGION} --quiet
# Delete secret if we created it before
if gcloud secrets list --filter="name ~ .*run-gmp-config.*" | grep run-gmp-sidecar
then
  gcloud secrets delete run-gmp-config
fi
gcloud artifacts repositories delete run-gmp \
  --location=${REGION} \
  --quiet
gcloud iam service-accounts delete ${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com
