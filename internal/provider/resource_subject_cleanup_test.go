// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestSubjectCleanupLatestHappyFlow(t *testing.T) {

	subject_name = "keep-latest"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						schemasToAdd := []int{1, 2, 3, 4, 5}
						err := addSubjectVersions(subject_name, schemasToAdd)
						if err != nil {
							return err
						}

						schemasToRemove := []int{1, 2}
						err = removeSubjectVersions(subject_name, schemasToRemove)
						if err != nil {
							return err
						}

						return nil
					},
				),
			},
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "latest" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  cleanup_method = "KEEP_LATEST_ONLY"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "rest_endpoint", rest_endpoint),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "cleanup_method", "KEEP_LATEST_ONLY"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "last_deleted.#", "4"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "last_deleted.0", "1"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "last_deleted.1", "2"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "last_deleted.2", "3"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "last_deleted.3", "4"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "latest_schema_version", "5"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "cleanup_needed", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_cleanup.latest", "last_updated"),
				),
			},
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						expected := "[5]"
						err := validateSubjectVersions(subject_name, expected)
						if err != nil {
							return err
						}
						return nil
					},
				),
			},
		},
	})
}

func TestSubjectCleanupActiveHappyFlow(t *testing.T) {

	subject_name = "keep-active"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						schemasToAdd := []int{1, 2, 3, 4, 5}
						err := addSubjectVersions(subject_name, schemasToAdd)
						if err != nil {
							return err
						}

						schemasToRemove := []int{1, 2, 3}
						err = removeSubjectVersions(subject_name, schemasToRemove)
						if err != nil {
							return err
						}

						return nil
					},
				),
			},
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "active" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  cleanup_method = "KEEP_ACTIVE_ONLY"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.active", "rest_endpoint", rest_endpoint),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.active", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.active", "cleanup_method", "KEEP_ACTIVE_ONLY"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.active", "last_deleted.#", "3"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.active", "last_deleted.0", "1"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.active", "last_deleted.1", "2"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.active", "last_deleted.2", "3"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.active", "latest_schema_version", "5"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.active", "cleanup_needed", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_cleanup.active", "last_updated"),
				),
			},
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						expected := "[4,5]"
						err := validateSubjectVersions(subject_name, expected)
						if err != nil {
							return err
						}
						return nil
					},
				),
			},
		},
	})
}

func TestSubjectCleanupKeepLatestIfOneVersion(t *testing.T) {

	subject_name = "one"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "latest" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  cleanup_method = "KEEP_LATEST_ONLY"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "rest_endpoint", rest_endpoint),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "cleanup_method", "KEEP_LATEST_ONLY"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "last_deleted.#", "0"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "latest_schema_version", "1"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "cleanup_needed", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_cleanup.latest", "last_updated"),
				),
			},
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						expected := "[1]"
						err := validateSubjectVersions(subject_name, expected)
						if err != nil {
							return err
						}
						return nil
					},
				),
			},
		},
	})
}

func TestSubjectCleanupKeepActiveIfOneVersion(t *testing.T) {

	subject_name = "one"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "latest" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  cleanup_method = "KEEP_ACTIVE_ONLY"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "rest_endpoint", rest_endpoint),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "cleanup_method", "KEEP_ACTIVE_ONLY"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "last_deleted.#", "0"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "latest_schema_version", "1"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "cleanup_needed", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_cleanup.latest", "last_updated"),
				),
			},
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						expected := "[1]"
						err := validateSubjectVersions(subject_name, expected)
						if err != nil {
							return err
						}
						return nil
					},
				),
			},
		},
	})
}

func TestSubjectCleanupKeepActiveIfOneVersionWrongProvider(t *testing.T) {

	subject_name = "one"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: schemaProviderWrongConfig + `
resource "foxcon_subject_cleanup" "latest" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  cleanup_method = "KEEP_ACTIVE_ONLY"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "rest_endpoint", rest_endpoint),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "cleanup_method", "KEEP_ACTIVE_ONLY"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "last_deleted.#", "0"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "latest_schema_version", "1"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "cleanup_needed", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_cleanup.latest", "last_updated"),
				),
			},
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						expected := "[1]"
						err := validateSubjectVersions(subject_name, expected)
						if err != nil {
							return err
						}
						return nil
					},
				),
			},
		},
	})
}

func TestSubjectCleanupKeepActiveIfOneVersionUsingProvider(t *testing.T) {

	subject_name = "one"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: schemaProviderConfig + `
resource "foxcon_subject_cleanup" "latest" {
  subject_name = "` + subject_name + `"
  cleanup_method = "KEEP_ACTIVE_ONLY"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "cleanup_method", "KEEP_ACTIVE_ONLY"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "last_deleted.#", "0"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "latest_schema_version", "1"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "cleanup_needed", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_cleanup.latest", "last_updated"),
				),
			},
			{
				Config: schemaProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						expected := "[1]"
						err := validateSubjectVersions(subject_name, expected)
						if err != nil {
							return err
						}
						return nil
					},
				),
			},
		},
	})
}

func TestSubjectCleanupKeepActiveIfOneVersionNoProviderNoConfigErrorHandling(t *testing.T) {

	subject_name = "test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: emptyProviderConfig + `
resource "foxcon_subject_cleanup" "latest" {
  subject_name = "` + subject_name + `"
  cleanup_method = "KEEP_ACTIVE_ONLY"
}
`,
				ExpectError: regexp.MustCompile(`Error creating http client`),
			},
		},
	})
}

func TestSubjectCleanupKeepActiveIfOneVersionEnvVarsProviderConfig(t *testing.T) {

	t.Setenv("SCHEMA_REGISTRY_REST_ENDPOINT", rest_endpoint)
	t.Setenv("SCHEMA_REGISTRY_API_KEY", api_key)
	t.Setenv("SCHEMA_REGISTRY_API_SECRET", api_secret)

	subject_name = "one"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: emptyProviderConfig + `
resource "foxcon_subject_cleanup" "latest" {
  subject_name = "` + subject_name + `"
  cleanup_method = "KEEP_ACTIVE_ONLY"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "cleanup_method", "KEEP_ACTIVE_ONLY"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "last_deleted.#", "0"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "latest_schema_version", "1"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "cleanup_needed", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_cleanup.latest", "last_updated"),
				),
			},
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						expected := "[1]"
						err := validateSubjectVersions(subject_name, expected)
						if err != nil {
							return err
						}
						return nil
					},
				),
			},
		},
	})
}

func TestSubjectCleanupAddSchemaVersionWhenActive(t *testing.T) {

	subject_name = "one-to-two-active"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						schemasToAdd := []int{1}
						err := addSubjectVersions(subject_name, schemasToAdd)
						if err != nil {
							return err
						}

						return nil
					},
				),
			},
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "latest" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  cleanup_method = "KEEP_ACTIVE_ONLY"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "rest_endpoint", rest_endpoint),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "cleanup_method", "KEEP_ACTIVE_ONLY"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "last_deleted.#", "0"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "latest_schema_version", "1"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "cleanup_needed", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_cleanup.latest", "last_updated"),
				),
			},
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						schemasToAdd := []int{2}
						err := addSubjectVersions(subject_name, schemasToAdd)
						if err != nil {
							return err
						}

						expected := "[1,2]"
						err = validateSubjectVersions(subject_name, expected)
						if err != nil {
							return err
						}

						return nil
					},
				),
			},
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "latest" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  cleanup_method = "KEEP_ACTIVE_ONLY"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "rest_endpoint", rest_endpoint),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "cleanup_method", "KEEP_ACTIVE_ONLY"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "last_deleted.#", "0"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "latest_schema_version", "2"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "cleanup_needed", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_cleanup.latest", "last_updated"),
				),
			},
		},
	})
}

func TestSubjectCleanupAddSchemaVersionWhenLatest(t *testing.T) {

	subject_name = "one-to-two-latest"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						schemasToAdd := []int{1}
						err := addSubjectVersions(subject_name, schemasToAdd)
						if err != nil {
							return err
						}

						return nil
					},
				),
			},
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "latest" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  cleanup_method = "KEEP_LATEST_ONLY"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "rest_endpoint", rest_endpoint),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "cleanup_method", "KEEP_LATEST_ONLY"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "last_deleted.#", "0"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "latest_schema_version", "1"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "cleanup_needed", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_cleanup.latest", "last_updated"),
				),
			},
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						schemasToAdd := []int{2}
						err := addSubjectVersions(subject_name, schemasToAdd)
						if err != nil {
							return err
						}

						expected := "[1,2]"
						err = validateSubjectVersions(subject_name, expected)
						if err != nil {
							return err
						}

						return nil
					},
				),
			},
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "latest" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  cleanup_method = "KEEP_LATEST_ONLY"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "rest_endpoint", rest_endpoint),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "cleanup_method", "KEEP_LATEST_ONLY"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "last_deleted.#", "1"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "last_deleted.0", "1"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "latest_schema_version", "2"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.latest", "cleanup_needed", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_cleanup.latest", "last_updated"),
				),
			},
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						expected := "[2]"
						err := validateSubjectVersions(subject_name, expected)
						if err != nil {
							return err
						}
						return nil
					},
				),
			},
		},
	})
}

func TestSubjectCleanupFromActiveToLatestHappyFlow(t *testing.T) {

	subject_name = "switch-cleanup-mode"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						schemasToAdd := []int{1, 2, 3, 4, 5}
						err := addSubjectVersions(subject_name, schemasToAdd)
						if err != nil {
							return err
						}

						schemasToRemove := []int{1, 2}
						err = removeSubjectVersions(subject_name, schemasToRemove)
						if err != nil {
							return err
						}

						return nil
					},
				),
			},
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  cleanup_method = "KEEP_ACTIVE_ONLY"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "rest_endpoint", rest_endpoint),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "cleanup_method", "KEEP_ACTIVE_ONLY"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "last_deleted.#", "2"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "last_deleted.0", "1"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "last_deleted.1", "2"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "latest_schema_version", "5"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "cleanup_needed", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_cleanup.test", "last_updated"),
				),
			},
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						expected := "[3,4,5]"
						err := validateSubjectVersions(subject_name, expected)
						if err != nil {
							return err
						}
						return nil
					},
				),
			},
			{
				Config: cloudProviderConfig + `
			resource "foxcon_subject_cleanup" "test" {
			  rest_endpoint = "` + rest_endpoint + `"
			  subject_name = "` + subject_name + `"
			  cleanup_method = "KEEP_LATEST_ONLY"
			  credentials {
			    key = "` + api_key + `"
			    secret = "` + api_secret + `"
			  }
			}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "rest_endpoint", rest_endpoint),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "cleanup_method", "KEEP_LATEST_ONLY"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "last_deleted.#", "2"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "last_deleted.0", "3"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "last_deleted.1", "4"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "latest_schema_version", "5"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "cleanup_needed", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_cleanup.test", "last_updated"),
				),
			},
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						expected := "[5]"
						err := validateSubjectVersions(subject_name, expected)
						if err != nil {
							return err
						}
						return nil
					},
				),
			},
		},
	})
}

func TestSubjectCleanupNoCleanupMethodErrorHandling(t *testing.T) {

	subject_name = "dummy"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestSubjectCleanupNoCredentialsErrorHandling(t *testing.T) {

	subject_name = "dummy"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  cleanup_method = "KEEP_LATEST_ONLY"
}
`,
				ExpectError: regexp.MustCompile(`Missing Required Attribute`),
			},
		},
	})
}

func TestSubjectCleanupNoRestEndpointErrorHandling(t *testing.T) {

	subject_name = "dummy"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "test" {
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
  subject_name = "` + subject_name + `"
  cleanup_method = "KEEP_LATEST_ONLY"
}
`,
				ExpectError: regexp.MustCompile(`Missing Required Attribute`),
			},
		},
	})
}

func TestSubjectCleanupNoSecretErrorHandling(t *testing.T) {

	subject_name = "dummy"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  credentials {
    key = "` + api_key + `"
  }
  subject_name = "` + subject_name + `"
  cleanup_method = "KEEP_LATEST_ONLY"
}
`,
				ExpectError: regexp.MustCompile(`Missing Required Attribute`),
			},
		},
	})
}

func TestSubjectCleanupNoKeyErrorHandling(t *testing.T) {

	subject_name = "dummy"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  credentials {
    secret = "` + api_secret + `"
  }
  subject_name = "` + subject_name + `"
  cleanup_method = "KEEP_LATEST_ONLY"
}
`,
				ExpectError: regexp.MustCompile(`Missing Required Attribute`),
			},
		},
	})
}

func TestSubjectCleanupKeepActiveProviderMigration111To121Setup(t *testing.T) {

	subject_name = "one"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  cleanup_method = "KEEP_ACTIVE_ONLY"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "cleanup_method", "KEEP_ACTIVE_ONLY"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "last_deleted.#", "0"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "latest_schema_version", "1"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "cleanup_needed", "false"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "rest_endpoint", rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_cleanup.test", "last_updated"),
				),
			},
			{
				Config: schemaProviderConfig + `
resource "foxcon_subject_cleanup" "test" {
  subject_name = "` + subject_name + `"
  cleanup_method = "KEEP_ACTIVE_ONLY"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "cleanup_method", "KEEP_ACTIVE_ONLY"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "last_deleted.#", "0"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "latest_schema_version", "1"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "cleanup_needed", "false"),
					resource.TestCheckNoResourceAttr("foxcon_subject_cleanup.test", "rest_endpoint"),
					resource.TestCheckNoResourceAttr("foxcon_subject_cleanup.test", "credentials"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_cleanup.test", "last_updated"),
				),
			},
		},
	})
}

func TestSubjectCleanupNonExistingMethodErrorHandling(t *testing.T) {

	subject_name = "test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "latest" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  cleanup_method = "METHOD_THAT_DOESNOT_EXIST"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				ExpectError: regexp.MustCompile(`Attribute cleanup_method value must be one of: \[\"KEEP_LATEST_ONLY\"
\"KEEP_ACTIVE_ONLY\" \"MAX_STORED_SCHEMAS\"\]`),
			},
		},
	})
}

func TestSubjectCleanupKeepNActiveNoNumberErrorHandling(t *testing.T) {

	subject_name = "keep-n-active"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "latest" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  cleanup_method = "MAX_STORED_SCHEMAS"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				ExpectError: regexp.MustCompile(`Number of schemas must be more than 0 when cleanup_method is set to
\'MAX_STORED_SCHEMAS\'`),
			},
		},
	})
}

func TestSubjectCleanupMaxStoredHappyPath(t *testing.T) {

	subject_name = "subj-keep-n"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						schemasToAdd := []int{1, 2, 3, 4, 5}
						err := addSubjectVersions(subject_name, schemasToAdd)
						if err != nil {
							return err
						}

						return nil
					},
				),
			},
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  cleanup_method = "MAX_STORED_SCHEMAS"
  number_of_schemas_to_keep = 3
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "cleanup_method", "MAX_STORED_SCHEMAS"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "last_deleted.#", "2"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "last_deleted.0", "1"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "last_deleted.1", "2"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "latest_schema_version", "5"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "cleanup_needed", "false"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "rest_endpoint", rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_cleanup.test", "last_updated"),
				),
			},
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						expected := "[3,4,5]"
						err := validateSubjectVersions(subject_name, expected)
						if err != nil {
							return err
						}
						return nil
					},
				),
			},
		},
	})
}

func TestSubjectCleanupMaxStoredWithDeletedHappyPath(t *testing.T) {

	subject_name = "keep-n"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						schemasToAdd := []int{1, 2, 3, 4, 5}
						err := addSubjectVersions(subject_name, schemasToAdd)
						if err != nil {
							return err
						}

						schemasToRemove := []int{1, 2, 3}
						err = removeSubjectVersions(subject_name, schemasToRemove)
						if err != nil {
							return err
						}

						return nil
					},
				),
			},
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  cleanup_method = "MAX_STORED_SCHEMAS"
  number_of_schemas_to_keep = 4
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "cleanup_method", "MAX_STORED_SCHEMAS"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "last_deleted.#", "1"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "last_deleted.0", "1"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "latest_schema_version", "5"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "cleanup_needed", "false"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "rest_endpoint", rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_cleanup.test", "last_updated"),
				),
			},
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						expected := "[2,3,4,5]"
						err := validateSubjectVersions(subject_name, expected)
						if err != nil {
							return err
						}
						return nil
					},
				),
			},
		},
	})
}

func TestSubjectCleanupActiveOnlyToNOnly(t *testing.T) {

	subject_name = "subj-cleanup-new-method"

	resource.Test(t, resource.TestCase{

		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Config:                   cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						schemasToAdd := []int{1, 2, 3, 4, 5}
						err := addSubjectVersions(subject_name, schemasToAdd)
						if err != nil {
							return err
						}

						schemasToRemove := []int{1, 2}
						err = removeSubjectVersions(subject_name, schemasToRemove)
						if err != nil {
							return err
						}

						return nil
					},
				),
			},
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"foxcon": {
						VersionConstraint: "1.3.2",
						Source:            "registry.terraform.io/fox-md/foxcon",
					},
				},
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  cleanup_method = "KEEP_ACTIVE_ONLY"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "cleanup_method", "KEEP_ACTIVE_ONLY"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "last_deleted.#", "2"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "last_deleted.0", "1"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "last_deleted.1", "2"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "latest_schema_version", "5"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "cleanup_needed", "false"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "rest_endpoint", rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_cleanup.test", "last_updated"),
				),
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  cleanup_method = "MAX_STORED_SCHEMAS"
  number_of_schemas_to_keep = 2
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "cleanup_method", "MAX_STORED_SCHEMAS"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "last_deleted.#", "1"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "last_deleted.0", "3"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "latest_schema_version", "5"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "cleanup_needed", "false"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "rest_endpoint", rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_cleanup.test", "last_updated"),
				),
			},
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						expected := "[4,5]"
						err := validateSubjectVersions(subject_name, expected)
						if err != nil {
							return err
						}
						return nil
					},
				),
			},
		},
	})
}

func TestSubjectCleanupEmptySubjectErrorHandling(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = ""
  cleanup_method = "MAX_STORED_SCHEMAS"
  number_of_schemas_to_keep = 2
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				ExpectError: regexp.MustCompile(`Attribute subject_name string length must be at least 1`),
			},
		},
	})
}

func TestSubjectCleanupInvalidRestEndpointErrorHandling(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "test" {
  rest_endpoint = "localhost"
  subject_name = "test"
  cleanup_method = "MAX_STORED_SCHEMAS"
  number_of_schemas_to_keep = 2
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				ExpectError: regexp.MustCompile(`The value must start with 'http://' or 'https://'`),
			},
		},
	})
}

func TestSubjectCleanupBigMaxStoredVersions(t *testing.T) {

	subject_name = "keep-all"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						subjectsToAdd := []int{1, 2, 3, 4, 5}
						err := addSubjectVersions(subject_name, subjectsToAdd)
						if err != nil {
							return err
						}
						return nil
					},
				),
			},
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  cleanup_method = "MAX_STORED_SCHEMAS"
  number_of_schemas_to_keep = 1000
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "subject_name", subject_name),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "cleanup_method", "MAX_STORED_SCHEMAS"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "last_deleted.#", "0"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "cleanup_needed", "false"),
					resource.TestCheckResourceAttr("foxcon_subject_cleanup.test", "rest_endpoint", rest_endpoint),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("foxcon_subject_cleanup.test", "last_updated"),
				),
			},
			{
				Config: cloudProviderConfig + "",
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						expected := "[1,2,3,4,5]"
						err := validateSubjectVersions(subject_name, expected)
						if err != nil {
							return err
						}
						return nil
					},
				),
			},
		},
	})
}

func TestSubjectCleanupZeroMaxStoredVersionsErrorHandling(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
resource "foxcon_subject_cleanup" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "test"
  cleanup_method = "MAX_STORED_SCHEMAS"
  number_of_schemas_to_keep = 0
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				ExpectError: regexp.MustCompile(`Number of schemas must be more than 0 when cleanup_method is set to
'MAX_STORED_SCHEMAS'`),
			},
		},
	})
}
