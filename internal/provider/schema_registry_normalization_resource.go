// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	// "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &schemaRegistryNormalizationResource{}
	_ resource.ResourceWithConfigure   = &schemaRegistryNormalizationResource{}
	_ resource.ResourceWithImportState = &schemaRegistryNormalizationResource{}
)

// NewschemaRegistryNormalizationResource is a helper function to simplify the provider implementation.
func NewSchemaRegistryNormalizationResource() resource.Resource {
	return &schemaRegistryNormalizationResource{}
}

// schemaRegistryNormalizationResource is the resource implementation.
type schemaRegistryNormalizationResource struct {
	client *Client
}

// Metadata returns the resource type name.
func (r *schemaRegistryNormalizationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema_registry_normalization"
}

// Schema defines the schema for the resource.
func (r *schemaRegistryNormalizationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"rest_endpoint": schema.StringAttribute{
				Required:    true,
				Description: "Schema registry rest endpoint",
			},
			"normalization_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Normalization toggle",
			},
			// "credentials": schema.SingleNestedAttribute{
			// 	Required: true,
			// 	Attributes: map[string]schema.Attribute{
			// 		"key": schema.StringAttribute{
			// 			Required: true,
			// 		},
			// 		"secret": schema.StringAttribute{
			// 			Required:  true,
			// 			Sensitive: true,
			// 		},
			// 	},
			// },
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"credentials": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"key": schema.StringAttribute{
						Required: true,
					},
					"secret": schema.StringAttribute{
						Required:  true,
						Sensitive: true,
					},
				},
			},
		},
	}
}

type schemaRegistryNormalizationResourceModel struct {
	RestEndpoint types.String      `tfsdk:"rest_endpoint"`
	Normalize    types.Bool        `tfsdk:"normalization_enabled"`
	Credentials  *credentialsModel `tfsdk:"credentials"`
	LastUpdated  types.String      `tfsdk:"last_updated"`
}

// Create creates the resource and sets the initial Terraform state.
// Create a new resource.
func (r *schemaRegistryNormalizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan schemaRegistryNormalizationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	schemaAPIClient, err := NewClient(plan.RestEndpoint.ValueStringPointer(), plan.Credentials.Key.ValueStringPointer(), plan.Credentials.Secret.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating http client",
			"Could not create http client. Unexpected error: "+err.Error(),
		)
		return
	}

	var normalizationPayload = NormalizeRequest{
		Normalize: plan.Normalize.ValueBoolPointer(),
	}

	// Set Normalization
	schemaConfig, err := SetNormalization(schemaAPIClient, "", normalizationPayload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting normalization",
			"Could not set normalization unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.Normalize = types.BoolPointerValue(schemaConfig.Normalize)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
// Read resource information.
func (r *schemaRegistryNormalizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state schemaRegistryNormalizationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	schemaAPIClient, err := NewClient(state.RestEndpoint.ValueStringPointer(), state.Credentials.Key.ValueStringPointer(), state.Credentials.Secret.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating http client",
			"Could not create http client. Unexpected error: "+err.Error(),
		)
		return
	}

	// Get schema config
	schemaConfig, err := GetSchemaConfig(schemaAPIClient, "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Schema config",
			"Could not read Schema config :"+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.Normalize = types.BoolPointerValue(schemaConfig.Normalize)
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *schemaRegistryNormalizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan schemaRegistryNormalizationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	schemaAPIClient, err := NewClient(plan.RestEndpoint.ValueStringPointer(), plan.Credentials.Key.ValueStringPointer(), plan.Credentials.Secret.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating http client",
			"Could not create http client. Unexpected error: "+err.Error(),
		)
		return
	}

	var normalizationPayload = NormalizeRequest{
		Normalize: plan.Normalize.ValueBoolPointer(),
	}

	schemaConfig, err := SetNormalization(schemaAPIClient, "", normalizationPayload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Schema config",
			"Could not read Schema config :"+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.Normalize = types.BoolPointerValue(schemaConfig.Normalize)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *schemaRegistryNormalizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state schemaRegistryNormalizationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	schemaAPIClient, err := NewClient(state.RestEndpoint.ValueStringPointer(), state.Credentials.Key.ValueStringPointer(), state.Credentials.Secret.ValueStringPointer())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating http client",
			"Could not create http client. Unexpected error: "+err.Error(),
		)
		return
	}

	// Set normalization to null
	tflog.Debug(ctx, fmt.Sprintf("Deleting normalization toggle for schema %s", state.RestEndpoint.ValueString()))

	var normalizationPayload = NormalizeRequest{
		Normalize: nil,
	}

	_, err = SetNormalization(schemaAPIClient, "", normalizationPayload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting normalization value to null",
			"Could not set normalization value to null: "+err.Error(),
		)
		return
	}

}

// Configure adds the provider configured client to the resource.
func (r *schemaRegistryNormalizationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *schemaRegistryNormalizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve subject config and save it to state
	// resource.ImportStatePassthroughID(ctx, path.Root("subject_name"), req, resp)

	// parts := strings.Split(req.ID, "/")
	// if len(parts) != 2 {
	//     resp.Diagnostics.AddError(
	//         "Unexpected import identifier",
	//         "Expected format: <Schema Registry cluster ID>/<Subject name>",
	//     )
	//     return
	// }

	// resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("schema_registry_id"), parts[0])...)
	// resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subject_name"), parts[1])...)

	if os.Getenv("IMPORT_SCHEMA_REGISTRY_REST_ENDPOINT") == "" {
		resp.Diagnostics.AddError(
			"Import error",
			"'IMPORT_SCHEMA_REGISTRY_REST_ENDPOINT' environment variable is not configured",
		)
		return
	}

	if os.Getenv("IMPORT_SCHEMA_REGISTRY_API_KEY") == "" {
		resp.Diagnostics.AddError(
			"Import error",
			"'IMPORT_SCHEMA_REGISTRY_API_KEY' environment variable is not configured",
		)
		return
	}

	if os.Getenv("IMPORT_SCHEMA_REGISTRY_API_SECRET") == "" {
		resp.Diagnostics.AddError(
			"Import error",
			"'IMPORT_SCHEMA_REGISTRY_API_SECRET' environment variable is not configured",
		)
		return
	}

	var credentials = &credentialsModel{
		Key:    types.StringValue(os.Getenv("IMPORT_SCHEMA_REGISTRY_API_KEY")),
		Secret: types.StringValue(os.Getenv("IMPORT_SCHEMA_REGISTRY_API_SECRET")),
	}

	var rest_endpoint = types.StringValue(os.Getenv("IMPORT_SCHEMA_REGISTRY_REST_ENDPOINT"))

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("rest_endpoint"), rest_endpoint)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("credentials"), credentials)...)
}
