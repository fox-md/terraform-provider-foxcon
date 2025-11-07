// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSubjectCleanupLatestHappyFlow(t *testing.T) {

	subject_name = "keep-latest"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
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
		},
	})
}

func TestSubjectCleanupActiveHappyFlow(t *testing.T) {

	subject_name = "keep-active"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
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
		},
	})
}

func TestSubjectCleanupKeepLatestIfOneVersion(t *testing.T) {

	subject_name = "one"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
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
		},
	})
}

func TestSubjectCleanupKeepActiveIfOneVersion(t *testing.T) {

	subject_name = "one"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
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
		},
	})
}

func TestSubjectCleanupKeepActiveIfOneVersionWrongProvider(t *testing.T) {

	subject_name = "one"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
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
		},
	})
}

func TestSubjectCleanupKeepActiveIfOneVersionUsingProvider(t *testing.T) {

	subject_name = "one"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
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
		},
	})
}

func TestSubjectCleanupKeepActiveIfOneVersionNoProviderNoConfig(t *testing.T) {

	subject_name = "one"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
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
			// Create and Read testing
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
		},
	})
}

func TestSubjectCleanupAddSchemaVersionWhenActive(t *testing.T) {

	subject_name = "one-to-two-active"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
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
				PreConfig: func() {
					// Get the root of the Git repo
					gitRootCmd := exec.Command("git", "rev-parse", "--show-toplevel")
					output, err := gitRootCmd.Output()
					if err != nil {
						fmt.Printf("Failed to get git root: %v\n", err)
						return
					}
					gitRoot := strings.TrimSpace(string(output))
					cmd := exec.Command("make", "add2oneactive")
					cmd.Dir = gitRoot
					output, err = cmd.CombinedOutput()
					if err != nil {
						fmt.Printf("Command output:\n%s\n", string(output))
						fmt.Printf("error:%s", err.Error())
						//panic("error:" + err.Error())
					}
				},
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
		},
	})
}

func TestSubjectCleanupAddSchemaVersionWhenLatest(t *testing.T) {

	subject_name = "one-to-two-latest"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
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
				PreConfig: func() {
					// Get the root of the Git repo
					gitRootCmd := exec.Command("git", "rev-parse", "--show-toplevel")
					output, err := gitRootCmd.Output()
					if err != nil {
						fmt.Printf("Failed to get git root: %v\n", err)
						return
					}
					gitRoot := strings.TrimSpace(string(output))
					cmd := exec.Command("make", "add2onelatest")
					cmd.Dir = gitRoot
					output, err = cmd.CombinedOutput()
					if err != nil {
						fmt.Printf("Command output:\n%s\n", string(output))
						fmt.Printf("error:%s", err.Error())
						//panic("error:" + err.Error())
					}
				},
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
		},
	})
}

func TestSubjectCleanupFromActiveToLatestHappyFlow(t *testing.T) {

	subject_name = "switch-cleanup-mode"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
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
			// Update and Read testing
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
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestSubjectCleanupNoCleanupMethod(t *testing.T) {

	subject_name = "dummy"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
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

func TestSubjectCleanupNoCredentials(t *testing.T) {

	subject_name = "dummy"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
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

func TestSubjectCleanupNoRestEndpoint(t *testing.T) {

	subject_name = "dummy"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
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

func TestSubjectCleanupNoSecret(t *testing.T) {

	subject_name = "dummy"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
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

func TestSubjectCleanupNoKey(t *testing.T) {

	subject_name = "dummy"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
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
			// Create and Read testing
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
