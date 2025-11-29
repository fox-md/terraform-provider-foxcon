// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

type SubjectModeRequest struct {
	Mode string `json:"mode"`
}

type SubjectModeResponse struct {
	Mode string `json:"mode"`
}
