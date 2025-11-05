// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type schemaRegistryCredentials struct {
	RestEndpoint types.String      `tfsdk:"rest_endpoint"`
	Credentials  *credentialsModel `tfsdk:"credentials"`
}

func (config *schemaRegistryCredentials) ValidateResourceConfig(resp *resource.ValidateConfigResponse) {
	if config.RestEndpoint.IsNull() && config.Credentials == nil {
		// Expected configuration without any schema registry configuration inside for a resource
		return
	}

	if !config.RestEndpoint.IsNull() && config.Credentials == nil {
		// Resource is partially configured. That is not expected
		resp.Diagnostics.AddAttributeError(
			path.Root("credentials"),
			"Missing Required Attribute \"credentials\"",
			"credentials must be set since you set rest_endpoint inside of a resource",
		)
		return
	}

	// Check if any of the fields is set
	anySet := !config.RestEndpoint.IsNull() || !config.Credentials.Key.IsNull() || !config.Credentials.Secret.IsNull()

	// If any is set, all must be set
	if anySet {
		if config.RestEndpoint.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("rest_endpoint"),
				"Missing Required Attribute \"rest_endpoint\"",
				"If any of 'credentials.key', 'credentials.secret', or 'rest_endpoint' is set, all must be set.",
			)
		}
		if config.Credentials.Key.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("credentials.key"),
				"Missing Required Attribute \"credentials.key\"",
				"If any of 'credentials.key', 'credentials.secret', or 'rest_endpoint' is set, all must be set.",
			)
		}
		if config.Credentials.Secret.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("credentials.secret"),
				"Missing Required Attribute \"credentials.secret\"",
				"If any of 'credentials.key', 'credentials.secret', or 'rest_endpoint' is set, all must be set.",
			)
		}
	}
}

func (config *schemaRegistryCredentials) ValidateDataSourceConfig(resp *datasource.ReadResponse) {
	if config.RestEndpoint.IsNull() && config.Credentials == nil {
		// Expected configuration without any schema registry configuration inside for a data source
		return
	}

	if !config.RestEndpoint.IsNull() && config.Credentials == nil {
		// Resource is partially configured. That is not expected
		resp.Diagnostics.AddAttributeError(
			path.Root("credentials"),
			"Missing Required Attribute \"credentials\"",
			"credentials must be set since you set rest_endpoint inside of a data source",
		)
		return
	}

	// Check if any of the fields is set
	anySet := !config.RestEndpoint.IsNull() || !config.Credentials.Key.IsNull() || !config.Credentials.Secret.IsNull()

	// If any is set, all must be set
	if anySet {
		if config.RestEndpoint.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("rest_endpoint"),
				"Missing Required Attribute \"rest_endpoint\"",
				"If any of 'credentials.key', 'credentials.secret', or 'rest_endpoint' is set, all must be set.",
			)
		}
		if config.Credentials.Key.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("credentials.key"),
				"Missing Required Attribute \"credentials.key\"",
				"If any of 'credentials.key', 'credentials.secret', or 'rest_endpoint' is set, all must be set.",
			)
		}
		if config.Credentials.Secret.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("credentials.secret"),
				"Missing Required Attribute \"credentials.secret\"",
				"If any of 'credentials.key', 'credentials.secret', or 'rest_endpoint' is set, all must be set.",
			)
		}
	}
}

func schemaRegistryClientFactory(providerClient *Client, model *schemaRegistryCredentials) (*Client, error) {

	// Local resource config takes precedence over provider client
	if model != nil {
		if !model.RestEndpoint.IsNull() && !model.Credentials.Key.IsNull() && !model.Credentials.Secret.IsNull() {
			schemaAPIClient, err := NewClient(model.RestEndpoint.ValueStringPointer(), model.Credentials.Key.ValueStringPointer(), model.Credentials.Secret.ValueStringPointer())
			if err != nil {
				return nil, err
			}
			return schemaAPIClient, nil
		}
	}
	// Fallback to provider client
	if providerClient != nil {
		return providerClient, nil
	}

	return nil, fmt.Errorf("could not create schema registry client. Make sure rest endpoint and credentials are configured for this resource as there is no schema registry client either configured in the provider settings")
}
