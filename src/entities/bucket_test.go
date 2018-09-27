package entities

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-couchbase/src/client"
	"github.com/newrelic/nri-couchbase/src/definition"
	"github.com/newrelic/nri-couchbase/src/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_FunctionalCollection(t *testing.T) {
	endpointMap := map[string]string{
		"/pools/default/buckets/beer-sample/stats": filepath.Join("..", "testdata", "input", "bucket-stats.json"),
	}
	testServer := testutils.GetTestServer(t, endpointMap)
	defer testServer.Close()

	i := getTestingIntegration(t)

	testClient := testServer.Client()
	collector := createBucketCollector(i, testClient, testServer.URL)

	collector.Collect(true, true)

	output, _ := i.MarshalJSON()
	goldenPath := filepath.Join("..", "testdata", "bucket.json")
	writeGoldenFile(t, goldenPath, output)

	expected, _ := ioutil.ReadFile(goldenPath)
	assert.Equal(t, expected, output)
}

func Test_BucketMetrics(t *testing.T) {
	bucketStats := createBucketResponse()

	var bucketExtended definition.BucketStats
	data, _ := ioutil.ReadFile(filepath.Join("..", "testdata", "input", "bucket-stats.json"))
	_ = json.Unmarshal(data, &bucketExtended)

	i := getTestingIntegration(t)
	e, _ := i.Entity("test", "testEntity")

	metricSet := collectBucketMetrics(e, bucketStats)
	collectExtendedBucketMetrics(metricSet, &bucketExtended, "test")

	output, _ := i.MarshalJSON()

	goldenFile := filepath.Join("..", "testdata", "bucket-metrics.json")
	writeGoldenFile(t, goldenFile, output)

	expected, _ := ioutil.ReadFile(goldenFile)
	assert.Equal(t, output, expected)
}

func Test_BucketInventory(t *testing.T) {
	var bucketStats definition.PoolsDefaultBucket
	data, _ := ioutil.ReadFile(filepath.Join("..", "testdata", "input", "bucket.json"))
	_ = json.Unmarshal(data, &bucketStats)

	i := getTestingIntegration(t)
	e, _ := i.Entity("test", "testEntity")

	collectBucketInventory(e, &bucketStats)

	output, _ := i.MarshalJSON()

	goldenFile := filepath.Join("..", "testdata", "bucket-inventory.json")
	writeGoldenFile(t, goldenFile, output)

	expected, _ := ioutil.ReadFile(goldenFile)
	assert.Equal(t, output, expected)
}

func createBucketCollector(i *integration.Integration, httpClient *http.Client, url string) *bucketCollector {
	return &bucketCollector{
		defaultCollector{
			name:        "beer-sample",
			integration: i,
			client: &client.HTTPClient{
				Client:   httpClient,
				Username: "testUser",
				Password: "testPass",
				BaseURL:  url,
			},
		},
		createBucketResponse(),
		true,
	}
}

func createBucketExtendedResponse() *definition.BucketStats {
	bucketExtended := definition.BucketStats{
		Op: &definition.OpStats{
			Samples: &definition.SampleStats{
				BytesRead: &[]float64{5.5, 7.5, 1.0, 0.0},
			},
		},
	}
	return &bucketExtended
}

func createBucketResponse() *definition.PoolsDefaultBucket {
	var bucketStats definition.PoolsDefaultBucket
	data, _ := ioutil.ReadFile(filepath.Join("..", "testdata", "input", "bucket.json"))
	_ = json.Unmarshal(data, &bucketStats)

	return &bucketStats
}
