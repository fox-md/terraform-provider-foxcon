// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var rest_endpoint string = "http://localhost:8081"
var api_key string = "admin"
var api_secret string = "admin-secret"
var subject_name string = "test"
var subject_name_imported string = "test_import"
var normalization_enabled_true string = "true"
var normalization_enabled_false string = "false"

func TestSubjectNormalizationResourceCRUDHappyFlow(t *testing.T) {

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
			// ImportState testing
			// {
			// 	ResourceName:      "foxcon_subject_normalization.test",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	// The last_updated attribute does not exist in the HashiCups
			// 	// API, therefore there is no value for it during import.
			// 	ImportStateVerifyIgnore: []string{"last_updated"},
			// },
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
