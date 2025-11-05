// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	cloudProviderConfig = `
provider "foxcon" {
  cloud_api_key = "test"
  cloud_api_secret = "test"
}
`
	emptyProviderConfig = `
provider "foxcon" {
}
`
	schemaProviderWrongConfig = `
provider "foxcon" {
  schema_registry_rest_endpoint = "http://1.1.1.1"
  schema_registry_api_key = "dummy_value"
  schema_registry_api_secret = "dummy_value"
}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"foxcon": providerserver.NewProtocol6WithError(New("test")()),
	}
)

var rest_endpoint string = "http://localhost:8081"
var api_key string = "admin"
var api_secret string = "admin-secret"
var subject_name string
var subject_name_imported string = "test_import"
var normalization_enabled_true string = "true"
var normalization_enabled_false string = "false"

var schemaProviderConfig = `
provider "foxcon" {
  schema_registry_rest_endpoint = "` + rest_endpoint + `"
  schema_registry_api_key = "` + api_key + `"
  schema_registry_api_secret = "` + api_secret + `"
}
`
