// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type SchemaConfigMetadata struct {
	Properties map[string]interface{} `json:"properties"`
}

type SchemaConfigResponse struct {
	Alias              *string               `json:"alias"`
	Normalize          *bool                 `json:"normalize"`
	CompatibilityLevel *string               `json:"compatibilityLevel"`
	CompatibilityGroup *string               `json:"compatibilityGroup"`
	DefaultMetadata    *SchemaConfigMetadata `json:"defaultMetadata"`
	OverrideMetadata   *SchemaConfigMetadata `json:"overrideMetadata"`
	DefaultRuleSet     *SchemaConfigMetadata `json:"defaultRuleSet"`
	OverrideRuleSet    *SchemaConfigMetadata `json:"overrideRuleSet"`
}

type credentialsModel struct {
	Key    types.String `tfsdk:"key"`
	Secret types.String `tfsdk:"secret"`
}
