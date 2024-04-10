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

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/GoogleCloudPlatform/run-gmp-sidecar/confgenerator"
)

// Create channel to listen for signals.
var signalChan chan (os.Signal) = make(chan os.Signal, 1)
var userConfigFile = "/etc/rungmp/config.yaml"
var otelConfigFile = "/tmp/rungmp/otel.yaml"
var configRefreshInterval = 20 * time.Second
var selfMetricsPort = 0

func getRawUserConfig(userConfigFile string) (string, error) {
	_, err := os.Stat(userConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("failed to stat file %q: %v", userConfigFile, err)
	}

	data, err := ioutil.ReadFile(userConfigFile)
	if err != nil {
		return "", fmt.Errorf("failed to read file %q: %v", userConfigFile, err)
	}

	return string(data), nil
}

// Generate OTel config from RunMonitoring config. Returns an error if
// generation of OTel configs failed.
func generateOtelConfig(ctx context.Context, userConfigFile string) error {
	// Pick up RunMonitoring configuration from mounted volume that is tied to
	// secret manager.  Translate it from RunMonitoring to OTel.
	c, err := confgenerator.ReadConfigFromFile(ctx, userConfigFile, os.Getenv("GMP_CFG_SPEC"), os.Getenv("GMP_SCRAPE_PORTS"))
	if err != nil {
		log.Fatal(err)
	}

	if selfMetricsPort == 0 {
		selfMetricsPort, err = confgenerator.GetFreePort()
		if err != nil {
			return err
		}
	}

	// Create the OTel config and write it to disk
	otel, err := c.GenerateOtelConfig(ctx, selfMetricsPort)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(otelConfigFile), 0755); err != nil {
		return fmt.Errorf("failed to create directory for %q: %v", otelConfigFile, err)
	}
	if err := ioutil.WriteFile(otelConfigFile, []byte(otel), 0644); err != nil {
		return fmt.Errorf("failed to write file to %q: %v", otelConfigFile, err)
	}

	return nil
}

func main() {
	// SIGINT handles Ctrl+C locally.
	// SIGTERM handles Cloud Run termination signal.
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	ctx := context.Background()

	lastRawConfig, err := getRawUserConfig(userConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	// Generate the OTel config for the first time.
	err = generateOtelConfig(ctx, userConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	// Spin up new-subprocess that runs the OTel collector and store the PID.
	// This OTel collector should use the generated config.
	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{nil, /* stdin is not needed for the collector */
		os.Stdout, os.Stderr}
	collectorProcess, err := os.StartProcess("./rungmpcol", []string{"./rungmpcol", "--config", otelConfigFile}, &procAttr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("entrypoint: started OTel successfully")

	refreshTicker := time.NewTicker(configRefreshInterval)
	for {
		select {
		case <-refreshTicker.C:
			rawConfig, err := getRawUserConfig(userConfigFile)
			if err != nil {
				log.Fatal(err)
			}

			// Check if we're using the default config. Only reload if something
			// has changed since the last time we checked.
			if rawConfig == lastRawConfig {
				continue
			}

			// Something changed since the last time we checked the config.
			err = generateOtelConfig(ctx, userConfigFile)
			if err != nil {
				log.Fatal(err)
			}
			lastRawConfig = rawConfig

			// Signal the OTel collector to reload its config
			collectorProcess.Signal(syscall.SIGHUP)
			log.Println("entrypoint: reloaded OTel config")
		case sig := <-signalChan:
			// Wait for signals from Cloud Run. Signal the sub process appropriately
			// after making relevant changes to the config and/or health signals.
			// TODO(b/307317433): Consider having a timeout to shutdown the subprocess
			// non-gracefully.
			log.Printf("entrypoint: %s signal caught", sig)

			collectorProcess.Signal(sig)
			processState, err := collectorProcess.Wait()
			if err != nil {
				log.Fatalf(processState.String(), err)
			}
			log.Print("entrypoint: sidecar exited")
			return
		}
	}
}
