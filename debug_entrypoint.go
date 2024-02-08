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
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Create channel to listen for signals.
var signalChan chan (os.Signal) = make(chan os.Signal, 1)
var otelConfigFile = "otel.yaml"
var configRefreshInterval = 10 * time.Second
var debugPort = 41285

func main() {
	// SIGINT handles Ctrl+C locally.
	// SIGTERM handles Cloud Run termination signal.
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Spin up new-subprocess that runs the OTel collector and store the PID.
	// This OTel collector should use the generated config.
	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{nil, /* stdin is not needed for the collector */
		os.Stdout, os.Stderr}
	collectorProcess, err := os.StartProcess("./bin/rungmpcol", []string{"./bin/rungmpcol", "--config", otelConfigFile}, &procAttr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("entrypoint: started OTel successfully")

	refreshTicker := time.NewTicker(configRefreshInterval)
	debugTicker := time.NewTicker(35 * time.Second)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case <-refreshTicker.C:
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
		case <-debugTicker.C:
			// make a get request to localhost:2020/metrics and print the result.
			url := fmt.Sprintf("http://localhost:%d/metrics", debugPort)
			req, _ := http.NewRequest("GET", url, nil)
			res, _ := http.DefaultClient.Do(req)
			body, _ := ioutil.ReadAll(res.Body)
			log.Println("COLLECTED METRICS: ", string(body))
			res.Body.Close()
		}
	}
}
