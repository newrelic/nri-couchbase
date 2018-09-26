package entities

import (
	"flag"
	"io/ioutil"
	"testing"
	"net/http/httptest"
	"net/http"
	"path/filepath"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-couchbase/src/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/newrelic/nri-couchbase/src/arguments"
	"github.com/newrelic/nri-couchbase/src/client"
)

var (
	update = flag.Bool("update", false, "update .golden files")
)

func getTestServer(t *testing.T, dataMap map[string]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		username, password, ok := req.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, username, "testUser")
		assert.Equal(t, password, "testPass")

		endpoint := req.RequestURI
		filepath, ok := dataMap[endpoint]
		if !ok {
			t.Errorf("bad request, was not expecting request for endpoint %s", endpoint)
		}
		data, _ := ioutil.ReadFile(filepath)
		res.Write(data)
	}))
}

func getTestingIntegration(t *testing.T) *integration.Integration {
	payload, err := integration.New("Test", "0.0.1", integration.Logger(&testutils.TestLogger{F: t.Logf}))
	require.NoError(t, err)
	require.NotNil(t, payload)
	return payload
}

func writeGoldenFile(t *testing.T, goldenPath string, data []byte) {
	if *update {
		t.Log("Writing .golden file")
		err := ioutil.WriteFile(goldenPath, data, 0644)
		assert.NoError(t, err)
	}
}

func Test_GetCollectors(t *testing.T) {
	endpointMap := map[string]string{
		"/pools/default": filepath.Join("..", "testdata", "input", "cluster.json"),
		"/pools/default/buckets": filepath.Join("..", "testdata", "input", "buckets.json"),
	}
	testServer := getTestServer(t, endpointMap)
	defer testServer.Close()

	i := getTestingIntegration(t)
	args := arguments.ArgumentList{}
	client := &client.HTTPClient{
		Client: testServer.Client(),
		Username: "testUser",
		Password: "testPass",
		BaseURL: testServer.URL,
	}

	clusterCollectors, err := GetClusterCollectors(&args, i, client)
	assert.NoError(t, err)
	assert.Equal(t, 7, len(clusterCollectors))

	bucketCollectors, err := GetBucketCollectors(&args, i, client)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(bucketCollectors))
}