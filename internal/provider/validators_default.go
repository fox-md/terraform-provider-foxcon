// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var restEndpointValidators = []validator.String{
	EndpointValidator{},
	stringvalidator.AlsoRequires(
		path.MatchRoot("credentials").AtName("key"),
	),
	stringvalidator.AlsoRequires(
		path.MatchRoot("credentials").AtName("secret"),
	),
}

var credentialsKeyValidators = []validator.String{
	stringvalidator.LengthAtLeast(1),
	stringvalidator.AlsoRequires(
		path.MatchRoot("credentials").AtName("secret"),
	),
	stringvalidator.AlsoRequires(
		path.MatchRoot("rest_endpoint"),
	),
}

var credentialsSecretValidators = []validator.String{
	stringvalidator.LengthAtLeast(1),
	stringvalidator.AlsoRequires(
		path.MatchRoot("credentials").AtName("key"),
	),
	stringvalidator.AlsoRequires(
		path.MatchRoot("rest_endpoint"),
	),
}

var subjectNameValidators = []validator.String{
	stringvalidator.LengthAtLeast(1),
}
