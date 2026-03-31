// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"sort"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func ListSubjectVersions(client *Client, subject_name string, deleted bool) ([]int, error) {

	if subject_name == "" {
		return nil, fmt.Errorf("subject name not configured")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/subjects/%s/versions?deleted=%t", client.HostURL, subject_name, deleted), nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(client.Auth.Username, client.Auth.Password)

	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Subject config does not exist
	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get subject versions. Response code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response []int

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func DeleteSchemaVersion(client *Client, subject_name string, version int, permanent bool) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/subjects/%s/versions/%d?permanent=%t",
		client.HostURL, subject_name, version, permanent), nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(client.Auth.Username, client.Auth.Password)

	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
		return fmt.Errorf("unexpected response code received on deleting schema version '%d' for subject '%s'. Expected [%d, %d] Received %d", version, subject_name, http.StatusOK, http.StatusNotFound, res.StatusCode)
	}

	return nil
}

func GetSchemaVersions(model subjectCleanupResourceModel, client *Client) ([]int, []int, []int, error) {
	all, err := ListSubjectVersions(client, model.SubjectName.ValueString(), true)
	if err != nil {
		return nil, nil, nil, err
	}

	active, err := ListSubjectVersions(client, model.SubjectName.ValueString(), false)
	if err != nil {
		return nil, nil, nil, err
	}

	var softDeleted []int

	for _, v := range all {
		if !slices.Contains(active, v) {
			softDeleted = append(softDeleted, v)
		}
	}

	sort.Ints(all)
	sort.Ints(active)
	return all, active, softDeleted, nil
}

func DeleteSchemaVersions(versions *[]int, ctx context.Context, client *Client, model subjectCleanupResourceModel, soft bool) error {
	for _, v := range *versions {
		tflog.Debug(ctx, fmt.Sprintf("Deleting %s version %v", model.SubjectName.ValueString(), v))
		if soft {
			err := DeleteSchemaVersion(client, model.SubjectName.ValueString(), v, false)
			if err != nil {
				return fmt.Errorf("could not soft delete schema version. Unexpected error: %s", err.Error())
			}
		}

		err := DeleteSchemaVersion(client, model.SubjectName.ValueString(), v, true)
		if err != nil {
			return fmt.Errorf("could not hard delete schema version. Unexpected error: %s", err.Error())
		}

	}
	return nil
}

func SubjectCleanup(ctx context.Context, client *Client, model *subjectCleanupResourceModel) (diag.Diagnostics, error) {
	var subjectVersions schemaVersions
	var diags diag.Diagnostics
	var latestVersion int

	creds := schemaRegistryCredentials{
		RestEndpoint: model.RestEndpoint,
		Credentials:  model.Credentials,
	}

	schemaAPIClient, err := schemaRegistryClientFactory(client, &creds)
	if err != nil {
		return diags, err
	}

	subjectVersions.client = schemaAPIClient
	err = subjectVersions.get(*model)
	if err != nil {
		return diags, err
	}

	subjectVersions.countSchemasToKeep(*model)
	subjectVersions.calculateDeleteCandidates()
	err = subjectVersions.cleanDeleteCandidates(ctx, *model)
	if err != nil {
		return diags, err
	}

	if len(subjectVersions.all) > 0 {
		latestVersion = subjectVersions.all[len(subjectVersions.all)-1]
	} else {
		latestVersion = 0
	}

	var lastDeleted []attr.Value
	for _, id := range subjectVersions.deleteCandidates {
		lastDeleted = append(lastDeleted, types.Int32Value(int32(id)))
	}

	model.SchemasToKeep = types.Int64Value(int64(subjectVersions.schemasToKeep))
	model.LastSchemaVersion = types.Int32Value(int32(latestVersion))
	model.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))
	model.CleanupNeeded = types.BoolValue(false)
	model.LastDeleted, diags = types.ListValue(types.Int32Type, lastDeleted)

	return diags, nil
}

func ReadSubjectVersions(ctx context.Context, client *Client, model subjectCleanupResourceModel) (schemaVersions, error) {

	var subjectVersions schemaVersions

	creds := schemaRegistryCredentials{
		RestEndpoint: model.RestEndpoint,
		Credentials:  model.Credentials,
	}

	schemaAPIClient, err := schemaRegistryClientFactory(client, &creds)
	if err != nil {
		return subjectVersions, err
	}

	subjectVersions.client = schemaAPIClient
	err = subjectVersions.get(model)
	if err != nil {
		return subjectVersions, err
	}

	subjectVersions.countSchemasToKeep(model)
	subjectVersions.calculateDeleteCandidates()

	return subjectVersions, nil
}
