// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type EndpointValidator struct {
	basetypes.StringValue
}

func (v EndpointValidator) Description(_ context.Context) string {
	return "String must start with 'http://' or 'https://'"
}

func (v EndpointValidator) MarkdownDescription(_ context.Context) string {
	return "String must start with 'http://' or 'https://'"
}

func (v EndpointValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if !v.isValid(req.ConfigValue.ValueString()) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid URL",
			"The value must start with 'http://' or 'https://'.",
		)
	}
}

func (v EndpointValidator) ValidateAttribute(ctx context.Context, req xattr.ValidateAttributeRequest, resp *xattr.ValidateAttributeResponse) {
	if v.IsNull() || v.IsUnknown() {
		return
	}

	if !v.isValid(v.ValueString()) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid URL",
			"The value must start with 'http://' or 'https://'.",
		)
		return
	}
}

func (v EndpointValidator) ValidateParameter(ctx context.Context, req function.ValidateParameterRequest, resp *function.ValidateParameterResponse) {
	if v.IsNull() || v.IsUnknown() {
		return
	}

	if !v.isValid(v.ValueString()) {
		resp.Error = function.NewArgumentFuncError(
			req.Position,
			"Invalid URL. The value must start with 'http://' or 'https://'.",
		)
		return
	}
}

func (v EndpointValidator) isValid(in string) bool {
	return strings.HasPrefix(in, "http://") || strings.HasPrefix(in, "https://")
}
