// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

// NewuserResource is a helper function to simplify the provider implementation.
func NewUserResource() resource.Resource {
	return &userResource{}
}

// userResource is the resource implementation.
type userResource struct {
	client *Client
}

// Metadata returns the resource type name.
func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_read_confluent_user"
	//resp.TypeName = "confluent_user"
}

// Schema defines the schema for the resource.
func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"invitation_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("invitation_id"),
						path.MatchRoot("user_email"),
						path.MatchRoot("user_id"),
					),
				},
				Description: "Confluent invitation id.",
			},
			"user_email": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("invitation_id"),
						path.MatchRoot("user_email"),
						path.MatchRoot("user_id"),
					),
				},
				Description: "User email address.",
			},
			"user_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("invitation_id"),
						path.MatchRoot("user_email"),
						path.MatchRoot("user_id"),
					),
				},
				Description: "Confluent user id.",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp of the last apply execution.",
			},
		},
		MarkdownDescription: "Gets user id, user email and invitation id by setting one of this parameters. On deletion, resource also removes user from Confluent Cloud.",
	}
}

// Create creates the resource and sets the initial Terraform state.
// Create a new resource.
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan userResourceModel
	var searchType string
	var searchValue string
	var invitation *Invitation
	var err error

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.InvitationId.IsNull() && !plan.InvitationId.IsUnknown() {
		searchType = "invitation_id"
		searchValue = plan.InvitationId.ValueString()
		invitation, err = r.client.GetUserInvitationById(plan.InvitationId.ValueString())
	}

	if !plan.UserEmail.IsNull() && !plan.UserEmail.IsUnknown() {
		searchType = "user_email"
		searchValue = plan.UserEmail.ValueString()
		invitation, err = r.client.GetUserInvitationByParameter("email", plan.UserEmail.ValueString())
	}

	if !plan.UserId.IsNull() && !plan.UserId.IsUnknown() {
		searchType = "user_id"
		searchValue = plan.UserId.ValueString()
		invitation, err = r.client.GetUserInvitationByParameter("user", plan.UserId.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading invitation",
			fmt.Sprintf("Could not read invitation for %s by %s type. Unexpected error: %s", searchValue, searchType, err.Error()),
		)
		return
	}

	if invitation == nil {
		resp.Diagnostics.AddError(
			"Invitation not found",
			fmt.Sprintf("Invitation not found for %s by %s type.", searchValue, searchType),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.UserEmail = types.StringValue(invitation.Email)
	plan.InvitationId = types.StringValue(invitation.ID)
	plan.UserId = types.StringValue(invitation.User.ID)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	// // Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
// Read resource information.
func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state userResourceModel
	var searchType string
	var searchValue string
	var invitation *Invitation
	var err error

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.InvitationId.IsNull() {
		searchType = "invitation_id"
		searchValue = state.InvitationId.ValueString()
		invitation, err = r.client.GetUserInvitationById(state.InvitationId.ValueString())
	}

	if !state.UserEmail.IsNull() {
		searchType = "user_email"
		searchValue = state.InvitationId.ValueString()
		invitation, err = r.client.GetUserInvitationByParameter("email", state.UserEmail.ValueString())
	}

	if !state.UserId.IsNull() {
		searchType = "user_id"
		searchValue = state.InvitationId.ValueString()
		invitation, err = r.client.GetUserInvitationByParameter("user", state.UserId.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Confluent invitation",
			fmt.Sprintf("Could not read invitation for %s by %s. Unexpected error: %s", searchValue, searchType, err.Error()),
		)
		return
	}

	if invitation == nil {
		tflog.Debug(ctx, fmt.Sprintf("User with %s of type %s does not exist in Confluent. Removing resource from state file.", searchValue, searchType))
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddWarning(
			"User does not exist in Confluent",
			fmt.Sprintf("User with %s of type %s does not exist in Confluent. Removing resource from state file.", searchValue, searchType),
		)
		return
	}

	// Overwrite items with refreshed state
	state.UserEmail = types.StringValue(invitation.Email)
	state.InvitationId = types.StringValue(invitation.ID)
	state.UserId = types.StringValue(invitation.User.ID)
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan userResourceModel
	var invitation *Invitation
	var err error

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.InvitationId.IsNull() {
		invitation, err = r.client.GetUserInvitationById(plan.InvitationId.ValueString())
	}

	if !plan.UserEmail.IsNull() {
		invitation, err = r.client.GetUserInvitationByParameter("email", plan.UserEmail.ValueString())
	}

	if !plan.UserId.IsNull() {
		invitation, err = r.client.GetUserInvitationByParameter("user", plan.UserId.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading invitation",
			"Could not read invitation unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.UserEmail = types.StringValue(invitation.Email)
	plan.InvitationId = types.StringValue(invitation.ID)
	plan.UserId = types.StringValue(invitation.User.ID)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing user
	err := r.client.DeleteUser(state.UserId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Confluent User",
			"Could not delete user, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.AddError(
		"Import not implemented",
		"Import for this resource is not available since the resource itself performs read operation on create",
	)

	// resource.ImportStatePassthroughID(ctx, path.Root("user_email"), req, resp)

	// invitationIdPattern := regexp.MustCompile(`^i-[a-z0-9]{6,}$`)
	// userIdPattern := regexp.MustCompile(`^u-[a-z0-9]{6,}$`)
	// emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	// if invitationIdPattern.MatchString(req.ID) {
	// 	tflog.Debug(ctx, fmt.Sprintf("Importing user by invitation id:%s", req.ID))
	// 	resource.ImportStatePassthroughID(ctx, path.Root("invitation_id"), req, resp)
	// }

	// if userIdPattern.MatchString(req.ID) {
	// 	tflog.Debug(ctx, fmt.Sprintf("Importing user by user id:%s", req.ID))
	// 	resource.ImportStatePassthroughID(ctx, path.Root("user_id"), req, resp)
	// }

	// if emailPattern.MatchString(req.ID) {
	// 	tflog.Debug(ctx, fmt.Sprintf("Importing user by user email id:%s", req.ID))
	// 	resource.ImportStatePassthroughID(ctx, path.Root("user_email"), req, resp)
	// }
}

type userResourceModel struct {
	InvitationId types.String `tfsdk:"invitation_id"`
	UserEmail    types.String `tfsdk:"user_email"`
	UserId       types.String `tfsdk:"user_id"`
	LastUpdated  types.String `tfsdk:"last_updated"`
}
