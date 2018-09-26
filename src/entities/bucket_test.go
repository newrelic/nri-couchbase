package entities

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-couchbase/src/client"
	"github.com/newrelic/nri-couchbase/src/definition"
	"github.com/stretchr/testify/assert"
)

func Test_BucketMetrics(t *testing.T) {
	bucketStats := definition.PoolsDefaultBucket{
		ReplicaNumber: new(int),
	}
	bucketExtended := definition.BucketStats{
		Op: &definition.OpStats{
			Samples: &definition.SampleStats{
				BytesRead: &[]float64{5.5, 7.5, 1.0, 0.0},
			},
		},
	}

	i := getTestingIntegration(t)
	e, _ := i.Entity("test", "testEntity")

	metricSet := collectBucketMetrics(e, &bucketStats)
	collectExtendedBucketMetrics(metricSet, &bucketExtended, "test")

	output, _ := i.MarshalJSON()

	goldenFile := filepath.Join("..", "testdata", "bucket-metrics.json")
	writeGoldenFile(t, goldenFile, output)

	expected, _ := ioutil.ReadFile(goldenFile)
	assert.Equal(t, output, expected)
}

func createBucketCollector(i *integration.Integration) *bucketCollector {
	bucketStats := definition.PoolsDefaultBucket{
		ReplicaNumber: new(int),
	}
	return &bucketCollector{
		defaultCollector{
			name:        "test-collector",
			integration: i,
			client:      &client.HTTPClient{},
		},
		&bucketStats,
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
