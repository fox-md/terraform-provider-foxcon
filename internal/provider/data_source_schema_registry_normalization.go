// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &schemaRegistryNormalizationDataSource{}
	_ datasource.DataSourceWithConfigure = &schemaRegistryNormalizationDataSource{}
)

// NewSchemaRegistryNormalizationDataSource is a helper function to simplify the provider implementation.
func NewSchemaRegistryNormalizationDataSource() datasource.DataSource {
	return &schemaRegistryNormalizationDataSource{}
}

// schemaRegistryNormalizationDataSource is the data source implementation.
type schemaRegistryNormalizationDataSource struct {
	client *Client
}

// Metadata returns the data source type name.
func (d *schemaRegistryNormalizationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema_registry_normalization"
}

// Schema defines the schema for the data source.
func (d *schemaRegistryNormalizationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads schema registry normalization value.",
		Attributes: map[string]schema.Attribute{
			"rest_endpoint": schema.StringAttribute{
				Optional:    true,
				Description: restEndpointDescription,
			},
			"normalization_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Normalization value",
			},
		},
		Blocks: map[string]schema.Block{
			"credentials": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"key": schema.StringAttribute{
						Optional:    true,
						Description: schemaRegistryKeyDescription,
					},
					"secret": schema.StringAttribute{
						Optional:    true,
						Description: schemaRegistrySecretDescription,
						Sensitive:   true,
					},
				},
			},
		},
	}
}

func (d *schemaRegistryNormalizationDataSource) ValidateConfig(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config schemaRegistryNormalizationDataSourceModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	creds := schemaRegistryCredentials{
		RestEndpoint: config.RestEndpoint,
		Credentials:  config.Credentials,
	}

	creds.ValidateDataSourceConfig(resp)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *schemaRegistryNormalizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config schemaRegistryNormalizationDataSourceModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	creds := schemaRegistryCredentials{
		RestEndpoint: config.RestEndpoint,
		Credentials:  config.Credentials,
	}

	schemaAPIClient, err := schemaRegistryClientFactory(d.client, &creds)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating http client",
			"Could not create http client. Unexpected error: "+err.Error(),
		)
		return
	}

	// Get schema config
	schemaConfig, err := GetSubjectConfig(schemaAPIClient, "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Schema config",
			"Could not read Schema config :"+err.Error(),
		)
		return
	}

	config.Normalize = types.BoolPointerValue(schemaConfig.Normalize)

	// Set refreshed state
	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (d *schemaRegistryNormalizationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	clients, ok := req.ProviderData.(*providerClients)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = clients.SchemaRegistryClient
}

type schemaRegistryNormalizationDataSourceModel struct {
	RestEndpoint types.String      `tfsdk:"rest_endpoint"`
	Credentials  *credentialsModel `tfsdk:"credentials"`
	Normalize    types.Bool        `tfsdk:"normalization_enabled"`
}
