package entities

import (
	"github.com/newrelic/nri-couchbase/src/definition"
	"testing"
	"net/http"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"github.com/stretchr/testify/assert"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-couchbase/src/client"
)

func Test_ClusterCollection(t *testing.T) {
	endpointMap := map[string]string{
		"/pools": filepath.Join("..", "testdata", "input", "pools.json"),
		"/settings/autoFailover": filepath.Join("..", "testdata", "input", "auto-failover.json"),
	}
	testServer := getTestServer(t, endpointMap)
	defer testServer.Close()

	i := getTestingIntegration(t)
	collector := createClusterCollector(i, testServer.Client(), testServer.URL)

	collector.Collect(true, true)

	output, _ := i.MarshalJSON()
	goldenPath := filepath.Join("..", "testdata", "cluster.json")
	writeGoldenFile(t, goldenPath, output)

	expected, _ := ioutil.ReadFile(goldenPath)
	assert.Equal(t, expected, output)
}

func createClusterCollector(i *integration.Integration, httpClient *http.Client, url string) *clusterCollector {
	var poolsDefault definition.PoolsDefaultResponse
	data, _ := ioutil.ReadFile(filepath.Join("..", "testdata", "input", "pools-default.json"))
	json.Unmarshal(data, &poolsDefault)

	return &clusterCollector{
		defaultCollector{
			name:        "cluster-name",
			integration: i,
			client: &client.HTTPClient{
				Client:   httpClient,
				Username: "testUser",
				Password: "testPass",
				BaseURL:  url,
			},
		},
		url,
		&poolsDefault,
	}
}