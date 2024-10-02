package main

import (
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/newrelic/infra-integrations-sdk/v3/data/metric"
	"github.com/newrelic/infra-integrations-sdk/v3/integration"
	"github.com/newrelic/nri-couchbase/src/arguments"
	"github.com/newrelic/nri-couchbase/src/client"
	"github.com/newrelic/nri-couchbase/src/entities"
	"github.com/newrelic/nri-couchbase/src/testutils"
	"github.com/stretchr/testify/assert"
)

func TestStartCollectorWorkerPool(t *testing.T) {
	numWorkers := 10
	var wg sync.WaitGroup
	entitiesChan := StartCollectorWorkerPool(numWorkers, &wg)
	close(entitiesChan)

	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		return
	case <-time.After(time.Second):
		assert.FailNow(t, "Wait group close timed out")
	}
}

type testCollector struct {
	name        string
	integration *integration.Integration
	client      *client.HTTPClient
}

func (t *testCollector) GetEntity() (*integration.Entity, error) {
	if t.integration != nil {
		return t.integration.Entity(t.name, "test")
	}

	return nil, assert.AnError
}

func (t *testCollector) GetName() string {
	return t.name
}

func (t *testCollector) GetIntegration() *integration.Integration {
	return t.integration
}

func (t *testCollector) GetClient() *client.HTTPClient {
	return t.client
}

func (t *testCollector) Collect(collectInventory, collectMetrics bool) error {
	e, err := t.GetEntity()
	if err != nil {
		return err
	}
	if err := e.SetInventoryItem("testitem", "value", "some-attribute"); err != nil {
		return err
	}

	ms := e.NewMetricSet("testSample")
	return ms.SetMetric("test-metric", 17, metric.GAUGE)
}

func Test_collectorWorker(t *testing.T) {
	collectorChan := make(chan entities.Collector)
	var wg sync.WaitGroup
	i, _ := integration.New("testIntegration", "testVersion")

	wg.Add(1)
	go collectorWorker(collectorChan, &wg)

	collectorChan <- &testCollector{
		"testName",
		i,
		&client.HTTPClient{},
	}
	close(collectorChan)

	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		assert.Len(t, i.Entities, 1, "Expected one entity")
		assert.Len(t, i.Entities[0].Metrics[0].Metrics, 2, "Expected one metric in the set")
		assert.Len(t, i.Entities[0].Inventory.Items(), 1, "Expected one inventory item")
	case <-time.After(time.Second):
		assert.FailNow(t, "Collector worker took too long to close.")
	}
}

func Test_FeedWorkerPool(t *testing.T) {
	entities.ClusterName = "couch5"
	endpointMap := map[string]string{
		"/pools/default":         filepath.Join("testdata", "input", "cluster.json"),
		"/pools/default/buckets": filepath.Join("testdata", "input", "buckets.json"),
	}
	testServer := testutils.GetTestServer(t, endpointMap)

	mockedClient := client.HTTPClient{
		Client:       testServer.Client(),
		Username:     "testUser",
		Password:     "testPass",
		BaseURL:      testServer.URL,
		BaseQueryURL: testServer.URL,
	}

	args = arguments.ArgumentList{
		QueryPort:             8093,
		EnableBuckets:         true,
		EnableClusterAndNodes: true,
	}

	collChan := make(chan entities.Collector)
	i, _ := integration.New("test", "0.0.0")

	go FeedWorkerPool(&args, &mockedClient, collChan, i)

	wgDone := make(chan struct{})
	var collectors []entities.Collector
	go func() {
		for {
			coll, ok := <-collChan
			if !ok {
				break
			} else {
				collectors = append(collectors, coll)
			}
		}
		close(wgDone)
	}()

	expectedCollectorNames := map[string]bool{
		// cluster
		"couch5": true,
		// nodes
		"10.33.106.59:8091":                   true,
		"cb50-rh7-1.bluemedora.localnet:8091": true,
		"cb50-rh7-2.bluemedora.localnet:8091": true,
		// query engines
		"10.33.106.59:8093":                   true,
		"cb50-rh7-1.bluemedora.localnet:8093": true,
		"cb50-rh7-2.bluemedora.localnet:8093": true,
		// buckets
		"beer-sample":    true,
		"gamesim-sample": true,
		"travel-sample":  true,
	}

	select {
	case <-wgDone:
		assert.Len(t, collectors, len(expectedCollectorNames))
		for _, coll := range collectors {
			_, ok := expectedCollectorNames[coll.GetName()]
			assert.True(t, ok, "Expected collector name is missing: %s", coll.GetName())
			assert.Equal(t, i, coll.GetIntegration())
			e, err := coll.GetEntity()
			assert.NoError(t, err)
			assert.NotNil(t, e)
		}
	case <-time.After(time.Second):
		assert.FailNow(t, "Timed out waiting for Mongoses")
	}
}
