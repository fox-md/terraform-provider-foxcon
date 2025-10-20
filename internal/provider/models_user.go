// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "time"

type User struct {
	APIVersion string       `json:"api_version"`
	Kind       string       `json:"kind"`
	ID         string       `json:"id"`
	Metadata   UserMetadata `json:"metadata"`
	Email      string       `json:"email"`
	FullName   string       `json:"full_name"`
	AuthType   string       `json:"auth_type"`
}

type UserMetadata struct {
	Self         string    `json:"self"`
	ResourceName string    `json:"resource_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}
