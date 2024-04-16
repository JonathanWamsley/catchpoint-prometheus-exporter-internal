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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockServer helps in creating a mock HTTP server for testing purposes.
func mockServer(response string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		fmt.Fprint(w, response)
	}))
}

// TestFetchNodeStatusSuccess tests the successful fetching of node status.
func TestFetchNodeStatusSuccess(t *testing.T) {
	mockResponse := `{
        "data": {
            "nodes": [{"id": 1, "name": "Test Node", "status": {"id": 1, "name": "active"}}]
        },
        "completed": true
    }`

	server := mockServer(mockResponse, http.StatusOK)
	defer server.Close()

	client := NewCatchpointClient("testToken")
	client.BaseURL = server.URL

	nodeStatusResponse, err := client.FetchNodeStatus("1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if nodeStatusResponse == nil {
		t.Errorf("Expected node status response to be populated")
	}
}

// TestFetchNodeStatusError tests error handling when fetching node status.
func TestFetchNodeStatusError(t *testing.T) {
	mockResponse := `{
        "errors": [{"message": "You do not have permission to access this information."}],
        "completed": false
    }`

	server := mockServer(mockResponse, http.StatusForbidden)
	defer server.Close()

	client := NewCatchpointClient("testToken")
	client.BaseURL = server.URL

	_, err := client.FetchNodeStatus("1")
	if err == nil {
		t.Errorf("Expected an error, got none")
	}
}

func TestCatchpointClient_FetchSLAPurgeItems_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"slaItems":[{"id":0,"name":"string","reason":"string","statusType":{"id":0,"name":"string"},"intervalStart":"2024-04-11T00:33:41.497Z","intervalEnd":"2024-04-11T00:33:41.497Z","purgeRuns":{"id":0,"name":"string"},"tests":[{"id":0,"name":"string"}]}]},"completed":true}`))
	}))
	defer server.Close()

	client := NewCatchpointClient("testToken")
	client.BaseURL = server.URL

	_, err := client.FetchSLAPurgeItems()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCatchpointClient_FetchSLAPurgeItems_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"errors":[{"id":"string","message":"string"}],"completed":false}`))
	}))
	defer server.Close()

	client := NewCatchpointClient("testToken")
	client.BaseURL = server.URL

	_, err := client.FetchSLAPurgeItems()
	if err == nil {
		t.Errorf("Expected an error, got none")
	}
}

func TestCatchpointClient_FetchTestErrorsRaw_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
            "data": {
                "responseItems": [{
                    "startTimeUtc": "2024-04-10T23:13:49.2933704Z",
                    "endTimeUtc": "2024-04-11T02:13:49.2933704Z",
                    "summaryItems": [{
                        "values": [5],
                        "dimensions": [{
                            "type": {"id": 41, "name": "ErrorType"},
                            "id": 2,
                            "name": "Connection"
                        }]
                    }]
                }]
            },
            "completed": true
        }`))
	}))
	defer server.Close()

	client := NewCatchpointClient("testToken")
	client.BaseURL = server.URL

	resp, err := client.FetchTestErrorsRaw()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected non-nil response, got nil")
	}

	if !resp.Completed {
		t.Errorf("Expected response to be completed, got %v", resp.Completed)
	}
}

func TestCatchpointClient_FetchTestErrorsRaw_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{
            "errors": [{"id":"string","message":"You do not have permission to access this information."}],
            "completed": false
        }`))
	}))
	defer server.Close()

	client := NewCatchpointClient("testToken")
	client.BaseURL = server.URL

	_, err := client.FetchTestErrorsRaw()
	if err == nil {
		t.Errorf("Expected an error, got none")
	}
}

func TestCatchpointClient_FetchAlerts_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
            "data": {
                "alerts": [{
                    "id": "1",
                    "reportTime": "2024-04-11T05:02:55.710Z",
                    "level": {"id": 1, "name": "Critical"},
                    "test": {"id": 101, "name": "Home Page Performance"},
                    "node": {"id": 10, "name": "New York Node"}
                }],
                "hasMore": false,
                "next": "",
                "previous": ""
            },
            "completed": true
        }`))
	}))
	defer server.Close()

	client := NewCatchpointClient("testToken")
	client.BaseURL = server.URL

	resp, err := client.FetchAlerts()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected non-nil response, got nil")
	}

	if !resp.Completed {
		t.Errorf("Expected response to be completed, got %v", resp.Completed)
	}

}

func TestCatchpointClient_FetchAlerts_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{
            "errors": [{"id":"forbidden","message":"You do not have permission to access this resource."}],
            "completed": false
        }`))
	}))
	defer server.Close()

	client := NewCatchpointClient("testToken")
	client.BaseURL = server.URL

	_, err := client.FetchAlerts()
	if err == nil {
		t.Errorf("Expected an error, got none")
	}
}

func TestCatchpointClient_FetchNodeTestRuns_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
            "data": {
                "node": {
                    "id": 0,
                    "name": "string",
                    "isPaused": true,
                    "instances": [
                        {
                            "id": 0,
                            "hostName": "string",
                            "status": {
                                "id": 0,
                                "name": "string"
                            }
                        }
                    ]
                },
                "testRuns": [
                    {
                        "testId": 0,
                        "testName": "string",
                        "divisionName": "string",
                        "runs": 0,
                        "usagePercentage": 0
                    }
                ],
                "totalTests": 0,
                "hasMore": false
            },
            "completed": true
        }`))
	}))
	defer server.Close()

	client := NewCatchpointClient("testToken")
	client.BaseURL = server.URL

	nodeId := 0
	resp, err := client.FetchNodeTestRuns(nodeId)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected non-nil response, got nil")
	}

	if !resp.Completed {
		t.Errorf("Expected response to be completed, got %v", resp.Completed)
	}

}

func TestCatchpointClient_FetchNodeTestRuns_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{
            "errors": [{"id":"forbidden","message":"You do not have permission to access this resource."}],
            "completed": false
        }`))
	}))
	defer server.Close()

	client := NewCatchpointClient("testToken")
	client.BaseURL = server.URL

	nodeId := 0
	_, err := client.FetchNodeTestRuns(nodeId)
	if err == nil {
		t.Errorf("Expected an error, got none")
	}
}

func TestCatchpointClient_FetchNodeTestRunCount_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"node": {
					"id": 1,
					"name": "Node 1",
					"isPaused": false,
					"runRate": 75,
					"instanceCount": 10,
					"activeInstanceCount": 5
				},
				"allTestRuns": [
					{
						"monitorSetType": {"id": 1, "name": "Browser"},
						"data": [{"reportTime": "2024-04-12T10:16:07.219Z", "value": 25}]
					}
				],
				"uniqueTestRuns": [
					{
						"monitorSetType": {"id": 1, "name": "Browser"},
						"data": [{"reportTime": "2024-04-12T10:16:07.219Z", "value": 15}]
					}
				],
				"hasMore": false
			},
			"completed": true
		}`))
	}))
	defer server.Close()

	client := NewCatchpointClient("testToken")
	client.BaseURL = server.URL

	nodeId := 1
	resp, err := client.FetchNodeTestRunCount(nodeId)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected non-nil response, got nil")
	}

	if !resp.Completed {
		t.Errorf("Expected response to be completed, got %v", resp.Completed)
	}

	if len(resp.Data.AllTestRuns) != 1 {
		t.Errorf("Expected 1 all test run entry, got %d", len(resp.Data.AllTestRuns))
	}
}

func TestCatchpointClient_FetchNodeTestRunCount_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{
            "errors": [{"id":"forbidden","message":"You do not have permission to access this resource."}],
            "completed": false
        }`))
	}))
	defer server.Close()

	client := NewCatchpointClient("testToken")
	client.BaseURL = server.URL

	nodeId := 1
	_, err := client.FetchNodeTestRunCount(nodeId)
	if err == nil {
		t.Errorf("Expected an error, got none")
	}
}
