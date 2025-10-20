// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "context"

type schemaVersions struct {
	softDeleted *[]int
	active      *[]int
	all         *[]int
	client      *Client
}

func (r *schemaVersions) get(model subjectCleanupResourceModel) error {
	var err error
	r.all, r.active, r.softDeleted, err = getSchemaVersions(model, r.client)
	return err
}

func (r *schemaVersions) cleanSoftDeleted(ctx context.Context, model subjectCleanupResourceModel) error {
	soft := false
	err := deleteSchemaVersions(r.softDeleted, ctx, r.client, model, soft)
	return err
}

func (r *schemaVersions) cleanActiveNoneLatest(ctx context.Context, model subjectCleanupResourceModel) error {
	soft := true
	activeNoneLatest := (*r.active)[:len(*r.active)-1]
	err := deleteSchemaVersions(&activeNoneLatest, ctx, r.client, model, soft)
	return err
}
