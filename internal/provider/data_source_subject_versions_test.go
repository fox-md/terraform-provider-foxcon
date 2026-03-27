// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestSubjectVersionsDataSourceRead(t *testing.T) {

	subject_name = "data-source-read"

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
data "foxcon_subject_versions" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.foxcon_subject_versions.test", "latest", "5"),
					resource.TestCheckResourceAttr("data.foxcon_subject_versions.test", "soft_deleted.#", "2"),
					resource.TestCheckResourceAttr("data.foxcon_subject_versions.test", "active.#", "3"),
					resource.TestCheckResourceAttr("data.foxcon_subject_versions.test", "all.#", "5"),
				),
			},
		},
	})
}

func TestSubjectVersionsDataSourceReadAfterV1Delete(t *testing.T) {

	subject_name = "data-source-hard"

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

						req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/subjects/%s/versions/1?permanent=true", rest_endpoint, subject_name), nil)
						req.SetBasicAuth(api_key, api_secret)
						resp, err := http.DefaultClient.Do(req)
						if err != nil {
							return fmt.Errorf("failed to send HTTP request: %s", err.Error())
						}

						if resp.StatusCode != http.StatusOK {
							return fmt.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
						}

						return nil
					},
				),
			},
			{
				Config: cloudProviderConfig + `
data "foxcon_subject_versions" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.foxcon_subject_versions.test", "latest", "5"),
					resource.TestCheckResourceAttr("data.foxcon_subject_versions.test", "soft_deleted.#", "1"),
					resource.TestCheckResourceAttr("data.foxcon_subject_versions.test", "soft_deleted.0", "2"),
					resource.TestCheckResourceAttr("data.foxcon_subject_versions.test", "active.#", "3"),
					resource.TestCheckResourceAttr("data.foxcon_subject_versions.test", "all.#", "4"),
				),
			},
		},
	})
}

func TestSubjectVersionsDataSourceReadProviderMigration111To121Setup(t *testing.T) {

	subject_name = "data-source-migration"

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
data "foxcon_subject_versions" "test" {
  rest_endpoint = "` + rest_endpoint + `"
  subject_name = "` + subject_name + `"
  credentials {
    key = "` + api_key + `"
    secret = "` + api_secret + `"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.foxcon_subject_versions.test", "rest_endpoint", rest_endpoint),
					resource.TestCheckResourceAttr("data.foxcon_subject_versions.test", "latest", "5"),
					resource.TestCheckResourceAttr("data.foxcon_subject_versions.test", "soft_deleted.#", "2"),
					resource.TestCheckResourceAttr("data.foxcon_subject_versions.test", "active.#", "3"),
					resource.TestCheckResourceAttr("data.foxcon_subject_versions.test", "all.#", "5"),
				),
			},
			{
				Config: schemaProviderConfig + `
data "foxcon_subject_versions" "test" {
  subject_name = "` + subject_name + `"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("data.foxcon_subject_versions.test", "rest_endpoint"),
					resource.TestCheckResourceAttr("data.foxcon_subject_versions.test", "latest", "5"),
					resource.TestCheckResourceAttr("data.foxcon_subject_versions.test", "soft_deleted.#", "2"),
					resource.TestCheckResourceAttr("data.foxcon_subject_versions.test", "active.#", "3"),
					resource.TestCheckResourceAttr("data.foxcon_subject_versions.test", "all.#", "5"),
				),
			},
		},
	})
}

func TestDataSourceReadSubjectInvalidRestEndpointErrorHandling(t *testing.T) {

	subject_name = "test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
data "foxcon_subject_versions" "test" {
  rest_endpoint = "localhost"
  subject_name = "` + subject_name + `"
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

func TestDataSourceReadSubjectEmptySubjectErrorHandling(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
data "foxcon_subject_versions" "test" {
  rest_endpoint = "http://localhost"
  subject_name = ""
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

func TestDataSourceReadSubjectNoKeyErrorHandling(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
data "foxcon_subject_versions" "test" {
  rest_endpoint = "http://localhost"
  subject_name = "test"
  credentials {
    secret = "` + api_secret + `"
  }
}
`,
				ExpectError: regexp.MustCompile(`Attribute \"credentials.key\" must be specified when \"credentials.secret\" is
specified`),
			},
		},
	})
}

func TestDataSourceReadSubjectNoSecretErrorHandling(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
data "foxcon_subject_versions" "test" {
  rest_endpoint = "http://localhost"
  subject_name = "test"
  credentials {
    key = "` + api_key + `"
  }
}
`,
				ExpectError: regexp.MustCompile(`Attribute \"credentials.secret\" must be specified when \"credentials.key\" is
specified`),
			},
		},
	})
}
