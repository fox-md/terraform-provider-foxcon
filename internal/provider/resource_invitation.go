// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &invitationResource{}
	_ resource.ResourceWithConfigure   = &invitationResource{}
	_ resource.ResourceWithImportState = &invitationResource{}
)

// NewinvitationResource is a helper function to simplify the provider implementation.
func NewInvitationResource() resource.Resource {
	return &invitationResource{}
}

// invitationResource is the resource implementation.
type invitationResource struct {
	client *Client
}

// Metadata returns the resource type name.
func (r *invitationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_confluent_invitation"
}

// Schema defines the schema for the resource.
func (r *invitationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				Required:    true,
				Description: "User's/invitee's email address.",
			},
			"auth_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("AUTH_TYPE_SSO", "AUTH_TYPE_LOCAL"),
				},
				Default:     stringdefault.StaticString("AUTH_TYPE_SSO"),
				Description: "User's/invitee's authentication type. Must be set to `AUTH_TYPE_SSO` or `AUTH_TYPE_LOCAL`. Defauts to `AUTH_TYPE_SSO`.",
			},
			"invitation_id": schema.StringAttribute{
				Computed:    true,
				Description: "Confluent invitation id.",
			},
			"user_id": schema.StringAttribute{
				Computed:    true,
				Description: "Confluent user id.",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp of the last apply execution.",
			},
		},
		MarkdownDescription: "Provides an invitation resource that enables creating, reading, and deleting invitation on Confluent Cloud. On deleting invitation also deletes user from Confluent Cloud.",
	}
}

// Create creates the resource and sets the initial Terraform state.
// Create a new resource.
func (r *invitationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan invitationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var invitationPayload = InvitationItem{
		Email:    plan.Email.ValueString(),
		AuthType: plan.AuthType.ValueString(),
	}

	// Create Invitation
	invitation, err := r.client.CreateInvitation(invitationPayload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating invitation",
			"Could not create invitation unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.UserId = types.StringValue(invitation.User.ID)
	plan.InvitationId = types.StringValue(invitation.ID)
	plan.AuthType = types.StringValue(invitation.AuthType)
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
func (r *invitationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state invitationResourceModel
	var invitation *Invitation
	var err error

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.InvitationId.IsNull() {
		invitation, err = r.client.GetUserInvitationById(state.InvitationId.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Confluent invitation by invitation id",
			"Could not read Confluent invitation info by invitation id "+state.InvitationId.ValueString()+": "+err.Error(),
		)
		return
	}

	if !state.Email.IsNull() {
		invitation, err = r.client.GetUserInvitationByParameter("email", state.Email.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Confluent invitation by email",
			"Could not read Confluent invitation by email "+state.Email.ValueString()+": "+err.Error(),
		)
		return
	}

	if invitation == nil {
		tflog.Debug(ctx, fmt.Sprintf("Invitation with ID %s does not exist in Confluent. Removing resource from state file.", state.InvitationId.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	// Overwrite items with refreshed state
	state.Email = types.StringValue(invitation.Email)
	state.UserId = types.StringValue(invitation.User.ID)
	state.InvitationId = types.StringValue(invitation.ID)
	state.AuthType = types.StringValue(invitation.AuthType)
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *invitationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan invitationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get User Invitation
	invitation, err := r.client.GetUserInvitationById(plan.InvitationId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading invitation",
			"Could not read invitation unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.UserId = types.StringValue(invitation.User.ID)
	plan.InvitationId = types.StringValue(invitation.ID)
	plan.AuthType = types.StringValue(invitation.AuthType)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *invitationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state invitationResourceModel
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
func (r *invitationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	clients, ok := req.ProviderData.(*providerClients)

	//client, ok := req.ProviderData.(*Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	ValidateCloudApiClient(clients.CloudApiClient, resp)

	if resp.Diagnostics.HasError() {
		return
	}

	r.client = clients.CloudApiClient
}

func (r *invitationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import email and save to email attribute
	// resource.ImportStatePassthroughID(ctx, path.Root("invitation_id"), req, resp)

	invitation_id_pattern := regexp.MustCompile(`^i-[a-z0-9]{6,}$`)
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	if invitation_id_pattern.MatchString(req.ID) {
		tflog.Debug(ctx, fmt.Sprintf("Importing invitation by invitation id:%s", req.ID))
		resource.ImportStatePassthroughID(ctx, path.Root("invitation_id"), req, resp)
	}

	if emailPattern.MatchString(req.ID) {
		tflog.Debug(ctx, fmt.Sprintf("Importing invitation by user email:%s", req.ID))
		resource.ImportStatePassthroughID(ctx, path.Root("email"), req, resp)
	}
}

type invitationResourceModel struct {
	Email        types.String `tfsdk:"email"`
	AuthType     types.String `tfsdk:"auth_type"`
	InvitationId types.String `tfsdk:"invitation_id"`
	UserId       types.String `tfsdk:"user_id"`
	LastUpdated  types.String `tfsdk:"last_updated"`
}
