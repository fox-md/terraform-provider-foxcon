// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) DeleteUser(userID string) error {

	if userID == "" {
		return fmt.Errorf("UserID is empty")
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/iam/v2/users/%s", c.HostURL, userID), nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.Auth.Username, c.Auth.Password)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected response code received on User DELETE request. Expected %d Received %d", http.StatusNoContent, res.StatusCode)
	}

	return nil
}

func (c *Client) ReadUser(userId string) (*User, error) {

	user := User{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/iam/v2/users/%s", c.HostURL, userId), nil)

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

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response code received on User DELETE request. Expected %d Received %d", http.StatusOK, res.StatusCode)
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
