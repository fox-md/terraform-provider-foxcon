// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type SchemasNumberValidator struct {
	basetypes.Int64Value
}

func (v SchemasNumberValidator) Description(_ context.Context) string {
	return "Number of schemas must be more than 0 when cleanup_method is set to 'MAX_STORED_SCHEMAS'"
}

func (v SchemasNumberValidator) MarkdownDescription(_ context.Context) string {
	return "Number of schemas must be more than 0 when cleanup_method is set to 'MAX_STORED_SCHEMAS'"
}

func (v SchemasNumberValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	var cleanupMethod types.String

	diags := req.Config.GetAttribute(ctx, path.Root("cleanup_method"), &cleanupMethod)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if (req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown()) && cleanupMethod.ValueString() == "MAX_STORED_SCHEMAS" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid value",
			"Number of schemas must be more than 0 when cleanup_method is set to 'MAX_STORED_SCHEMAS'",
		)
	}

	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if !v.isValid(int(req.ConfigValue.ValueInt64())) && cleanupMethod.ValueString() == "MAX_STORED_SCHEMAS" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid value",
			"Number of schemas must be more than 0 when cleanup_method is set to 'MAX_STORED_SCHEMAS'",
		)
	}
}

func (v SchemasNumberValidator) isValid(in int) bool {
	return in > 0
}
