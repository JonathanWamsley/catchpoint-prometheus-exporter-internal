// Copyright 2024 Grafana Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type HttpClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// Ensure CatchpointClient's HttpClient field conforms to this interface.
var _ HttpClientInterface = &http.Client{}

type CatchpointClientInterface interface {
	FetchNodeStatus(nodeId string) (*NodeStatusResponse, error)
	FetchSLAPurgeItems() (*SLAPurgeItemsResponse, error)
	FetchTestErrorsRaw() (*TestErrorsRawResponse, error)
	FetchAlerts() (*AlertsResponse, error)
	FetchNodeTestRuns(nodeId int) (*NodeTestRunResponse, error)
	FetchNodeRunRate(nodeId int) (*NodeRunRateResponse, error)
	FetchNodeTestRunCount(nodeId int) (*TestRunCountResponse, error)
}

type CatchpointClient struct {
	HttpClient  HttpClientInterface
	BaseURL     string
	BearerToken string
}

func NewCatchpointClient(bearerToken string) *CatchpointClient {
	return &CatchpointClient{
		HttpClient:  &http.Client{Timeout: 10 * time.Second},
		BaseURL:     "https://io.catchpoint.com/api/v2",
		BearerToken: bearerToken,
	}
}

// Ensure CatchpointClient implements CatchpointClientInterface.
var _ CatchpointClientInterface = (*CatchpointClient)(nil)

func parseAPIErrorResponse(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var apiError struct { // Struct to parse the errors
			Errors []struct {
				Message string `json:"message"`
			} `json:"errors"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&apiError); err == nil {
			var errMsgs []string
			for _, e := range apiError.Errors {
				errMsgs = append(errMsgs, e.Message)
			}
			return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, strings.Join(errMsgs, ", "))
		}
		// Fallback if the error parsing fails
		return fmt.Errorf("API request failed with status %d and could not parse error body", resp.StatusCode)
	}
	return nil
}

// doRequest prepares and executes an HTTP request and unmarshals the response.
func (c *CatchpointClient) doRequest(method, url string, headers map[string]string, result interface{}) error {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := parseAPIErrorResponse(resp); err != nil {
		return err
	}

	// Decode the body into the provided result interface if no error
	return json.NewDecoder(resp.Body).Decode(result)
}

func (c *CatchpointClient) FetchNodeStatus(nodeId string) (*NodeStatusResponse, error) {
	url := fmt.Sprintf("%s/nodes/status/%s", c.BaseURL, nodeId)
	headers := map[string]string{
		"Authorization": "Bearer " + c.BearerToken,
		"Accept":        "application/json",
	}
	var response NodeStatusResponse
	if err := c.doRequest("GET", url, headers, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *CatchpointClient) FetchSLAPurgeItems() (*SLAPurgeItemsResponse, error) {
	url := fmt.Sprintf("%s/slapurgeitems", c.BaseURL)
	headers := map[string]string{
		"Authorization": "Bearer " + c.BearerToken,
		"Accept":        "application/json",
	}
	var response SLAPurgeItemsResponse
	if err := c.doRequest("GET", url, headers, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *CatchpointClient) FetchTestErrorsRaw() (*TestErrorsRawResponse, error) {
	url := fmt.Sprintf("%s/tests/errors/raw", c.BaseURL)
	headers := map[string]string{
		"Authorization": "Bearer " + c.BearerToken,
		"Accept":        "application/json",
	}
	var response TestErrorsRawResponse
	if err := c.doRequest("GET", url, headers, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *CatchpointClient) FetchAlerts() (*AlertsResponse, error) {
	url := fmt.Sprintf("%s/tests/alerts", c.BaseURL)
	headers := map[string]string{
		"Authorization": "Bearer " + c.BearerToken,
		"Accept":        "application/json",
	}
	var response AlertsResponse
	if err := c.doRequest("GET", url, headers, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *CatchpointClient) FetchNodeTestRuns(nodeId int) (*NodeTestRunResponse, error) {
	url := fmt.Sprintf("%s/nodes/testrun/%d", c.BaseURL, nodeId)
	headers := map[string]string{
		"Authorization": "Bearer " + c.BearerToken,
		"Accept":        "application/json",
	}
	var response NodeTestRunResponse
	if err := c.doRequest("GET", url, headers, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *CatchpointClient) FetchNodeRunRate(nodeId int) (*NodeRunRateResponse, error) {
	url := fmt.Sprintf("%s/nodes/runrate/%d", c.BaseURL, nodeId)
	headers := map[string]string{
		"Authorization": "Bearer " + c.BearerToken,
		"Accept":        "application/json",
	}
	var response NodeRunRateResponse
	if err := c.doRequest("GET", url, headers, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *CatchpointClient) FetchNodeTestRunCount(nodeId int) (*TestRunCountResponse, error) {
	url := fmt.Sprintf("%s/nodes/testruncount/%d", c.BaseURL, nodeId)
	headers := map[string]string{
		"Authorization": "Bearer " + c.BearerToken,
		"Accept":        "application/json",
	}
	var response TestRunCountResponse
	if err := c.doRequest("GET", url, headers, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
