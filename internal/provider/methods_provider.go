// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "github.com/hashicorp/terraform-plugin-framework/resource"

func ValidateCloudApiClient(client *Client, resp *resource.ConfigureResponse) {
	if client == nil {
		resp.Diagnostics.AddError(
			"Missing Confluent API client configuration",
			"The provider cannot create the Confluent API client as there is a missing or empty value for the Confluent API key or secret. "+
				"Set the api_key value in the configuration or use the CONFLUENT_CLOUD_API_KEY environment variable. "+
				"Set the api_secret value in the configuration or use the CONFLUENT_CLOUD_API_SECRET environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
}
