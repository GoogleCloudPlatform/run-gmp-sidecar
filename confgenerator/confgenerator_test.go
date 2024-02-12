// Copyright 2023 Google LLC
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

package confgenerator_test

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/run-gmp-sidecar/confgenerator"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/golden"
)

const (
	builtinTestdataDirName = "builtin"
	goldenDir              = "golden"
	errorGolden            = goldenDir + "/error"
	inputFileName          = "input.yaml"
)

func testMetadata() *confgenerator.CloudRunEnvironment {
	return &confgenerator.CloudRunEnvironment{
		Service:       "test_service",
		Revision:      "test_revision",
		Configuration: "test_configuration",
	}
}

func TestGoldens(t *testing.T) {
	t.Parallel()
	testNames := getTestsInDir(t)

	for _, testName := range testNames {
		// https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		testName := testName
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			t.Run(testName, func(t *testing.T) {
				testDir := filepath.Join(testName)
				got, err := generateConfigs(testDir)
				if strings.HasPrefix(testName, "invalid-") {
					assert.Assert(t, err != nil, "expected test config to be invalid, but was successful")
				}
				if err := testGeneratedFiles(t, got, filepath.Join(testDir, goldenDir)); err != nil {
					t.Errorf("Failed to check generated configs: %v", err)
				}
			})
		})
	}
}

func getTestsInDir(t *testing.T) []string {
	t.Helper()

	testdataDir := filepath.Join("testdata")
	testDirEntries, err := os.ReadDir(testdataDir)
	if os.IsNotExist(err) {
		// No tests for this combination.
		return nil
	}
	assert.NilError(t, err, "couldn't read directory %s: %v", testdataDir, err)
	testNames := []string{}
	for _, testDirEntry := range testDirEntries {
		if !testDirEntry.IsDir() {
			continue
		}
		userSpecifiedConfPath := filepath.Join(testdataDir, testDirEntry.Name(), inputFileName)
		if _, err := os.Stat(userSpecifiedConfPath + ".missing"); err == nil {
			// Intentionally missing
		} else if _, err := os.Stat(userSpecifiedConfPath); errors.Is(err, os.ErrNotExist) {
			// Empty directory; probably a leftover with backup files.
			continue
		}
		testNames = append(testNames, testDirEntry.Name())
	}
	return testNames
}

func generateConfigs(testDir string) (got map[string]string, err error) {
	ctx := context.Background()

	got = make(map[string]string)
	defer func() {
		if err != nil {
			got["error"] = err.Error()
		}
	}()

	c, err := confgenerator.ReadConfigFromFile(ctx, filepath.Join("testdata", testDir, inputFileName))
	if err != nil {
		return
	}

	// Use deterministic metadata and self metrics port for tests
	c.Env = testMetadata()
	selfMetricsPort := 42

	// Otel configs
	otelGeneratedConfig, err := c.GenerateOtelConfig(ctx, selfMetricsPort)
	if err != nil {
		return
	}
	got["otel.yaml"] = otelGeneratedConfig
	return
}

func testGeneratedFiles(t *testing.T, generatedFiles map[string]string, testDir string) error {
	// Find all files currently in this test directory
	existingFiles := map[string]struct{}{}
	goldenPath := filepath.Join("testdata", testDir)
	err := filepath.Walk(
		goldenPath,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Mode().IsRegular() {
				existingFiles[info.Name()] = struct{}{}
			}
			return nil
		},
	)
	if golden.FlagUpdate() && os.IsNotExist(err) {
		if err := os.Mkdir(goldenPath, 0777); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Assert the goldens of all the generated files. Either the generated file
	// matches a file already present in the directory, or the file is new.
	// If the file is new, the test will fail if not currently doing a golden
	// update (`-update` flag).
	for file, content := range generatedFiles {
		golden.Assert(t, content, filepath.Join(testDir, file))
		delete(existingFiles, file)
	}

	// If there are any files left in the existing file map, then that means the
	// test generated new files and we're currently in an update run. We now need
	// to clean up the existing lua files left aren't being generated anymore.
	for file := range existingFiles {
		if golden.FlagUpdate() {
			err := os.Remove(filepath.Join("testdata", testDir, file))
			if err != nil {
				return err
			}
		} else {
			t.Errorf("unexpected existing file: %q", file)
		}
	}

	return nil
}
