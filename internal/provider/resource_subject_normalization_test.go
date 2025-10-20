// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestSubjectNormalizationResourceTrueToFalse(t *testing.T) {

	subject_name = "true-false"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "foxcon_subject_normalization" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  normalization_enabled = ` + normalization_enabled_true + `
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "normalization_enabled", normalization_enabled_true),
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "rest_endpoint", rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_normalization.test", "last_updated"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "foxcon_subject_normalization" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  normalization_enabled = ` + normalization_enabled_false + `
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "normalization_enabled", normalization_enabled_false),
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "rest_endpoint", rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_normalization.test", "last_updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestSubjectNormalizationResourceFalseToTrueToNull(t *testing.T) {

	subject_name = "false-true-null"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "foxcon_subject_normalization" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  normalization_enabled = ` + normalization_enabled_false + `
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "normalization_enabled", normalization_enabled_false),
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "rest_endpoint", rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_normalization.test", "last_updated"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "foxcon_subject_normalization" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  normalization_enabled = ` + normalization_enabled_true + `
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "normalization_enabled", normalization_enabled_true),
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "rest_endpoint", rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_normalization.test", "last_updated"),
				),
			},
			{
				Config: providerConfig + `
resource "foxcon_subject_normalization" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"foxcon_subject_normalization.test",
						tfjsonpath.New("normalization_enabled"),
						knownvalue.Null(),
					),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "rest_endpoint", rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_normalization.test", "last_updated"),
				),
			},
			{
				PreConfig: func() {
					time.Sleep(1 * time.Second)
				},
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						req, _ := http.NewRequest("GET", fmt.Sprintf("%s/config/%s", rest_endpoint, subject_name), nil)
						req.SetBasicAuth(api_key, api_secret)
						resp, err := http.DefaultClient.Do(req)
						if err != nil {
							return fmt.Errorf("failed to send HTTP request: %s", err)
						}
						defer resp.Body.Close()

						if resp.StatusCode != http.StatusNotFound {
							return fmt.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusNotFound)
						}
						return nil
					},
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestSubjectNormalizationResourceNullToFalse(t *testing.T) {

	subject_name = "null-false"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "foxcon_subject_normalization" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"foxcon_subject_normalization.test",
						tfjsonpath.New("normalization_enabled"),
						knownvalue.Null(),
					),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "rest_endpoint", rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_normalization.test", "last_updated"),
				),
			},
			{
				Config: providerConfig + `
resource "foxcon_subject_normalization" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  normalization_enabled = ` + normalization_enabled_false + `
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "normalization_enabled", normalization_enabled_false),
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "rest_endpoint", rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_normalization.test", "last_updated"),
				),
			},
		},
	})
}

func TestSubjectNormalizationResourceSubjectConfigDeletionOnDelete(t *testing.T) {

	subject_name = "deletion"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "foxcon_subject_normalization" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  normalization_enabled = ` + normalization_enabled_false + `
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "normalization_enabled", normalization_enabled_false),
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "rest_endpoint", rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_normalization.test", "last_updated"),
				),
			},
			{
				Config: providerConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						req, _ := http.NewRequest("GET", fmt.Sprintf("%s/config/%s", rest_endpoint, subject_name), nil)
						req.SetBasicAuth(api_key, api_secret)
						resp, err := http.DefaultClient.Do(req)
						if err != nil {
							return fmt.Errorf("failed to send HTTP request: %s", err)
						}

						defer resp.Body.Close()

						if resp.StatusCode != http.StatusNotFound {
							return fmt.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusNotFound)
						}
						return nil
					},
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestSubjectNormalizationResourceKeepSubjectConfigOnDelete(t *testing.T) {

	subject_name = "no-deletion"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "foxcon_subject_normalization" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  normalization_enabled = ` + normalization_enabled_false + `
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "normalization_enabled", normalization_enabled_false),
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "rest_endpoint", rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_normalization.test", "last_updated"),
				),
			},
			{
				PreConfig: func() {
					jsonPayload := []byte(`{"compatibilityGroup": "test"}`)
					req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/config/%s", rest_endpoint, subject_name), bytes.NewBuffer(jsonPayload))
					req.Header.Set("Content-Type", "application/json")
					req.SetBasicAuth(api_key, api_secret)
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						panic("failed to send HTTP request:" + err.Error())
					}

					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						panic(fmt.Sprintf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK))
					}

				},
				Config: providerConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						req, _ := http.NewRequest("GET", fmt.Sprintf("%s/config/%s", rest_endpoint, subject_name), nil)
						req.SetBasicAuth(api_key, api_secret)
						resp, err := http.DefaultClient.Do(req)
						if err != nil {
							return fmt.Errorf("failed to send HTTP request: %s", err)
						}

						defer resp.Body.Close()

						if resp.StatusCode != http.StatusOK {
							return fmt.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
						}
						return nil
					},
				),
			},
		},
	})
}

func TestSubjectNormalizationResourceKeepSubjectConfigOnDeleteFullCompatLevel(t *testing.T) {

	subject_name = "full-compat"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "foxcon_subject_normalization" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  normalization_enabled = ` + normalization_enabled_false + `
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "normalization_enabled", normalization_enabled_false),
					resource.TestCheckResourceAttr("foxcon_subject_normalization.test", "rest_endpoint", rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_normalization.test", "last_updated"),
				),
			},
			{
				PreConfig: func() {
					jsonPayload := []byte(`{"compatibility": "FULL"}`)
					req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/config/%s", rest_endpoint, subject_name), bytes.NewBuffer(jsonPayload))
					req.Header.Set("Content-Type", "application/json")
					req.SetBasicAuth(api_key, api_secret)
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						panic("failed to send HTTP request:" + err.Error())
					}

					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						panic(fmt.Sprintf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK))
					}

				},
				Config: providerConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						req, _ := http.NewRequest("GET", fmt.Sprintf("%s/config/%s", rest_endpoint, subject_name), nil)
						req.SetBasicAuth(api_key, api_secret)
						resp, err := http.DefaultClient.Do(req)
						if err != nil {
							return fmt.Errorf("failed to send HTTP request: %s", err)
						}
						defer resp.Body.Close()

						if resp.StatusCode != http.StatusOK {
							return fmt.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
						}
						return nil
					},
				),
			},
		},
	})
}

func TestSubjectNormalizationResourceImportHappyFlow(t *testing.T) {

	t.Setenv("IMPORT_SCHEMA_REGISTRY_REST_ENDPOINT", rest_endpoint)
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_KEY", api_key)
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_SECRET", api_secret)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Define resource
			{
				Config: providerConfig + `
resource "foxcon_subject_normalization" "test_import" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name_imported + `"
  normalization_enabled = ` + normalization_enabled_true + `
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
			},
			// ImportState testing
			{
				ResourceName: "foxcon_subject_normalization.test_import",
				ImportState:  true,
				//ImportStateVerify: true,
				ImportStateKind: resource.ImportBlockWithID,
				ImportStateId:   subject_name_imported,
				//ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestSubjectNormalizationResourceNoSubjectNameParameter(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "foxcon_subject_normalization" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  normalization_enabled = ` + normalization_enabled_true + `
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				ExpectError: regexp.MustCompile(`The argument "subject_name" is required`),
			},
		},
	})
}

func TestSubjectNormalizationResourceNoRestEndpointParameter(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "foxcon_subject_normalization" "test" {
  subject_name = "` + subject_name + `"
  normalization_enabled = ` + normalization_enabled_true + `
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				ExpectError: regexp.MustCompile(`The argument "rest_endpoint" is required`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestSubjectNormalizationResourceNoCredentialsConfigBlock(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "foxcon_subject_normalization" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  normalization_enabled = ` + normalization_enabled_true + `
}
`,
				ExpectError: regexp.MustCompile(`Missing Configuration for Required Attribute`),
			},
		},
	})
}

func TestSubjectNormalizationResourceImportNoRestEndpointSet(t *testing.T) {

	t.Setenv("IMPORT_SCHEMA_REGISTRY_REST_ENDPOINT", "")
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_KEY", api_key)
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_SECRET", api_secret)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Define resource
			{
				Config: providerConfig + `
resource "foxcon_subject_normalization" "test_import" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name_imported + `"
  normalization_enabled = ` + normalization_enabled_true + `
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
			},
			// ImportState testing
			{
				ResourceName: "foxcon_subject_normalization.test_import",
				ImportState:  true,
				//ImportStateVerify: true,
				ImportStateKind: resource.ImportBlockWithID,
				ImportStateId:   subject_name_imported,
				//ImportStateVerifyIdentifierAttribute: subject_name_imported,
				//ImportStateVerifyIgnore: []string{"last_updated"},
				ExpectError: regexp.MustCompile(`'IMPORT_SCHEMA_REGISTRY_REST_ENDPOINT' environment variable is not configured`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestSubjectNormalizationResourceImportNoApiSecretSet(t *testing.T) {

	t.Setenv("IMPORT_SCHEMA_REGISTRY_REST_ENDPOINT", rest_endpoint)
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_KEY", api_key)
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_SECRET", "")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Define resource
			{
				Config: providerConfig + `
resource "foxcon_subject_normalization" "test_import" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name_imported + `"
  normalization_enabled = ` + normalization_enabled_true + `
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
			},
			// ImportState testing
			{
				ResourceName: "foxcon_subject_normalization.test_import",
				ImportState:  true,
				//ImportStateVerify: true,
				ImportStateKind: resource.ImportBlockWithID,
				ImportStateId:   subject_name_imported,
				//ImportStateVerifyIdentifierAttribute: subject_name_imported,
				//ImportStateVerifyIgnore: []string{"last_updated"},
				ExpectError: regexp.MustCompile(`'IMPORT_SCHEMA_REGISTRY_API_SECRET' environment variable is not configured`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestSubjectNormalizationResourceImportNoApiKeySet(t *testing.T) {

	t.Setenv("IMPORT_SCHEMA_REGISTRY_REST_ENDPOINT", rest_endpoint)
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_KEY", "")
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_SECRET", api_secret)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Define resource
			{
				Config: providerConfig + `
resource "foxcon_subject_normalization" "test_import" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name_imported + `"
  normalization_enabled = ` + normalization_enabled_true + `
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
			},
			// ImportState testing
			{
				ResourceName: "foxcon_subject_normalization.test_import",
				ImportState:  true,
				//ImportStateVerify: true,
				ImportStateKind: resource.ImportBlockWithID,
				ImportStateId:   subject_name_imported,
				//ImportStateVerifyIdentifierAttribute: subject_name_imported,
				//ImportStateVerifyIgnore: []string{"last_updated"},
				ExpectError: regexp.MustCompile(`'IMPORT_SCHEMA_REGISTRY_API_KEY' environment variable is not configured`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
