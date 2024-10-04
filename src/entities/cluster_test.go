package entities

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/newrelic/nri-couchbase/src/definition"
	"github.com/stretchr/testify/assert"

	"github.com/newrelic/infra-integrations-sdk/v3/integration"
	"github.com/newrelic/nri-couchbase/src/client"

	"github.com/newrelic/nri-couchbase/src/testutils"
)

func Test_ClusterCollection(t *testing.T) {
	endpointMap := map[string]string{
		"/pools":                 filepath.Join("..", "testdata", "input", "pools.json"),
		"/settings/autoFailover": filepath.Join("..", "testdata", "input", "auto-failover.json"),
	}
	testServer := testutils.GetTestServer(t, endpointMap)
	defer testServer.Close()

	i := getTestingIntegration(t)
	collector := createClusterCollector(i, testServer.Client(), testServer.URL)

	// set the ClusterName global variable for the test
	ClusterName = "test"
	assert.NoError(t, collector.Collect(true, true))

	output, _ := i.MarshalJSON()
	goldenPath := filepath.Join("..", "testdata", "cluster.json")
	writeGoldenFile(t, goldenPath, output)

	expected, _ := ioutil.ReadFile(goldenPath)
	assert.Equal(t, expected, output)
}

func createClusterCollector(i *integration.Integration, httpClient *http.Client, url string) *clusterCollector {
	var poolsDefault definition.PoolsDefaultResponse
	data, _ := ioutil.ReadFile(filepath.Join("..", "testdata", "input", "pools-default.json"))
	json.Unmarshal(data, &poolsDefault) //nolint

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
		"some-hostname",
		&poolsDefault,
	}
}

func Test_sanitizeAutocompactionSettings(t *testing.T) {
	input := &definition.PoolsDefaultResponse{
		AutoCompactionSettings: &definition.AutoCompactionSettings{
			DatabaseFragmentationThreshold: &definition.DatabaseFragmentationThreshold{
				Percentage: 46,
			},
			IndexFragmentationThreshold: &definition.IndexFragmentationThreshold{
				Percentage: "undefined",
			},
			ViewFragmentationThreshold: &definition.ViewFragmentationThreshold{
				Percentage: 68.8,
			},
		},
	}

	expected := &definition.PoolsDefaultResponse{
		AutoCompactionSettings: &definition.AutoCompactionSettings{
			DatabaseFragmentationThreshold: &definition.DatabaseFragmentationThreshold{
				Percentage: 46,
			},
			IndexFragmentationThreshold: &definition.IndexFragmentationThreshold{
				Percentage: nil,
			},
			ViewFragmentationThreshold: &definition.ViewFragmentationThreshold{
				Percentage: nil,
			},
		},
	}
	sanitizeAutocompactionSettings(input)

	assert.Equal(t, expected, input)

}
