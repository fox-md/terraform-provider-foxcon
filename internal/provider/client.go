// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"crypto/tls"
	"net/http"
	"strings"
	"time"
)

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Auth       AuthStruct
}

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewClient -
func NewClient(HostURL, username, password *string) (*Client, error) {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	if strings.HasPrefix(*HostURL, "https") {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: false,
		}
	}

	c := Client{
		HTTPClient: &http.Client{
			Timeout:   10 * time.Second,
			Transport: transport,
		},
		HostURL: *HostURL,
	}

	// If username or password not provided, return empty client
	if username == nil || password == nil {
		return &c, nil
	}

	c.Auth = AuthStruct{
		Username: *username,
		Password: *password,
	}

	return &c, nil
}
