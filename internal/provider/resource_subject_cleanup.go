// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &subjectCleanupResource{}
	_ resource.ResourceWithConfigure   = &subjectCleanupResource{}
	_ resource.ResourceWithImportState = &subjectCleanupResource{}
)

// NewSubjectCleanupResource is a helper function to simplify the provider implementation.
func NewSubjectCleanupResource() resource.Resource {
	return &subjectCleanupResource{}
}

// subjectCleanupResource is the resource implementation.
type subjectCleanupResource struct {
	client *Client
}

// Metadata returns the resource type name.
func (r *subjectCleanupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subject_cleanup"
}

// Schema defines the schema for the resource.
func (r *subjectCleanupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"rest_endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "Schema registry rest endpoint.",
			},
			"subject_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the subject.",
			},
			"cleanup_method": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("KEEP_LATEST_ONLY", "KEEP_ACTIVE_ONLY"),
				},
				Description: "Cleanup method type. Must be set to `KEEP_LATEST_ONLY` or `KEEP_ACTIVE_ONLY`.",
			},
			"latest_schema_version": schema.Int32Attribute{
				Computed:    true,
				Description: "Last schema version number.",
			},
			"cleanup_needed": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Toggle to control whether clean-up in needed. No need to set it manually.",
				Default:     booldefault.StaticBool(false),
			},
			"last_deleted": schema.ListAttribute{
				ElementType: types.Int32Type,
				Computed:    true,
				Description: "List of schema versions deleted on the last apply execution.",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp of the last apply execution.",
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
		MarkdownDescription: "Deletes schema versions depending on the configured clean-up method.",
	}
}

type subjectCleanupResourceModel struct {
	RestEndpoint      types.String      `tfsdk:"rest_endpoint"`
	SubjectName       types.String      `tfsdk:"subject_name"`
	Credentials       *credentialsModel `tfsdk:"credentials"`
	LastSchemaVersion types.Int32       `tfsdk:"latest_schema_version"`
	CleanupNeeded     types.Bool        `tfsdk:"cleanup_needed"`
	CleanupMethod     types.String      `tfsdk:"cleanup_method"`
	LastDeleted       types.List        `tfsdk:"last_deleted"`
	LastUpdated       types.String      `tfsdk:"last_updated"`
}

func (r *subjectCleanupResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config subjectCleanupResourceModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	creds := schemaRegistryCredentials{
		RestEndpoint: config.RestEndpoint,
		Credentials:  config.Credentials,
	}

	creds.ValidateResourceConfig(resp)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
// Create a new resource.
func (r *subjectCleanupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan subjectCleanupResourceModel
	var latestVersion int
	var deleteCandidates []int
	var subjectVersions schemaVersions

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	creds := schemaRegistryCredentials{
		RestEndpoint: plan.RestEndpoint,
		Credentials:  plan.Credentials,
	}

	schemaAPIClient, err := schemaRegistryClientFactory(r.client, &creds)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating http client",
			"Could not create http client. Unexpected error: "+err.Error(),
		)
		return
	}

	subjectVersions.client = schemaAPIClient

	err = subjectVersions.get(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting subject versions",
			"Could not get subject versions. Unexpected error: "+err.Error(),
		)
		return
	}

	err = subjectVersions.cleanSoftDeleted(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting soft-deleted versions",
			"Could not delete soft-deleted versions. Unexpected error: "+err.Error(),
		)
		return
	}

	deleteCandidates = (*subjectVersions.softDeleted)

	if plan.CleanupMethod == types.StringValue("KEEP_LATEST_ONLY") {
		err := subjectVersions.cleanActiveNoneLatest(ctx, plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error deleting active versions",
				"Could not delete active versions. Unexpected error: "+err.Error(),
			)
			return
		}
		deleteCandidates = (*subjectVersions.all)[:len(*subjectVersions.all)-1]
	}

	latestVersion = (*subjectVersions.all)[len(*subjectVersions.all)-1]

	var lastDeleted []attr.Value
	for _, id := range deleteCandidates {
		lastDeleted = append(lastDeleted, types.Int32Value(int32(id)))
	}

	plan.LastSchemaVersion = types.Int32Value(int32(latestVersion))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))
	plan.CleanupNeeded = types.BoolValue(false)
	plan.LastDeleted, diags = types.ListValue(types.Int32Type, lastDeleted)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
// Read resource information.
func (r *subjectCleanupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state subjectCleanupResourceModel
	var deleteCandidates []int
	var subjectVersions schemaVersions

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	creds := schemaRegistryCredentials{
		RestEndpoint: state.RestEndpoint,
		Credentials:  state.Credentials,
	}

	schemaAPIClient, err := schemaRegistryClientFactory(r.client, &creds)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating http client",
			"Could not create http client. Unexpected error: "+err.Error(),
		)
		return
	}

	subjectVersions.client = schemaAPIClient
	err = subjectVersions.get(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting subject versions",
			"Could not get subject versions. Unexpected error: "+err.Error(),
		)
		return
	}

	if state.CleanupMethod == types.StringValue("KEEP_LATEST_ONLY") {
		deleteCandidates = (*subjectVersions.all)[:len(*subjectVersions.all)-1]
	}

	if state.CleanupMethod == types.StringValue("KEEP_ACTIVE_ONLY") {
		deleteCandidates = (*subjectVersions.softDeleted)
	}

	if len(deleteCandidates) > 0 {
		state.CleanupNeeded = types.BoolValue(true)
	} else {
		state.CleanupNeeded = types.BoolValue(false)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *subjectCleanupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan subjectCleanupResourceModel
	var latestVersion int
	var deleteCandidates []int
	var subjectVersions schemaVersions

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	creds := schemaRegistryCredentials{
		RestEndpoint: plan.RestEndpoint,
		Credentials:  plan.Credentials,
	}

	schemaAPIClient, err := schemaRegistryClientFactory(r.client, &creds)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating http client",
			"Could not create http client. Unexpected error: "+err.Error(),
		)
		return
	}

	subjectVersions.client = schemaAPIClient
	err = subjectVersions.get(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting subject versions",
			"Could not get subject versions. Unexpected error: "+err.Error(),
		)
		return
	}

	// Execute flow for the KEEP_ACTIVE_ONLY. Needed by both cleanup methods.
	err = subjectVersions.cleanSoftDeleted(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting soft-deleted versions",
			"Could not delete soft-deleted versions. Unexpected error: "+err.Error(),
		)
		return
	}

	deleteCandidates = (*subjectVersions.softDeleted)

	if plan.CleanupMethod == types.StringValue("KEEP_LATEST_ONLY") {
		err := subjectVersions.cleanActiveNoneLatest(ctx, plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error deleting active versions",
				"Could not delete active versions. Unexpected error: "+err.Error(),
			)
			return
		}

		deleteCandidates = (*subjectVersions.all)[:len(*subjectVersions.all)-1]
	}

	latestVersion = (*subjectVersions.all)[len(*subjectVersions.all)-1]

	// Map response body to schema and populate Computed attribute values
	var lastDeleted []attr.Value
	for _, id := range deleteCandidates {
		lastDeleted = append(lastDeleted, types.Int32Value(int32(id)))
	}

	plan.LastSchemaVersion = types.Int32Value(int32(latestVersion))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))
	plan.CleanupNeeded = types.BoolValue(false)
	plan.LastDeleted, diags = types.ListValue(types.Int32Type, lastDeleted)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *subjectCleanupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state subjectCleanupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Deleting clean-up resource with effecting subject  %s", state.SubjectName.ValueString()))
}

// Configure adds the provider configured client to the resource.
func (r *subjectCleanupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = clients.SchemaRegistryClient
}

func (r *subjectCleanupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.AddError(
		"Import not implemented",
		"Import for this resource is not available since the resource itself does not create any objects.",
	)
}
