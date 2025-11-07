// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestSchemaRegistryNormalizationResourceCRUDHappyFlowWithProviderSwap(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: cloudProviderConfig + `
resource "foxcon_schema_registry_normalization" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  normalization_enabled = ` + normalization_enabled_true + `
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_schema_registry_normalization.test", "normalization_enabled", normalization_enabled_true),
					resource.TestCheckResourceAttr("foxcon_schema_registry_normalization.test", "rest_endpoint", rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_schema_registry_normalization.test", "last_updated"),
				),
			},
			{
				Config: schemaProviderConfig + `
resource "foxcon_schema_registry_normalization" "test" {
  normalization_enabled = ` + normalization_enabled_false + `
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"foxcon_schema_registry_normalization.test",
						tfjsonpath.New("rest_endpoint"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"foxcon_schema_registry_normalization.test",
						tfjsonpath.New("credentials"),
						knownvalue.Null(),
					),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_schema_registry_normalization.test", "normalization_enabled", normalization_enabled_false),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_schema_registry_normalization.test", "last_updated"),
				),
			},
			{
				Config: schemaProviderConfig + `
			resource "foxcon_schema_registry_normalization" "test" {
			  rest_endpoint = "` + rest_endpoint + `"
			  credentials {
			    key = "` + api_key + `"
			    secret = "` + api_secret + `"
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
					resource.TestCheckResourceAttr("foxcon_schema_registry_normalization.test", "rest_endpoint", rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_schema_registry_normalization.test", "last_updated"),
				),
			},
			{
				Config: cloudProviderConfig + `
resource "foxcon_schema_registry_normalization" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  normalization_enabled = ` + normalization_enabled_true + `
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_schema_registry_normalization.test", "normalization_enabled", normalization_enabled_true),
					resource.TestCheckResourceAttr("foxcon_schema_registry_normalization.test", "rest_endpoint", rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_schema_registry_normalization.test", "last_updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestSchemaRegistryNormalizationResourceImportHappyFlow(t *testing.T) {

	t.Setenv("IMPORT_SCHEMA_REGISTRY_REST_ENDPOINT", rest_endpoint)
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_KEY", api_key)
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_SECRET", api_secret)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Define resource
			{
				Config: cloudProviderConfig + `
resource "foxcon_schema_registry_normalization" "test_import" {
  rest_endpoint = "` + rest_endpoint + `"
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

func TestSchemaRegistryNormalizationResourceNoRestEndpointParameter(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
resource "foxcon_schema_registry_normalization" "test" {
  normalization_enabled = ` + normalization_enabled_true + `
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				ExpectError: regexp.MustCompile(`Missing Required Attribute "rest_endpoint"`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestSchemaRegistryNormalizationResourceNoCredentialsConfigBlock(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: cloudProviderConfig + `
resource "foxcon_schema_registry_normalization" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  normalization_enabled = ` + normalization_enabled_true + `
}
`,
				ExpectError: regexp.MustCompile(`Error: Missing Required Attribute "credentials"`),
			},
		},
	})
}

func TestSchemaRegistryNormalizationResourceImportNoRestEndpointSet(t *testing.T) {

	t.Setenv("IMPORT_SCHEMA_REGISTRY_REST_ENDPOINT", "")
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_KEY", api_key)
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_SECRET", api_secret)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Define resource
			{
				Config: cloudProviderConfig + `
resource "foxcon_schema_registry_normalization" "test_import" {
  rest_endpoint = "` + rest_endpoint + `"
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
				ResourceName: "foxcon_schema_registry_normalization.test_import",
				ImportState:  true,
				//ImportStateVerify: true,
				ImportStateKind: resource.ImportBlockWithID,
				ImportStateId:   "ok",
				ExpectError:     regexp.MustCompile(`'IMPORT_SCHEMA_REGISTRY_REST_ENDPOINT' environment variable is not configured`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestSchemaRegistryNormalizationResourceImportNoApiSecretSet(t *testing.T) {

	t.Setenv("IMPORT_SCHEMA_REGISTRY_REST_ENDPOINT", rest_endpoint)
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_KEY", api_key)
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_SECRET", "")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Define resource
			{
				Config: cloudProviderConfig + `
resource "foxcon_schema_registry_normalization" "test_import" {
  rest_endpoint = "` + rest_endpoint + `"
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
				ResourceName: "foxcon_schema_registry_normalization.test_import",
				ImportState:  true,
				//ImportStateVerify: true,
				ImportStateKind: resource.ImportBlockWithID,
				ImportStateId:   "ok",
				ExpectError:     regexp.MustCompile(`'IMPORT_SCHEMA_REGISTRY_API_SECRET' environment variable is not configured`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestSchemaRegistryNormalizationResourceImportNoApiKeySet(t *testing.T) {

	t.Setenv("IMPORT_SCHEMA_REGISTRY_REST_ENDPOINT", rest_endpoint)
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_KEY", "")
	t.Setenv("IMPORT_SCHEMA_REGISTRY_API_SECRET", api_secret)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Define resource
			{
				Config: cloudProviderConfig + `
resource "foxcon_schema_registry_normalization" "test_import" {
  rest_endpoint = "` + rest_endpoint + `"
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
				ResourceName: "foxcon_schema_registry_normalization.test_import",
				ImportState:  true,
				//ImportStateVerify: true,
				ImportStateKind: resource.ImportBlockWithID,
				ImportStateId:   "ok",
				ExpectError:     regexp.MustCompile(`'IMPORT_SCHEMA_REGISTRY_API_KEY' environment variable is not configured`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestSchemaRegistryNormalizationProviderConfiguredWithEnvVars(t *testing.T) {

	t.Setenv("SCHEMA_REGISTRY_REST_ENDPOINT", rest_endpoint)
	t.Setenv("SCHEMA_REGISTRY_API_KEY", api_key)
	t.Setenv("SCHEMA_REGISTRY_API_SECRET", api_secret)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Define resource
			{
				Config: emptyProviderConfig + `
resource "foxcon_schema_registry_normalization" "test" {
  normalization_enabled = ` + normalization_enabled_true + `
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_schema_registry_normalization.test", "normalization_enabled", normalization_enabled_true),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_schema_registry_normalization.test", "last_updated"),
				),
			},
		},
	})
}

func TestSchemaRegistryNormalizationResourceWrongCredentials(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: schemaProviderWrongConfig + `
resource "foxcon_schema_registry_normalization" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  normalization_enabled = ` + normalization_enabled_true + `
  credentials {
    key = "dummy_value"
    secret = "dummy_value"
  }
}
`,
				ExpectError: regexp.MustCompile(`Response code 401`),
			},
		},
	})
}

func TestSchemaRegistryNormalizationProviderMigration111To121Setup(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: cloudProviderConfig + `
resource "foxcon_schema_registry_normalization" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  normalization_enabled = ` + normalization_enabled_true + `
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_schema_registry_normalization.test", "normalization_enabled", normalization_enabled_true),
					resource.TestCheckResourceAttr("foxcon_schema_registry_normalization.test", "rest_endpoint", rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_schema_registry_normalization.test", "last_updated"),
				),
			},
			{
				Config: schemaProviderConfig + `
resource "foxcon_schema_registry_normalization" "test" {
  normalization_enabled = ` + normalization_enabled_true + `
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_schema_registry_normalization.test", "normalization_enabled", normalization_enabled_true),
					resource.TestCheckNoResourceAttr("foxcon_schema_registry_normalization.test", "rest_endpoint"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_schema_registry_normalization.test", "last_updated"),
				),
			},
		},
	})
}
