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

func SetSubjectMode(client *Client, subject_name string, payload SubjectModeRequest) (*SubjectModeResponse, error) {
	rb, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/mode/%s", client.HostURL, subject_name), strings.NewReader(string(rb)))
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

	var response SubjectModeResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func GetSubjectMode(client *Client, subject_name string) (*SubjectModeResponse, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/mode/%s", client.HostURL, subject_name), nil)
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
		return nil, fmt.Errorf("failed to get subject configuration. Response code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response SubjectModeResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
