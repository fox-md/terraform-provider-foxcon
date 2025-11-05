// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"sort"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

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
				Optional:    true,
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
		MarkdownDescription: "Sets subject normalization value",
	}
}

type subjectNormalizationResourceModel struct {
	RestEndpoint types.String      `tfsdk:"rest_endpoint"`
	SubjectName  types.String      `tfsdk:"subject_name"`
	Normalize    types.Bool        `tfsdk:"normalization_enabled"`
	Credentials  *credentialsModel `tfsdk:"credentials"`
	LastUpdated  types.String      `tfsdk:"last_updated"`
}

func (r *subjectNormalizationResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config subjectNormalizationResourceModel

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
func (r *subjectNormalizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan subjectNormalizationResourceModel
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
	subjectConfig, err := SetSubjectConfig(schemaAPIClient, plan.SubjectName.ValueString(), normalizationPayload)
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

	// Get subject config
	subjectConfig, err := GetSubjectConfig(schemaAPIClient, state.SubjectName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Subject config",
			"Could not read Subject config "+state.SubjectName.ValueString()+": "+err.Error(),
		)
		return
	}

	if subjectConfig == nil {
		state.Normalize = types.BoolNull()
	} else {
		state.Normalize = types.BoolPointerValue(subjectConfig.Normalize)
	}

	// Overwrite items with refreshed state
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
	var normalizationPayload NormalizeRequest
	var attrs []string

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

	// Get schema registry config
	schemaRegistryConfig, err := GetSubjectConfig(schemaAPIClient, "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Schema config",
			"Could not read Schema config:"+err.Error(),
		)
		return
	}

	// Get subject config
	subjectConfig, err := GetSubjectConfig(schemaAPIClient, plan.SubjectName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Subject config",
			"Could not read Subject config "+plan.SubjectName.ValueString()+": "+err.Error(),
		)
		return
	}

	if subjectConfig != nil {
		attrs = parseResponseAttrs(subjectConfig)
	}

	if plan.Normalize.IsNull() {
		normalizationPayload.Normalize = nil
	} else {
		normalizationPayload.Normalize = plan.Normalize.ValueBoolPointer()
	}

	if plan.Normalize.IsNull() {
		if *subjectConfig.CompatibilityLevel == *schemaRegistryConfig.CompatibilityLevel &&
			(reflect.DeepEqual(attrs, []string{"compatibilityLevel", "normalize"}) ||
				reflect.DeepEqual(attrs, []string{"compatibilityLevel"})) {

			// Delete subject config as CompatibilityLevels are identical (being inherited) and the second remaining value must be normalize
			tflog.Debug(ctx, fmt.Sprintf("Deleting entire %s subject config", plan.SubjectName.ValueString()))
			err = DeleteSubjectConfig(schemaAPIClient, plan.SubjectName.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error deleting subject configuration",
					"Could not delete subject configuration: "+err.Error(),
				)
				return
			}
		}
		plan.Normalize = types.BoolNull()
	} else {
		subjectConfig, err := SetSubjectConfig(schemaAPIClient, plan.SubjectName.ValueString(), normalizationPayload)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Setting Subject config",
				"Could not set Subject config "+plan.SubjectName.ValueString()+": "+err.Error(),
			)
			return
		}
		plan.Normalize = types.BoolPointerValue(subjectConfig.Normalize)
	}

	// Map response body to schema and populate Computed attribute values
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
	var attrs []string

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

	// Get schema registry config
	schemaRegistryConfig, err := GetSubjectConfig(schemaAPIClient, "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Schema config",
			"Could not read Schema config:"+err.Error(),
		)
		return
	}

	// Get subject config
	subjectConfig, err := GetSubjectConfig(schemaAPIClient, state.SubjectName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Subject config",
			"Could not read Subject config "+state.SubjectName.ValueString()+": "+err.Error(),
		)
		return
	}

	if subjectConfig == nil {
		return
	} else {
		attrs = parseResponseAttrs(subjectConfig)
	}

	if *subjectConfig.CompatibilityLevel == *schemaRegistryConfig.CompatibilityLevel &&
		(reflect.DeepEqual(attrs, []string{"compatibilityLevel", "normalize"}) ||
			reflect.DeepEqual(attrs, []string{"compatibilityLevel"})) {

		// Delete subject config as CompatibilityLevels are identical (being inherited) and the second remaining value must be normalize
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

		_, err = SetSubjectConfig(schemaAPIClient, state.SubjectName.ValueString(), normalizationPayload)
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

func parseResponseAttrs(resp *SchemaConfigResponse) []string {
	var attrs []string
	if resp.Alias != nil {
		attrs = append(attrs, "alias")
	}
	if resp.Normalize != nil {
		attrs = append(attrs, "normalize")
	}
	if resp.CompatibilityLevel != nil {
		attrs = append(attrs, "compatibilityLevel")
	}
	if resp.CompatibilityGroup != nil {
		attrs = append(attrs, "compatibilityGroup")
	}
	if resp.DefaultMetadata != nil {
		attrs = append(attrs, "defaultMetadata")
	}
	if resp.OverrideMetadata != nil {
		attrs = append(attrs, "overrideMetadata")
	}
	if resp.DefaultRuleSet != nil {
		attrs = append(attrs, "defaultRuleSet")
	}
	if resp.OverrideRuleSet != nil {
		attrs = append(attrs, "overrideRuleSet")
	}
	sort.Strings(attrs)
	return attrs
}

// func parseResponseAttrs(config *SchemaConfigResponse) (int, []string) {
// 	count := 0
// 	var attrs []string
// 	v := reflect.ValueOf(config)
// 	t := reflect.TypeOf(config)

// 	if v.Kind() == reflect.Ptr {
// 		v = v.Elem()
// 		t = t.Elem()
// 	}

// 	for i := 0; i < v.NumField(); i++ {
// 		fieldValue := v.Field(i)
// 		fieldType := t.Field(i)

// 		if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
// 			count++
// 			attrs = append(attrs, fieldType.Name)
// 		}
// 	}
// 	sort.Strings(attrs)
// 	return count, attrs
// }
