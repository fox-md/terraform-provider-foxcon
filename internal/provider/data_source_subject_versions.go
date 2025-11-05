// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &subjectVersionsDataSource{}
	_ datasource.DataSourceWithConfigure = &subjectVersionsDataSource{}
)

// NewSubjectVersionsDataSource is a helper function to simplify the provider implementation.
func NewSubjectVersionsDataSource() datasource.DataSource {
	return &subjectVersionsDataSource{}
}

// subjectVersionsDataSource is the data source implementation.
type subjectVersionsDataSource struct {
	client *Client
}

// Metadata returns the data source type name.
func (d *subjectVersionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subject_versions"
}

// Schema defines the schema for the data source.
func (d *subjectVersionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"rest_endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "Schema registry rest endpoint",
			},
			"subject_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the subject",
			},
			"latest": schema.Int32Attribute{
				Computed:    true,
				Description: "Latest schema version number.",
			},
			"active": schema.ListAttribute{
				ElementType: types.Int32Type,
				Computed:    true,
				Description: "List of all active versions.",
			},
			"soft_deleted": schema.ListAttribute{
				ElementType: types.Int32Type,
				Computed:    true,
				Description: "List of all soft-deleted versions.",
			},
			"all": schema.ListAttribute{
				ElementType: types.Int32Type,
				Computed:    true,
				Description: "List of all schema versions (active and soft-deleted).",
			},
		},
		Blocks: map[string]schema.Block{
			"credentials": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"key": schema.StringAttribute{
						Optional: true,
					},
					"secret": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
				},
			},
		},
		MarkdownDescription: "Reads subject schema versions",
	}
}

func (d *subjectVersionsDataSource) ValidateConfig(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config subjectVersionsDataSourceModel

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
func (d *subjectVersionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config subjectVersionsDataSourceModel
	var subjectVersions schemaVersions

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

	subjectVersions.client = schemaAPIClient

	var subject_config = subjectCleanupResourceModel{
		RestEndpoint: config.RestEndpoint,
		Credentials:  config.Credentials,
		SubjectName:  config.SubjectName,
	}

	err = subjectVersions.get(subject_config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating get subject versions",
			"Could not get subject versions. Unexpected error: "+err.Error(),
		)
		return
	}

	latestVersion := (*subjectVersions.all)[len(*subjectVersions.all)-1]
	config.LatestSchemaVersion = types.Int32Value(int32(latestVersion))

	var all []attr.Value
	for _, id := range *subjectVersions.all {
		all = append(all, types.Int32Value(int32(id)))
	}
	config.AllVersions, _ = types.ListValue(types.Int32Type, all)

	var active []attr.Value
	for _, id := range *subjectVersions.active {
		active = append(active, types.Int32Value(int32(id)))
	}
	config.ActiveVersions, _ = types.ListValue(types.Int32Type, active)

	var softDeleted []attr.Value
	for _, id := range *subjectVersions.softDeleted {
		softDeleted = append(softDeleted, types.Int32Value(int32(id)))
	}
	config.SoftDeletedVersions, _ = types.ListValue(types.Int32Type, softDeleted)

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (d *subjectVersionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type subjectVersionsDataSourceModel struct {
	RestEndpoint        types.String      `tfsdk:"rest_endpoint"`
	SubjectName         types.String      `tfsdk:"subject_name"`
	Credentials         *credentialsModel `tfsdk:"credentials"`
	LatestSchemaVersion types.Int32       `tfsdk:"latest"`
	AllVersions         types.List        `tfsdk:"all"`
	ActiveVersions      types.List        `tfsdk:"active"`
	SoftDeletedVersions types.List        `tfsdk:"soft_deleted"`
}
