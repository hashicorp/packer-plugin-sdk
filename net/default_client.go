// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package net

import (
	"net/http"
)

func HttpClientWithEnvironmentProxy() *http.Client {
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}
	return httpClient
}
