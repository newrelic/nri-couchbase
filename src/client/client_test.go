package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/nri-couchbase/src/arguments"
	"github.com/stretchr/testify/assert"
)

func Test_CreateClient(t *testing.T) {
	args := setupTestArgs()

	client, _ := CreateClient(args, "")

	assert.Equal(t, "http://testhostname:8091", client.BaseURL)
	assert.Equal(t, "http://testhostname:8093", client.BaseQueryURL)
	assert.Equal(t, "testuser", client.Username)
	assert.Equal(t, "testpass", client.Password)
}

func Test_SSL(t *testing.T) {
	args := setupTestArgs()
	args.UseSSL = true

	client, _ := CreateClient(args, "")

	assert.Equal(t, "https://testhostname:8091", client.BaseURL)
}

func Test_HostnameOverride(t *testing.T) {
	args := setupTestArgs()
	client, _ := CreateClient(args, "inventory-host")

	assert.Equal(t, "http://inventory-host:8091", client.BaseURL)
}

func Test_Request(t *testing.T) {
	// generate a test server so we can capture and inspect the request
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		username, password, ok := req.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, username, "testUser")
		assert.Equal(t, password, "testPass")
		res.Write([]byte("{\"ok\":true}"))
	}))
	defer func() { testServer.Close() }()

	testResult := struct {
		OK *bool `json:"ok"`
	}{}

	testCases := []struct {
		Client   *HTTPClient
		Endpoint string
	}{
		{
			// request to standard endpoint which uses the baseURL
			// requests to "bad-url" will fail
			&HTTPClient{
				Client:       testServer.Client(),
				Username:     "testUser",
				Password:     "testPass",
				BaseURL:      testServer.URL,
				BaseQueryURL: "bad-url",
			},
			"/some/endpoint",
		},
		{
			// request to query endpoint which uses the baseQueryURL
			// requests to "bad-url" will fail
			&HTTPClient{
				Client:       testServer.Client(),
				Username:     "testUser",
				Password:     "testPass",
				BaseURL:      "bad-url",
				BaseQueryURL: testServer.URL,
			},
			"/admin/endpoint",
		},
	}

	for _, tc := range testCases {
		err := tc.Client.Request(tc.Endpoint, &testResult)
		assert.NoError(t, err)
		assert.Equal(t, true, *testResult.OK)
	}
}

func Test_checkStatusCode(t *testing.T) {
	response := &http.Response{
		StatusCode: 404,
	}

	err := checkStatusCode(response)
	assert.Error(t, err)
}

func setupTestArgs() *arguments.ArgumentList {
	return &arguments.ArgumentList{
		DefaultArgumentList: args.DefaultArgumentList{
			Verbose:   false,
			Pretty:    false,
			Metrics:   false,
			Inventory: false,
			Events:    false,
		},
		Hostname:              "testhostname",
		Port:                  8091,
		QueryPort:             8093,
		Username:              "testuser",
		Password:              "testpass",
		UseSSL:                false,
		CABundleDir:           "",
		CABundleFile:          "",
		EnableClusterAndNodes: true,
		EnableBuckets:         true,
		EnableBucketStats:     true,
		Timeout:               30,
	}
}
