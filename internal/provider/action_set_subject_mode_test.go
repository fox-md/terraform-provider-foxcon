// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
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

func TestSubjectModeAction(t *testing.T) {

	subject_name = "test-action"

	mode_ro := "READONLY"
	mode_rw := "READWRITE"

	t.Setenv("SCHEMA_REGISTRY_ID", "lsrc-abc123")
	t.Setenv("SCHEMA_REGISTRY_REST_ENDPOINT", rest_endpoint)
	t.Setenv("SCHEMA_REGISTRY_API_KEY", api_key)
	t.Setenv("SCHEMA_REGISTRY_API_SECRET", api_secret)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_15_0),
		},
		ExternalProviders: map[string]resource.ExternalProvider{
			"confluent": {
				Source:            "confluentinc/confluent",
				VersionConstraint: "~> 2.0",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: schemaProviderConfig,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						err := setSubjectMode(subject_name, mode_ro)
						if err != nil {
							return err
						}

						err = getSubjectMode(subject_name, mode_ro)
						if err != nil {
							return err
						}

						return nil
					},
				),
			},
			{
				Config: schemaProviderConfig + `

			variable "subject_name" {
			  default = "` + subject_name + `"
			}

			locals {
			  json_schema = <<-EOT
			{
			  "$schema": "http://json-schema.org/draft-07/schema#",
			  "$id": "http://example.com/myURI.schema.json",
			  "title": "SampleRecord",
			  "description": "Sample schema to help",
			  "type": "object",
			  "additionalProperties": false,
			  "properties": {
			    "myField1": {
			      "type": "integer",
			      "description": "The integer type is used for integral numbers."
			    }
			  }
			}
			EOT
			}

			resource "confluent_schema" "this" {

			  subject_name  = var.subject_name

			  format = "JSON"
			  schema = jsonencode(local.json_schema)

			  lifecycle {
			    action_trigger {
			      events  = [after_create, after_update]
			      actions = [action.foxcon_set_subject_mode.ro]
			    }
			    action_trigger {
			      events  = [before_create, before_update]
			      actions = [action.foxcon_set_subject_mode.rw]
			    }
			  }

			}

			action "foxcon_set_subject_mode" "ro" {
			  config {
			    subject_name = var.subject_name
			    mode = "READONLY"
			  }
			}

			action "foxcon_set_subject_mode" "rw" {
			  config {
			    subject_name = var.subject_name
			    mode = "READWRITE"
			  }
			}
			`,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						err := getSubjectMode(subject_name, mode_ro)
						if err != nil {
							return err
						}

						err = setSubjectMode(subject_name, mode_rw)
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
