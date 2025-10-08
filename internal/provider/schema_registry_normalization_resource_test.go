// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var _rest_endpoint string = "http://localhost:8081"
var _api_key string = "admin"
var _api_secret string = "admin-secret"
var _normalization_enabled_true string = "true"
var _normalization_enabled_false string = "false"

func TestSchemaNormalizationResourceCRUDHappyFlow(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "foxcon_schema_registry_normalization" "test" {
  rest_endpoint = "` + _rest_endpoint + `"
  normalization_enabled = ` + _normalization_enabled_true + `
  credentials {
    key = "` + _api_key + `"
    secret = "` + _api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_schema_registry_normalization.test", "normalization_enabled", _normalization_enabled_true),
					resource.TestCheckResourceAttr("foxcon_schema_registry_normalization.test", "rest_endpoint", _rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_schema_registry_normalization.test", "last_updated"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "foxcon_schema_registry_normalization" "test" {
  rest_endpoint = "` + _rest_endpoint + `"
  normalization_enabled = ` + _normalization_enabled_false + `
  credentials {
    key = "` + _api_key + `"
    secret = "` + _api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_schema_registry_normalization.test", "normalization_enabled", _normalization_enabled_false),
					resource.TestCheckResourceAttr("foxcon_schema_registry_normalization.test", "rest_endpoint", _rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_schema_registry_normalization.test", "last_updated"),
				),
			},
			{
				Config: providerConfig + `
			resource "foxcon_schema_registry_normalization" "test" {
			  rest_endpoint = "` + _rest_endpoint + `"
			  credentials {
			    key = "` + _api_key + `"
			    secret = "` + _api_secret + `"
			  }
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["foxcon_schema_registry_normalization.test"]

						if !ok {
							return fmt.Errorf("resource not found")
						}

						if _, exists := rs.Primary.Attributes["normalization_enabled"]; exists {
							return fmt.Errorf("expected 'normalization_enabled' to be missing, but it exists")
						}

						return nil

					},
					resource.TestCheckResourceAttr("foxcon_schema_registry_normalization.test", "rest_endpoint", _rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_schema_registry_normalization.test", "last_updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestSchemaNormalizationResourceImportHappyFlow(t *testing.T) {

	t.Setenv("IMPORT_SCHEMA_REGISTRY_REST_ENDPOINT", _rest_endpoint)
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_KEY", _api_key)
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_SECRET", _api_secret)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Define resource
			{
				Config: providerConfig + `
resource "foxcon_schema_registry_normalization" "test_import" {
  rest_endpoint = "` + _rest_endpoint + `"
  normalization_enabled = ` + _normalization_enabled_true + `
  credentials {
    key = "` + _api_key + `"
    secret = "` + _api_secret + `"
  }
}
`,
			},
			// ImportState testing
			{
				ResourceName: "foxcon_schema_registry_normalization.test_import",
				ImportState:  true,
				//ImportStateVerify: true,
				ImportStateKind: resource.ImportBlockWithID,
				ImportStateId:   "ok",
				//ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestSchemaNormalizationResourceNoRestEndpointParameter(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "foxcon_schema_registry_normalization" "test" {
  normalization_enabled = ` + _normalization_enabled_true + `
  credentials {
    key = "` + _api_key + `"
    secret = "` + _api_secret + `"
  }
}
`,
				ExpectError: regexp.MustCompile(`The argument "rest_endpoint" is required`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestSchemaNormalizationResourceNoCredentialsConfigBlock(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "foxcon_schema_registry_normalization" "test" {
  rest_endpoint = "` + _rest_endpoint + `"
  normalization_enabled = ` + _normalization_enabled_true + `
}
`,
				ExpectError: regexp.MustCompile(`Missing Configuration for Required Attribute`),
			},
		},
	})
}

func TestSchemaNormalizationResourceImportNoRestEndpointSet(t *testing.T) {

	t.Setenv("IMPORT_SCHEMA_REGISTRY_REST_ENDPOINT", "")
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_KEY", _api_key)
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_SECRET", _api_secret)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Define resource
			{
				Config: providerConfig + `
resource "foxcon_schema_registry_normalization" "test_import" {
  rest_endpoint = "` + _rest_endpoint + `"
  normalization_enabled = ` + _normalization_enabled_true + `
  credentials {
    key = "` + _api_key + `"
    secret = "` + _api_secret + `"
  }
}
`,
			},
			// ImportState testing
			{
				ResourceName: "foxcon_schema_registry_normalization.test_import",
				ImportState:  true,
				//ImportStateVerify: true,
				ImportStateKind: resource.ImportBlockWithID,
				ImportStateId:   "ok",
				//ImportStateVerifyIdentifierAttribute: subject_name_imported,
				//ImportStateVerifyIgnore: []string{"last_updated"},
				ExpectError: regexp.MustCompile(`'IMPORT_SCHEMA_REGISTRY_REST_ENDPOINT' environment variable is not configured`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestSchemaNormalizationResourceImportNoApiSecretSet(t *testing.T) {

	t.Setenv("IMPORT_SCHEMA_REGISTRY_REST_ENDPOINT", _rest_endpoint)
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_KEY", _api_key)
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_SECRET", "")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Define resource
			{
				Config: providerConfig + `
resource "foxcon_schema_registry_normalization" "test_import" {
  rest_endpoint = "` + _rest_endpoint + `"
  normalization_enabled = ` + _normalization_enabled_true + `
  credentials {
    key = "` + _api_key + `"
    secret = "` + _api_secret + `"
  }
}
`,
			},
			// ImportState testing
			{
				ResourceName: "foxcon_schema_registry_normalization.test_import",
				ImportState:  true,
				//ImportStateVerify: true,
				ImportStateKind: resource.ImportBlockWithID,
				ImportStateId:   "ok",
				//ImportStateVerifyIdentifierAttribute: subject_name_imported,
				//ImportStateVerifyIgnore: []string{"last_updated"},
				ExpectError: regexp.MustCompile(`'IMPORT_SCHEMA_REGISTRY_API_SECRET' environment variable is not configured`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestSchemaNormalizationResourceImportNoApiKeySet(t *testing.T) {

	t.Setenv("IMPORT_SCHEMA_REGISTRY_REST_ENDPOINT", _rest_endpoint)
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_KEY", "")
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_SECRET", _api_secret)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Define resource
			{
				Config: providerConfig + `
resource "foxcon_schema_registry_normalization" "test_import" {
  rest_endpoint = "` + _rest_endpoint + `"
  normalization_enabled = ` + _normalization_enabled_true + `
  credentials {
    key = "` + _api_key + `"
    secret = "` + _api_secret + `"
  }
}
`,
			},
			// ImportState testing
			{
				ResourceName: "foxcon_schema_registry_normalization.test_import",
				ImportState:  true,
				//ImportStateVerify: true,
				ImportStateKind: resource.ImportBlockWithID,
				ImportStateId:   "ok",
				//ImportStateVerifyIdentifierAttribute: subject_name_imported,
				//ImportStateVerifyIgnore: []string{"last_updated"},
				ExpectError: regexp.MustCompile(`'IMPORT_SCHEMA_REGISTRY_API_KEY' environment variable is not configured`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
