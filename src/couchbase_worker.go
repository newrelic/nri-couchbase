package main

import (
	"sync"
	"github.com/newrelic/nri-couchbase/src/entities"
	"github.com/newrelic/nri-couchbase/src/client"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
)

// StartCollectorWorkerPool starts a pool of workers to handle collecting each entity
// and returns a channel of collectors which the workers read off of
func StartCollectorWorkerPool(numWorkers int, wg *sync.WaitGroup) chan entities.Collector {
	wg.Add(numWorkers)

	collectorChan := make(chan entities.Collector, 100)
	for j := 0; j < numWorkers; j++ {
		go collectorWorker(collectorChan, wg)
	}

	return collectorChan
}

// collectorWorker reads a collector from the collector chan, then asynchronously
// collects that object's inventory and metrics
func collectorWorker(collectorChan chan entities.Collector, wg *sync.WaitGroup) {
	defer wg.Done()

	// Loop until collectorChan is empty and closed
	for {
		collector, ok := <-collectorChan
		if !ok {
			return
		}

		collector.Collect(args.HasInventory(), args.HasMetrics())
	}
}

// FeedWorkerPool feeds the workers with the collectors that contain the info needed to collect each entity
func FeedWorkerPool(client *client.HTTPClient, collectorChan chan entities.Collector, integration *integration.Integration) {
	defer close(collectorChan)

	// Create a wait group for each of the get*Collectors calls
	getWg := new(sync.WaitGroup)

	getWg.Add(1)
	go createClusterAndNodeCollectors(getWg, client, collectorChan, integration)

	getWg.Wait()
}

func createClusterAndNodeCollectors(wg *sync.WaitGroup, client *client.HTTPClient, channel chan entities.Collector, integration *integration.Integration) {
	defer wg.Done()
	clusterAndNodeCollectors, err := entities.GetClusterCollectors(integration, client)
	if err != nil {
		log.Error("Could not create cluster and node collectors: %v", err)
	}
	for _, collector := range clusterAndNodeCollectors {
		channel <- collector
	}
}

func createNodeCollectors(wg *sync.WaitGroup, client *client.HTTPClient, channel chan entities.Collector, integration *integration.Integration) {
	defer wg.Done()

	// get list of nodes
	// spin up collectors for each.
}