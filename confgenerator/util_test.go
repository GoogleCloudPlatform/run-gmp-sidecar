// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package confgenerator

import (
	"os"
	"testing"
)

func TestFetchMetadata(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		expected *CloudRunEnvironment
	}{
		{
			name: "Standard Cloud Run service",
			env: map[string]string{
				"K_SERVICE":       "my-service",
				"K_REVISION":      "my-service-rev1",
				"K_CONFIGURATION": "my-config",
			},
			expected: &CloudRunEnvironment{
				Service:       "my-service",
				Revision:      "my-service-rev1",
				Configuration: "my-config",
			},
		},
		{
			name: "Cloud Run Worker Pool",
			env: map[string]string{
				"CLOUD_RUN_WORKER_POOL": "my-worker-pool",
				"CLOUD_RUN_REVISION":    "my-worker-pool-rev1",
			},
			expected: &CloudRunEnvironment{
				WorkerPool: "my-worker-pool",
				Revision:   "my-worker-pool-rev1",
			},
		},

	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Clear env first
			os.Unsetenv("K_SERVICE")
			os.Unsetenv("K_REVISION")
			os.Unsetenv("K_CONFIGURATION")
			os.Unsetenv("CLOUD_RUN_WORKER_POOL")
			os.Unsetenv("CLOUD_RUN_REVISION")

			for k, v := range tc.env {
				os.Setenv(k, v)
			}

			got := fetchMetadata()
			if got.Service != tc.expected.Service {
				t.Errorf("Service: got %q, want %q", got.Service, tc.expected.Service)
			}
			if got.Revision != tc.expected.Revision {
				t.Errorf("Revision: got %q, want %q", got.Revision, tc.expected.Revision)
			}
			if got.Configuration != tc.expected.Configuration {
				t.Errorf("Configuration: got %q, want %q", got.Configuration, tc.expected.Configuration)
			}
			if got.WorkerPool != tc.expected.WorkerPool {
				t.Errorf("WorkerPool: got %q, want %q", got.WorkerPool, tc.expected.WorkerPool)
			}
		})
	}
}
