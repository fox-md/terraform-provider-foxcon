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

func (c *Client) GetUserInvitationById(invitationId string) (*Invitation, error) {

	invitation := Invitation{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/iam/v2/invitations/%s", c.HostURL, invitationId), nil)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Auth.Username, c.Auth.Password)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	err = json.Unmarshal(body, &invitation)
	if err != nil {
		return nil, err
	}

	return &invitation, nil
}

func (c *Client) CreateInvitation(payload InvitationItem) (*Invitation, error) {
	rb, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/iam/v2/invitations", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Auth.Username, c.Auth.Password)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var invitation Invitation

	err = json.Unmarshal(body, &invitation)
	if err != nil {
		return nil, err
	}

	return &invitation, nil
}

func (c *Client) GetUserInvitationByParameter(search_type, search_parameter string) (*Invitation, error) {

	invitationList := InvitationList{}
	invitation := Invitation{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/iam/v2/invitations?%s=%s", c.HostURL, search_type, search_parameter), nil)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Auth.Username, c.Auth.Password)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get list of invitations. Response code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &invitationList)
	if err != nil {
		return nil, err
	}

	if len(invitationList.Data) == 0 {
		return nil, nil
	}

	invitation = invitationList.Data[0]

	return &invitation, nil
}
