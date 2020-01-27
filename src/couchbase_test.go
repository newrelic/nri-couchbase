package main

import (
	"flag"
	"io/ioutil"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-couchbase/src/arguments"
	"github.com/newrelic/nri-couchbase/src/client"
	"github.com/newrelic/nri-couchbase/src/entities"
	"github.com/newrelic/nri-couchbase/src/testutils"
	"github.com/stretchr/testify/assert"
)

var (
	update = flag.Bool("update", true, "update .golden files")
)

func writeGoldenFile(t *testing.T, goldenPath string, data []byte) {
	if *update {
		t.Log("Writing .golden file")
		err := ioutil.WriteFile(goldenPath, data, 0644)
		assert.NoError(t, err)
	}
}

func Test_EndToEnd(t *testing.T) {
	entities.ClusterName = "testcluster"
	testServ, testClient := getMappedMockServerAndClient(t)
	defer testServ.Close()

	testIntegration, _ := integration.New("test", "0.1.0")

	collect(testIntegration, testClient)

	// should have 4 total entities, one of each type
	assert.Equal(t, 4, len(testIntegration.Entities))
	counts := map[string]int{
		"cb-queryEngine": 0,
		"cb-node":        0,
		"cb-cluster":     0,
		"cb-bucket":      0,
	}
	for _, entity := range testIntegration.Entities {
		counts[entity.Metadata.Namespace]++
	}
	for e, count := range counts {
		assert.Equal(t, 1, count, "Wrong number of entities for type %s", e)
	}
}

func getMappedMockServerAndClient(t *testing.T) (*httptest.Server, *client.HTTPClient) {
	endpointMap := map[string]string{
		"/pools":                                     filepath.Join("testdata", "input", "end-to-end", "pools.json"),
		"/pools/default":                             filepath.Join("testdata", "input", "end-to-end", "pools-default.json"),
		"/pools/default/buckets":                     filepath.Join("testdata", "input", "end-to-end", "pools-default-buckets.json"),
		"/pools/default/buckets/sample-bucket/stats": filepath.Join("testdata", "input", "end-to-end", "bucket-stats.json"),
		"/admin/settings":                            filepath.Join("testdata", "input", "end-to-end", "admin-settings.json"),
		"/admin/vitals":                              filepath.Join("testdata", "input", "end-to-end", "admin-vitals.json"),
		"/settings/autoFailover":                     filepath.Join("testdata", "input", "end-to-end", "auto-failover.json"),
	}

	testServer := testutils.GetTestServer(t, endpointMap)
	hostnamePort := strings.Split(strings.Split(testServer.URL, "://")[1], ":")
	hostname := hostnamePort[0]
	port, _ := strconv.Atoi(hostnamePort[1])

	args = arguments.ArgumentList{
		Hostname:              hostname,
		Port:                  port,
		QueryPort:             port,
		Username:              "testUser",
		Password:              "testPass",
		EnableBuckets:         true,
		EnableBucketStats:     true,
		EnableClusterAndNodes: true,
	}

	client := client.HTTPClient{
		Client:       testServer.Client(),
		Username:     "testUser",
		Password:     "testPass",
		Hostname:     hostname,
		Port:         port,
		QueryPort:    port,
		BaseURL:      testServer.URL,
		BaseQueryURL: testServer.URL,
	}

	return testServer, &client
}
