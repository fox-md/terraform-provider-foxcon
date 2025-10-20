// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSubjectVersionsDataSourceRead(t *testing.T) {

	subject_name = "data-source"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
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

	subject_name = "data-source"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/subjects/%s/versions/1?permanent=true", rest_endpoint, subject_name), nil)
					req.SetBasicAuth(api_key, api_secret)
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						panic("failed to send HTTP request:" + err.Error())
					}

					if resp.StatusCode != http.StatusOK {
						panic(fmt.Sprintf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK))
					}
				},
				Config: providerConfig + `
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
