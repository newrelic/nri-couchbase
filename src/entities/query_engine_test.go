package entities

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/newrelic/nri-couchbase/src/client"
	"github.com/newrelic/nri-couchbase/src/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_QueryEngineCollection(t *testing.T) {
	dataMap := map[string]string{
		"/admin/settings": filepath.Join("..", "testdata", "input", "admin-settings.json"),
		"/admin/vitals":   filepath.Join("..", "testdata", "input", "admin-vitals.json"),
	}
	testServer := testutils.GetTestServer(t, dataMap)
	defer testServer.Close()

	i := getTestingIntegration(t)

	qeCollector := &queryEngineCollector{
		defaultCollector{
			name:        "test-query-engine",
			integration: i,
			client: &client.HTTPClient{
				Client:       testServer.Client(),
				Username:     "testUser",
				Password:     "testPass",
				BaseQueryURL: testServer.URL,
			},
		},
		"CouchCluser",
	}

	assert.NoError(t, qeCollector.Collect(true, true))

	output, _ := i.MarshalJSON()
	goldenFile := filepath.Join("..", "testdata", "query-engine.json")
	writeGoldenFile(t, goldenFile, output)

	expected, _ := ioutil.ReadFile(goldenFile)
	assert.Equal(t, expected, output)
}

func Test_TimeConversion(t *testing.T) {
	testCases := []struct {
		timeString string
		expectedMs float64
	}{
		{"1s", 1000},
		{"3s", 3000},
		{"1.03s", 1030},
		{"1m", 60000},
		{"1m1s", 61000},
		{"3h20m1.04s", 12001040},
		{"2d", 172800000},
		{"5ms", 5},
		{"7µs", 0.007},
		{"7us", 0.007},
		{"6700ns", 0.0067},
	}

	for _, tc := range testCases {
		actual, err := convertTimeUnits(tc.timeString)
		assert.NoError(t, err)
		assert.Equal(t, tc.expectedMs, actual)
	}
}

func Test_TimeConversionErrors(t *testing.T) {
	testCases := []string{
		"50ee7l",
	}

	for _, tc := range testCases {
		_, err := convertTimeUnits(tc)
		assert.Error(t, err)
	}
}
