// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const ApiEndpoint string = "https://api.confluent.cloud"

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &foxconProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &foxconProvider{
			version: version,
		}
	}
}

// foxconProvider is the provider implementation.
type foxconProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *foxconProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "foxcon"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *foxconProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "Confluent API endpoint. Can be configured using `CONFLUENT_CLOUD_API_ENDPOINT` environment variable. Defaults to: " + ApiEndpoint,
			},
			"cloud_api_key": schema.StringAttribute{
				Optional:    true,
				Description: "Confluent Cloud API Key. Can be configured using `CONFLUENT_CLOUD_API_KEY` environment variable.",
			},
			"cloud_api_secret": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Confluent Cloud API Secret. Can be configured using `CONFLUENT_CLOUD_API_SECRET` environment variable.",
			},
		},
	}
}

// Configure prepares a Confluent API client for data sources and resources.
func (p *foxconProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	tflog.Info(ctx, "Configuring foxcon client")

	// Retrieve provider data from configuration
	var config foxconProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating foxcon")

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.CloudApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("cloud_api_key"),
			"Unknown Confluent API Key",
			"The provider cannot create the HashiCups API client as there is an unknown configuration value for the HashiCups API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_USERNAME environment variable.",
		)
	}

	if config.CloudApiSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("cloud_api_secret"),
			"Unknown Confluent API Secret",
			"The provider cannot create the HashiCups API client as there is an unknown configuration value for the HashiCups API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	cloud_api_key := os.Getenv("CONFLUENT_CLOUD_API_KEY")
	cloud_api_secret := os.Getenv("CONFLUENT_CLOUD_API_SECRET")
	api_endpoint := os.Getenv("CONFLUENT_CLOUD_API_ENDPOINT")

	if !config.ApiEndpoint.IsNull() {
		api_endpoint = config.ApiEndpoint.ValueString()
	}

	if !config.CloudApiKey.IsNull() {
		cloud_api_key = config.CloudApiKey.ValueString()
	}

	if !config.CloudApiSecret.IsNull() {
		cloud_api_secret = config.CloudApiSecret.ValueString()
	}

	if api_endpoint == "" {
		api_endpoint = ApiEndpoint
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if cloud_api_key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("cloud_api_key"),
			"Missing Confluent API Key",
			"The provider cannot create the Confluent API client as there is a missing or empty value for the Confluent API username. "+
				"Set the username value in the configuration or use the CONFLUENT_CLOUD_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if cloud_api_secret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("cloud_api_secret"),
			"Missing Confluent API Secret",
			"The provider cannot create the Confluent API client as there is a missing or empty value for the Confluent API password. "+
				"Set the password value in the configuration or use the CONFLUENT_CLOUD_API_SECRET environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "api_endpoint", api_endpoint)
	ctx = tflog.SetField(ctx, "cloud_api_key", cloud_api_key)
	ctx = tflog.SetField(ctx, "cloud_api_secret", cloud_api_secret)

	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "cloud_api_secret")

	tflog.Debug(ctx, "Creating foxcon client")

	// Create a new HashiCups client using the configuration values
	client, err := NewClient(&api_endpoint, &cloud_api_key, &cloud_api_secret)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Confluent API Client",
			"An unexpected error occurred when creating the Confluent API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Confluent Client Error: "+err.Error(),
		)
		return
	}

	// Make the HashiCups client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured foxcon client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *foxconProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources defines the resources implemented in the provider.
func (p *foxconProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserResource,
		NewInvitationResource,
		NewSubjectNormalizationResource,
		NewSchemaRegistryNormalizationResource,
	}
}

type foxconProviderModel struct {
	ApiEndpoint    types.String `tfsdk:"api_endpoint"`
	CloudApiKey    types.String `tfsdk:"cloud_api_key"`
	CloudApiSecret types.String `tfsdk:"cloud_api_secret"`
}
