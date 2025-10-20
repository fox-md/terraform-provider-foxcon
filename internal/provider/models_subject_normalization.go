// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

type NormalizeRequest struct {
	Normalize *bool `json:"normalize"`
}

type NormalizeResponse struct {
	Normalize *bool `json:"normalize,omitempty"`
}
