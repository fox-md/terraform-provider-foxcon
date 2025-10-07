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

	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &subjectNormalizationResource{}
	_ resource.ResourceWithConfigure   = &subjectNormalizationResource{}
	_ resource.ResourceWithImportState = &subjectNormalizationResource{}
)

// NewsubjectNormalizationResource is a helper function to simplify the provider implementation.
func NewSubjectNormalizationResource() resource.Resource {
	return &subjectNormalizationResource{}
}

// subjectNormalizationResource is the resource implementation.
type subjectNormalizationResource struct {
	client *Client
}

// Metadata returns the resource type name.
func (r *subjectNormalizationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subject_normalization"
}

// Schema defines the schema for the resource.
func (r *subjectNormalizationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"rest_endpoint": schema.StringAttribute{
				Required:    true,
				Description: "Schema registry rest endpoint",
			},
			"subject_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the subject",
			},
			"normalization_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Normalization toggle",
				//Default:  booldefault.StaticBool(false),
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

type subjectNormalizationResourceModel struct {
	RestEndpoint types.String      `tfsdk:"rest_endpoint"`
	SubjectName  types.String      `tfsdk:"subject_name"`
	Normalize    types.Bool        `tfsdk:"normalization_enabled"`
	Credentials  *credentialsModel `tfsdk:"credentials"`
	LastUpdated  types.String      `tfsdk:"last_updated"`
}

// type credentialsModel struct {
// 	Key    types.String `tfsdk:"key"`
// 	Secret types.String `tfsdk:"secret"`
// }

// Create creates the resource and sets the initial Terraform state.
// Create a new resource.
func (r *subjectNormalizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan subjectNormalizationResourceModel
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

	if plan.Normalize.IsUnknown() {
		resp.Diagnostics.AddError(
			"Error reading normalization value",
			"Normalization value is unknown.",
		)
		return
	}

	var normalizationPayload = NormalizeRequest{
		Normalize: plan.Normalize.ValueBoolPointer(),
	}

	// Set Normalization
	subjectConfig, err := SetNormalization(schemaAPIClient, plan.SubjectName.ValueString(), normalizationPayload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting normalization",
			"Could not set normalization unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.Normalize = types.BoolPointerValue(subjectConfig.Normalize)
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
func (r *subjectNormalizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state subjectNormalizationResourceModel
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

	// Get subject config
	subjectConfig, err := GetSchemaConfig(schemaAPIClient, state.SubjectName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Subject config",
			"Could not read Subject config "+state.SubjectName.ValueString()+": "+err.Error(),
		)
		return
	}

	if subjectConfig == nil {
		tflog.Debug(ctx, fmt.Sprintf("%s subject config does not exist in Confluent. Removing resource from state file.", state.SubjectName.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	// Overwrite items with refreshed state

	state.Normalize = types.BoolPointerValue(subjectConfig.Normalize)
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *subjectNormalizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan subjectNormalizationResourceModel
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

	var normalizationPayload NormalizeRequest

	if plan.Normalize.IsNull() {
		normalizationPayload.Normalize = nil
	} else {
		normalizationPayload.Normalize = plan.Normalize.ValueBoolPointer()
	}

	subjectConfig, err := SetNormalization(schemaAPIClient, plan.SubjectName.ValueString(), normalizationPayload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Subject config",
			"Could not read Subject config "+plan.SubjectName.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.Normalize = types.BoolPointerValue(subjectConfig.Normalize)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *subjectNormalizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state subjectNormalizationResourceModel
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

	// Get schema registry config
	schemaRegistryConfig, err := GetSchemaConfig(schemaAPIClient, "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Schema config",
			"Could not read Schema config:"+err.Error(),
		)
		return
	}

	// Get subject config
	subjectConfig, err := GetSchemaConfig(schemaAPIClient, state.SubjectName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Subject config",
			"Could not read Subject config "+state.SubjectName.ValueString()+": "+err.Error(),
		)
		return
	}

	if *subjectConfig.CompatibilityLevel == *schemaRegistryConfig.CompatibilityLevel && countAttr(subjectConfig) < 3 {
		// Delete subject config as CompatibilityLevels are identical (being inherited) and the second remaining value is normalize
		tflog.Debug(ctx, fmt.Sprintf("Deleting entire %s subject config", state.SubjectName.ValueString()))
		err = DeleteSubjectConfig(schemaAPIClient, state.SubjectName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error deleting subject configuration",
				"Could not delete subject configuration: "+err.Error(),
			)
			return
		}
	} else {
		// Delete normalization only
		tflog.Debug(ctx, fmt.Sprintf("Deleting normalization toggle for %s subject config", state.SubjectName.ValueString()))
		var normalizationPayload = NormalizeRequest{
			Normalize: nil,
		}

		_, err = SetNormalization(schemaAPIClient, state.SubjectName.ValueString(), normalizationPayload)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error setting normalization value to null",
				"Could not set normalization value to null: "+err.Error(),
			)
			return
		}
	}

}

// Configure adds the provider configured client to the resource.
func (r *subjectNormalizationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *subjectNormalizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subject_name"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("rest_endpoint"), rest_endpoint)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("credentials"), credentials)...)
}

func countAttr(resp *SchemaConfigResponse) int {
	count := 0
	if resp.Alias != nil {
		count++
	}
	if resp.Normalize != nil {
		count++
	}
	if resp.CompatibilityLevel != nil {
		count++
	}
	if resp.CompatibilityGroup != nil {
		count++
	}
	if resp.DefaultMetadata != nil {
		count++
	}
	if resp.OverrideMetadata != nil {
		count++
	}
	if resp.DefaultRuleSet != nil {
		count++
	}
	if resp.OverrideRuleSet != nil {
		count++
	}
	return count
}
