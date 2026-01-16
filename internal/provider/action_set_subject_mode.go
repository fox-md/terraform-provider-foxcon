// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ action.Action              = (*subjectModeAction)(nil)
	_ action.ActionWithConfigure = &subjectModeAction{}
)

func SetSubjectModeAction() action.Action {
	return &subjectModeAction{}
}

type subjectModeAction struct {
	client *Client
}

func (r *subjectModeAction) Configure(_ context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
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

func (a *subjectModeAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_set_subject_mode"
}

func (a *subjectModeAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a Subject Mode action that sets Subject Mode on a Schema Registry cluster on Confluent Cloud.",
		Attributes: map[string]schema.Attribute{
			"rest_endpoint": schema.StringAttribute{
				Optional:    true,
				Description: restEndpointDescription,
				Validators: []validator.String{
					EndpointValidator{},
				},
			},
			"subject_name": schema.StringAttribute{
				Required:    true,
				Description: subjectNameDescription,
			},
			"mode": schema.StringAttribute{
				Required:    true,
				Description: "The mode of the specified subject. Accepted values are: `READWRITE`, `READONLY`, `READONLY_OVERRIDE` and `IMPORT`.",
				Validators: []validator.String{
					stringvalidator.OneOf("READWRITE", "READONLY", "READONLY_OVERRIDE", "IMPORT"),
				},
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
						Description: schemaRegistrySecretDescription + " Terraform actions do NOT support sensitive attributes. Please keep that in mind.",
					},
				},
			},
		},
	}
}

func (a *subjectModeAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var config subjectModeResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	creds := schemaRegistryCredentials{
		RestEndpoint: config.RestEndpoint,
		Credentials:  config.Credentials,
	}

	schemaAPIClient, err := schemaRegistryClientFactory(a.client, &creds)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating http client",
			"Could not create http client. Unexpected error: "+err.Error(),
		)
		return
	}

	var subjectModePayload = SubjectModeRequest{
		Mode: *config.Mode.ValueStringPointer(),
	}

	subjectMode, err := SetSubjectMode(schemaAPIClient, config.SubjectName.ValueString(), subjectModePayload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting normalization",
			"Could not set normalization unexpected error: "+err.Error(),
		)
		return
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("\n\nSubject '%s' has been set to the '%s' mode", config.SubjectName.ValueString(), subjectMode.Mode),
	})
}

type subjectModeResourceModel struct {
	RestEndpoint types.String      `tfsdk:"rest_endpoint"`
	SubjectName  types.String      `tfsdk:"subject_name"`
	Mode         types.String      `tfsdk:"mode"`
	Credentials  *credentialsModel `tfsdk:"credentials"`
}
