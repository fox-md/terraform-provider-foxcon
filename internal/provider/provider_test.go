// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	cloudProviderConfig = `
provider "foxcon" {
  cloud_api_key = "test"
  cloud_api_secret = "test"
}
`
	emptyProviderConfig = `
provider "foxcon" {
}
`
	schemaProviderWrongConfig = `
provider "foxcon" {
  schema_registry_rest_endpoint = "http://1.1.1.1"
  schema_registry_api_key = "dummy_value"
  schema_registry_api_secret = "dummy_value"
}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"foxcon": providerserver.NewProtocol6WithError(New("test")()),
	}
)

var rest_endpoint string = "http://localhost:8081"
var api_key string = "admin"
var api_secret string = "admin-secret"
var subject_name string
var subject_name_imported string = "test_import"
var normalization_enabled_true string = "true"
var normalization_enabled_false string = "false"

var schemaProviderConfig = `
provider "foxcon" {
  schema_registry_rest_endpoint = "` + rest_endpoint + `"
  schema_registry_api_key = "` + api_key + `"
  schema_registry_api_secret = "` + api_secret + `"
}
`

type Payload struct {
	Schema     string `json:"schema"`
	SchemaType string `json:"schemaType"`
}

func callSchemaRegistry(method string, endpoint string, body io.Reader) (string, int, error) {
	req, _ := http.NewRequest(method, endpoint, body)
	req.Header.Set("Content-Type", "application/vnd.schemaregistry.v1+json")
	req.SetBasicAuth(api_key, api_secret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", resp.StatusCode, fmt.Errorf("failed to send HTTP request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", resp.StatusCode, fmt.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
	}

	respBody, _ := io.ReadAll(resp.Body)
	strbody := string(respBody)
	return strbody, resp.StatusCode, nil
}

func validateSubjectVersions(subject string, expectedResponse string) error {
	strbody, respCode, err := callSchemaRegistry("GET", fmt.Sprintf("%s/subjects/%s/versions?deleted=true", rest_endpoint, subject), nil)

	if respCode == 404 && "[]" == expectedResponse {
		return nil
	}

	if err != nil {
		return err
	}

	if strbody != expectedResponse {
		return fmt.Errorf("unexpected body: got '%s', want '%s'", strbody, expectedResponse)
	}
	return nil
}

func removeSubjectVersions(subject string, schemasToRemove []int) error {

	for _, i := range schemasToRemove {
		_, _, err := callSchemaRegistry("DELETE", fmt.Sprintf("%s/subjects/%s/versions/%d", rest_endpoint, subject, i), nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func addSubjectVersions(subject string, schemasToAdd []int) error {
	schemasLocation := "tests/schemas"

	gitRootCmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := gitRootCmd.Output()
	if err != nil {
		fmt.Printf("Failed to get git root: %v\n", err)
	}
	gitRoot := strings.TrimSpace(string(output))

	jsonPayload := []byte(`{"compatibility": "NONE"}`)
	_, _, err = callSchemaRegistry("PUT", fmt.Sprintf("%s/config/%s", rest_endpoint, subject), bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	for _, i := range schemasToAdd {
		data, err := os.ReadFile(fmt.Sprintf("%s/%s/v%d.json", gitRoot, schemasLocation, i))
		if err != nil {
			return fmt.Errorf("failed to open file: %s", err)
		}

		payload := Payload{
			Schema:     string(data),
			SchemaType: "JSON",
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		_, _, err = callSchemaRegistry("POST", fmt.Sprintf("%s/subjects/%s/versions", rest_endpoint, subject), bytes.NewBuffer(jsonPayload))
		if err != nil {
			return err
		}
	}
	return nil
}
