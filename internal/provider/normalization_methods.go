// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func SetNormalization(client *Client, subject_name string, payload NormalizeRequest) (*NormalizeResponse, error) {
	rb, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/config/%s", client.HostURL, subject_name), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(client.Auth.Username, client.Auth.Password)

	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update subject configuration. Response code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response NormalizeResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func GetSchemaConfig(client *Client, subject_name string) (*SchemaConfigResponse, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/config/%s", client.HostURL, subject_name), nil)
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
		return nil, fmt.Errorf("failed to update subject configuration. Response code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response SchemaConfigResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func DeleteSubjectConfig(client *Client, subject_name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/config/%s", client.HostURL, subject_name), nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(client.Auth.Username, client.Auth.Password)

	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response code received on deleting schema config '%s' request. Expected %d Received %d", subject_name, http.StatusNoContent, res.StatusCode)
	}

	return nil
}
