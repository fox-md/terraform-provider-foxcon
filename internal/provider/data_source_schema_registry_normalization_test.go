// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestSchemaRegistryNormalizationDataSourceReadTrueValue(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Define resource
			{
				PreConfig: func() {

					jsonPayload := []byte(`{"normalize": true}`)
					req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/config", rest_endpoint), bytes.NewBuffer(jsonPayload))
					req.Header.Set("Content-Type", "application/json")
					req.SetBasicAuth(api_key, api_secret)
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						panic("failed to send HTTP request:" + err.Error())
					}

					body, _ := io.ReadAll(resp.Body)
					strbody := string(body)
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						panic(fmt.Sprintf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK))
					}

					if strbody != "{\"normalize\":true}" {
						panic(fmt.Errorf("unexpected body: got '%s', want '%s'", strbody, "{\"normalize\":true}"))
					}
				},
				Config: cloudProviderConfig + `
data "foxcon_schema_registry_normalization" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.foxcon_schema_registry_normalization.test", "normalization_enabled", "true"),
					resource.TestCheckResourceAttr("data.foxcon_schema_registry_normalization.test", "rest_endpoint", rest_endpoint),
				),
			},
		},
	})
}

func TestSchemaRegistryNormalizationDataSourceReadFalseValue(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Define resource
			{
				PreConfig: func() {
					jsonPayload := []byte(`{"normalize": false}`)
					req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/config", rest_endpoint), bytes.NewBuffer(jsonPayload))
					req.Header.Set("Content-Type", "application/json")
					req.SetBasicAuth(api_key, api_secret)
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						panic("failed to send HTTP request:" + err.Error())
					}

					body, _ := io.ReadAll(resp.Body)
					strbody := string(body)
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						panic(fmt.Sprintf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK))
					}

					if strbody != "{\"normalize\":false}" {
						panic(fmt.Errorf("unexpected body: got '%s', want '%s'", strbody, "{\"normalize\":false}"))
					}
				},
				Config: cloudProviderConfig + `
data "foxcon_schema_registry_normalization" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.foxcon_schema_registry_normalization.test", "normalization_enabled", "false"),
					resource.TestCheckResourceAttr("data.foxcon_schema_registry_normalization.test", "rest_endpoint", rest_endpoint),
				),
			},
		},
	})
}

func TestSchemaRegistryNormalizationDataSourceReadNullValue(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Define resource
			{
				PreConfig: func() {
					jsonPayload := []byte(`{"normalize": null}`)
					req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/config", rest_endpoint), bytes.NewBuffer(jsonPayload))
					req.Header.Set("Content-Type", "application/json")
					req.SetBasicAuth(api_key, api_secret)
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						panic("failed to send HTTP request:" + err.Error())
					}

					body, _ := io.ReadAll(resp.Body)
					strbody := string(body)
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						panic(fmt.Sprintf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK))
					}

					if strbody != "{}" {
						panic(fmt.Errorf("unexpected body: got '%s', want '%s'", strbody, "{}"))
					}
				},
				Config: cloudProviderConfig + `
data "foxcon_schema_registry_normalization" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.foxcon_schema_registry_normalization.test",
						tfjsonpath.New("normalization_enabled"),
						knownvalue.Null(),
					),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.foxcon_schema_registry_normalization.test", "rest_endpoint", rest_endpoint),
				),
			},
		},
	})
}
