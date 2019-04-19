package entities

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-couchbase/src/arguments"
	"github.com/newrelic/nri-couchbase/src/client"
	"github.com/newrelic/nri-couchbase/src/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	update = flag.Bool("update", false, "update .golden files")
)

func getTestingIntegration(t *testing.T) *integration.Integration {
	payload, err := integration.New("Test", "0.0.1")
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
	testCases := []struct {
		clusterDataPath     string
		expectedClusterName string
	}{
		{"cluster.json", "couch5"},
		{"pools-default-missing-clustername.json", "default"},
	}

	for _, tc := range testCases {
		endpointMap := map[string]string{
			"/pools/default":         filepath.Join("..", "testdata", "input", tc.clusterDataPath),
			"/pools/default/buckets": filepath.Join("..", "testdata", "input", "buckets.json"),
		}
		testServer := testutils.GetTestServer(t, endpointMap)
		defer testServer.Close()

		i := getTestingIntegration(t)
		args := arguments.ArgumentList{}
		client := &client.HTTPClient{
			Client:   testServer.Client(),
			Username: "testUser",
			Password: "testPass",
			BaseURL:  testServer.URL,
		}

		clusterCollectors, err := GetClusterCollectors(&args, i, client)
		assert.NoError(t, err)
		assert.Equal(t, 7, len(clusterCollectors))

		// find cluster collector
		for _, collector := range clusterCollectors {
			e, _ := collector.GetEntity()
			if e.Metadata.Namespace != "cluster" {
				continue
			}
			assert.Equal(t, tc.expectedClusterName, e.Metadata.Name)
		}

		bucketCollectors, err := GetBucketCollectors(&args, i, client)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(bucketCollectors))
	}
}
