package entities

import (
	"testing"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-couchbase/src/testutils"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"flag"
	"github.com/stretchr/testify/assert"
)

var (
	update = flag.Bool("update", false, "update .golden files")
)

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
