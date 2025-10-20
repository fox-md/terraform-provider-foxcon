// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

type InvitationList struct {
	APIVersion string       `json:"api_version"`
	Kind       string       `json:"kind"`
	Metadata   ListMeta     `json:"metadata"`
	Data       []Invitation `json:"data"`
}

type ListMeta struct {
	First     string `json:"first"`
	Last      string `json:"last"`
	Prev      string `json:"prev"`
	Next      string `json:"next"`
	TotalSize int    `json:"total_size"`
}

type Invitation struct {
	APIVersion string             `json:"api_version"`
	Kind       string             `json:"kind"`
	ID         string             `json:"id"`
	Metadata   InvitationMetadata `json:"metadata"`
	Email      string             `json:"email"`
	AuthType   string             `json:"auth_type"`
	Status     string             `json:"status"`
	AcceptedAt string             `json:"accepted_at"`
	ExpiresAt  string             `json:"expires_at"`
	User       UserEntity         `json:"user"`
	Creator    UserEntity         `json:"creator"`
}

type InvitationMetadata struct {
	Self         string `json:"self"`
	ResourceName string `json:"resource_name"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	DeletedAt    string `json:"deleted_at"`
}

type UserEntity struct {
	ID           string `json:"id"`
	Related      string `json:"related"`
	ResourceName string `json:"resource_name"`
}

type InvitationItem struct {
	Email    string `json:"email"`
	AuthType string `json:"auth_type"`
}
