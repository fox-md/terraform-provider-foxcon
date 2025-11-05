// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type EndpointValidator struct{}

func (v EndpointValidator) Description(_ context.Context) string {
	return "String must start with 'http'"
}

func (v EndpointValidator) MarkdownDescription(_ context.Context) string {
	return "String must start with `http`"
}

func (v EndpointValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if !strings.HasPrefix(req.ConfigValue.ValueString(), "http") {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid URL",
			"The value must start with 'http'.",
		)
	}
}
