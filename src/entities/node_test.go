package entities

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/newrelic/nri-couchbase/src/client"
	"github.com/newrelic/nri-couchbase/src/definition"
	"github.com/stretchr/testify/assert"
)

func Test_NodeCollection(t *testing.T) {
	i := getTestingIntegration(t)

	var node definition.Node
	data, _ := ioutil.ReadFile(filepath.Join("..", "testdata", "input", "node.json"))
	_ = json.Unmarshal(data, &node)

	nodeCollector := &nodeCollector{
		defaultCollector{
			name:        "test-node",
			integration: i,
			client:      &client.HTTPClient{},
		},
		node,
		"CouchCluster",
	}

	nodeCollector.Collect(true, true)

	output, _ := i.MarshalJSON()
	goldenFile := filepath.Join("..", "testdata", "node.json")
	writeGoldenFile(t, goldenFile, output)

	expected, _ := ioutil.ReadFile(goldenFile)
	assert.Equal(t, expected, output)
}
