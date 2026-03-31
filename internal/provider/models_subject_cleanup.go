// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type schemaVersions struct {
	softDeleted      []int
	active           []int
	all              []int
	client           *Client
	deleteCandidates []int
	schemasToKeep    int
}

func (r *schemaVersions) get(model subjectCleanupResourceModel) error {
	var err error
	r.all, r.active, r.softDeleted, err = GetSchemaVersions(model, r.client)
	return err
}

func (r *schemaVersions) countSchemasToKeep(model subjectCleanupResourceModel) {
	switch model.CleanupMethod {
	case types.StringValue("KEEP_ACTIVE_ONLY"):
		r.schemasToKeep = len(r.active)
	case types.StringValue("KEEP_LATEST_ONLY"):
		r.schemasToKeep = 1
	case types.StringValue("MAX_STORED_SCHEMAS"):
		r.schemasToKeep = int(model.SchemasToKeep.ValueInt64())
	}
}

func (r *schemaVersions) cleanDeleteCandidates(ctx context.Context, model subjectCleanupResourceModel) error {
	soft := true
	err := DeleteSchemaVersions(&r.deleteCandidates, ctx, r.client, model, soft)
	return err
}

func (r *schemaVersions) calculateDeleteCandidates() {
	if len(r.all) > r.schemasToKeep {
		r.deleteCandidates = r.all[:len(r.all)-r.schemasToKeep]
	}
}
