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
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Create channel to listen for signals.
var signalChan chan (os.Signal) = make(chan os.Signal, 1)

func main() {
	// SIGINT handles Ctrl+C locally.
	// SIGTERM handles Cloud Run termination signal.
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// 1. Pick up RunMonitoring configuration from mounted volume that
	// is tied to secret manager.
	// TODO(b/293137197)

	// 2. Translate from RunMonitoring to OTel.
	// TODO(b/293137197)

	// 3. Spin up new-subprocess that runs the OTel collector and store the PID.
	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{nil, /* stdin is not needed for the collector */
		os.Stdout, os.Stderr}
	collectorProcess, err := os.StartProcess("./rungmpcol", []string{"./rungmpcol", "--config", "/etc/rungmp/config.yml"}, &procAttr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("entrypoint: started OTel successfully")

	// 4. Wait for signals from Cloud Run. Signal the sub process appropriately
	// after making relevant changes to the config and/or health signals.
	// TODO(b/307317433): Consider having a timeout to shutdown the subprocess
	// non-gracefully.
	sig := <-signalChan
	log.Printf("entrypoint: %s signal caught", sig)

	collectorProcess.Signal(sig)
	processState, err := collectorProcess.Wait()
	if err != nil {
		log.Fatalf(processState.String(), err)
	}
	log.Print("entrypoint: sidecar exited")
}
