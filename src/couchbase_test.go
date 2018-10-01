package main

import (
	"net/http/httptest"
	"testing"

	"flag"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-couchbase/src/arguments"
	"github.com/newrelic/nri-couchbase/src/client"
	"github.com/newrelic/nri-couchbase/src/testutils"
	"github.com/stretchr/testify/assert"
)

var (
	update = flag.Bool("update", false, "update .golden files")
)

func writeGoldenFile(t *testing.T, goldenPath string, data []byte) {
	if *update {
		t.Log("Writing .golden file")
		err := ioutil.WriteFile(goldenPath, data, 0644)
		assert.NoError(t, err)
	}
}

func Test_EndToEnd(t *testing.T) {
	testServ, testClient := getMappedMockServerAndClient(t)
	defer testServ.Close()

	testIntegration, _ := integration.New("test", "0.1.0")

	collect(testIntegration, testClient)

	output, _ := testIntegration.MarshalJSON()
	// scrub test server hostname and port since it changes from run to run
	regex := regexp.MustCompile(`localhost:\d+`)
	output = regex.ReplaceAll(output, []byte("test-server"))

	regex = regexp.MustCompile(`config/port\":\{\"value\":\d+`)
	output = regex.ReplaceAll(output, []byte("config/port\":{\"value\":13131"))

	goldenFile := filepath.Join("testdata", "full-collection.json")
	writeGoldenFile(t, goldenFile, output)

	expected, _ := ioutil.ReadFile(goldenFile)
	assert.Equal(t, expected, output)
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
