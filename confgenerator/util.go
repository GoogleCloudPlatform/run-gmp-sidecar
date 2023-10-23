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

// Package confgenerator represents the Ops Agent configuration and provides functions to generate subagents configuration from unified agent.
package confgenerator

import (
	"fmt"
	"net"
	"os"
	"strings"
	"text/template"

	"github.com/prometheus/prometheus/model/relabel"
)

// CloudRunEnvironment captures some environment metadata that is captures using environment variables.
// See https://cloud.google.com/run/docs/container-contract#services-env-vars for more information.
// Note that PORT is not made available to sidecar containers, and so is omitted from this struct.
type CloudRunEnvironment struct {
	Service       string
	Revision      string
	Configuration string
}

func fetchMetadata() *CloudRunEnvironment {
	return &CloudRunEnvironment{
		Service:       os.Getenv("K_SERVICE"),
		Revision:      os.Getenv("K_REVISION"),
		Configuration: os.Getenv("K_CONFIGURATION"),
	}
}

var versionLabelTemplate = template.Must(template.New("versionlabel").Parse(`{{.Prefix}}@{{.AgentVersion}}`))
var userAgentTemplate = template.Must(template.New("useragent").Parse(`{{.Prefix}}/{{.AgentVersion}}; ShortName={{.ShortName}};ShortVersion={{.ShortVersion}}`))

var Version = "latest"

func expandTemplate(t *template.Template, prefix string, extraParams map[string]string) (string, error) {
	params := map[string]string{
		"Prefix":       prefix,
		"AgentVersion": Version,
	}
	for k, v := range extraParams {
		params[k] = v
	}
	var b strings.Builder
	if err := t.Execute(&b, params); err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return b.String(), nil
}

func VersionLabel(prefix string) (string, error) {
	return expandTemplate(versionLabelTemplate, prefix, nil)
}

func UserAgent(prefix, shortName, shortVersion string) (string, error) {
	extraParams := map[string]string{
		"ShortName":    shortName,
		"ShortVersion": shortVersion,
	}
	return expandTemplate(userAgentTemplate, prefix, extraParams)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

var protectedLabels = []string{
	"project_id",
	"location",
	"cluster",
	"namespace",
	"job",
	"cloud_run_instance",
	"__address__",
}

func isProtectedLabel(s string) bool {
	return contains(protectedLabels, s)
}

func matchesAnyProtectedLabel(re relabel.Regexp) bool {
	for _, pl := range protectedLabels {
		if re.MatchString(pl) {
			return true
		}
	}
	return false
}

func matchesAllProtectedLabels(re relabel.Regexp) bool {
	for _, pl := range protectedLabels {
		if !re.MatchString(pl) {
			return false
		}
	}
	return true
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
