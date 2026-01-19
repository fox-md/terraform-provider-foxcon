// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSetSubjectModeActionWrongRestEndpoint(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
action "foxcon_set_subject_mode" "ro" {
  config {
    rest_endpoint = "httpp://localhost"
    subject_name = "test"
    mode = "READONLY"
  }
}
`,
				ExpectError: regexp.MustCompile(`The value must start with 'http://' or 'https://'`),
			},
		},
	})
}

func TestSetSubjectModeActionWrongMode(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
action "foxcon_set_subject_mode" "ro" {
  config {
    subject_name = "test"
    mode = "WRONGMODE"
  }
}
`,
				ExpectError: regexp.MustCompile(`Attribute mode value must be one of: \[\"READWRITE\" \"READONLY\"
\"READONLY_OVERRIDE\" \"IMPORT\"\]`),
			},
		},
	})
}

func TestSetSubjectModeActionEmptySubject(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
action "foxcon_set_subject_mode" "ro" {
  config {
    subject_name = ""
    mode = "READONLY"
  }
}
`,
				ExpectError: regexp.MustCompile(`Attribute subject_name string length must be at least 1`),
			},
		},
	})
}

func TestSetSubjectModeActionNoCredentials(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
action "foxcon_set_subject_mode" "ro" {
  config {
    rest_endpoint = "http://localhost"
    subject_name = "test"
    mode = "READONLY"
  }
}
`,
				ExpectError: regexp.MustCompile(`Attribute "credentials" must be specified when "rest_endpoint" is specified`),
			},
		},
	})
}

func TestSetSubjectModeActionNoRestEndpoint(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudProviderConfig + `
action "foxcon_set_subject_mode" "ro" {
  config {
    credentials {
	  key = "admin"
      secret = "admin"
    }
    subject_name = "test"
    mode = "READONLY"
  }
}
`,
				ExpectError: regexp.MustCompile(`Attribute "rest_endpoint" must be specified when "credentials.key" is
specified`),
			},
		},
	})
}

// func TestSubjectModeAction(t *testing.T) {

// 	subject_name = "test-action"

// 	mode_ro := "READONLY"
// 	//mode_rw := "READWRITE"

// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
// 			//tfversion.RequireBelow(tfversion.Version1_14_0),
// 			tfversion.SkipBelow(tfversion.Version1_14_0),
// 		},
// 		ExternalProviders: map[string]resource.ExternalProvider{
// 			"confluent": {
// 				Source:            "confluentinc/confluent",
// 				VersionConstraint: "~> 2.0",
// 			},
// 		},
// 		Steps: []resource.TestStep{
// 			// Create and Read testing
// 			{
// 				Config: schemaProviderConfig + `

// provider "confluent" {
//   schema_registry_id            = "lsrc-abc123"
//   schema_registry_rest_endpoint = "` + rest_endpoint + `"
//   schema_registry_api_key       = "` + api_key + `"
//   schema_registry_api_secret    = "` + api_secret + `"
// }

// variable "subject_name" {
//   default = "` + subject_name + `"
// }

// resource "confluent_subject_mode" "test" {
//   subject_name  = var.subject_name
//   mode          = "READONLY"
// }

// locals {
//   json_schema = <<-EOT
// {
//   "$schema": "http://json-schema.org/draft-07/schema#",
//   "$id": "http://example.com/myURI.schema.json",
//   "title": "SampleRecord",
//   "description": "Sample schema to help",
//   "type": "object",
//   "additionalProperties": false,
//   "properties": {
//     "myField1": {
//       "type": "integer",
//       "description": "The integer type is used for integral numbers."
//     }
//   }
// }
// EOT
// }

// resource "confluent_schema" "this" {

//   depends_on = [ confluent_subject_mode.test ]

//   subject_name  = var.subject_name

//   format = "JSON"
//   schema = jsonencode(local.json_schema)

//   lifecycle {
//     action_trigger {
//       events  = [after_create, after_update]
//       actions = [action.foxcon_set_subject_mode.ro]
//     }
//     action_trigger {
//       events  = [before_create, before_update]
//       actions = [action.foxcon_set_subject_mode.rw]
//     }
//   }

// }

// action "foxcon_set_subject_mode" "ro" {
//   config {
//     subject_name = var.subject_name
//     mode = "READONLY"
//   }
// }

// action "foxcon_set_subject_mode" "rw" {
//   config {
//     subject_name = var.subject_name
//     mode = "READWRITE"
//   }
// }
// `,
// 				Check: resource.ComposeTestCheckFunc(
// 					func(s *terraform.State) error {
// 						req, _ := http.NewRequest("GET", fmt.Sprintf("%s/mode/%s", rest_endpoint, subject_name), nil)
// 						req.SetBasicAuth(api_key, api_secret)
// 						resp, err := http.DefaultClient.Do(req)
// 						if err != nil {
// 							return fmt.Errorf("failed to send HTTP request: %s", err)
// 						}

// 						if resp.StatusCode != http.StatusNotFound {
// 							return fmt.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
// 						}

// 						body, _ := io.ReadAll(resp.Body)
// 						strbody := string(body)
// 						defer resp.Body.Close()

// 						if strbody != "{\"mode\": \""+mode_ro+"\"}" {
// 							panic(fmt.Errorf("unexpected body: got '%s', want '%s'", strbody, "{\"mode\": \""+mode_ro+"\"}"))
// 						}

// 						return nil
// 					},
// 				),
// 			},
// 			// 			{
// 			// 				PreConfig: func() {
// 			// 					jsonPayload := []byte("{\"mode\": \"" + mode_rw + "\"}")
// 			// 					req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/mode/%s", rest_endpoint, subject_name), bytes.NewBuffer(jsonPayload))
// 			// 					req.Header.Set("Content-Type", "application/json")
// 			// 					req.SetBasicAuth(api_key, api_secret)
// 			// 					resp, err := http.DefaultClient.Do(req)
// 			// 					if err != nil {
// 			// 						panic("failed to send HTTP request:" + err.Error())
// 			// 					}

// 			// 					body, _ := io.ReadAll(resp.Body)
// 			// 					strbody := string(body)
// 			// 					defer resp.Body.Close()

// 			// 					if resp.StatusCode != http.StatusOK {
// 			// 						panic(fmt.Sprintf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK))
// 			// 					}

// 			// 					if strbody != "{\"mode\": \""+mode_rw+"\"}" {
// 			// 						panic(fmt.Errorf("unexpected body: got '%s', want '%s'", strbody, "{\"mode\": \""+mode_rw+"\"}"))
// 			// 					}
// 			// 				},
// 			// 				Config: schemaProviderConfig + `
// 			// resource "foxcon_subject_normalization" "test" {
// 			//   rest_endpoint = "` + rest_endpoint + `"
// 			//   subject_name = "` + subject_name + `"
// 			//   normalization_enabled = ` + normalization_enabled_true + `
// 			//   credentials {
// 			//     key = "` + api_key + `"
// 			//     secret = "` + api_secret + `"
// 			//   }

// 			//   lifecycle {
// 			//     action_trigger {
// 			//       events  = [after_create, after_update]
// 			//       actions = [action.foxcon_set_subject_mode.rw]
// 			//     }
// 			//   }
// 			// }

// 			// action "foxcon_set_subject_mode" "rw" {
// 			//   config {
// 			//     subject_name = "` + subject_name + `"
// 			//     mode = "` + mode_rw + `"
// 			//   }
// 			// }
// 			// `,
// 			// 				Check: resource.ComposeTestCheckFunc(
// 			// 					func(s *terraform.State) error {
// 			// 						req, _ := http.NewRequest("GET", fmt.Sprintf("%s/mode/%s", rest_endpoint, subject_name), nil)
// 			// 						req.SetBasicAuth(api_key, api_secret)
// 			// 						resp, err := http.DefaultClient.Do(req)
// 			// 						if err != nil {
// 			// 							return fmt.Errorf("failed to send HTTP request: %s", err)
// 			// 						}

// 			// 						if resp.StatusCode != http.StatusNotFound {
// 			// 							return fmt.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
// 			// 						}

// 			// 						body, _ := io.ReadAll(resp.Body)
// 			// 						strbody := string(body)
// 			// 						defer resp.Body.Close()

// 			// 						if strbody != "{\"mode\": \""+mode_rw+"\"}" {
// 			// 							panic(fmt.Errorf("unexpected body: got '%s', want '%s'", strbody, "{\"mode\": \""+mode_rw+"\"}"))
// 			// 						}

// 			// 						return nil
// 			// 					},
// 			// 				),
// 			// 			},
// 		},
// 	})
// }
