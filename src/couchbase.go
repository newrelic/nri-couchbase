package main

import (
	"sync"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-couchbase/src/arguments"
	"github.com/newrelic/nri-couchbase/src/client"
)

const (
	integrationName    = "com.newrelic.couchbase"
	integrationVersion = "0.1.0"
)

var (
	args arguments.ArgumentList
)

func main() {
	// Create Integration
	i, err := integration.New(integrationName, integrationVersion, integration.Args(&args))
	panicOnErr(err)

	client, err := client.CreateClient(&args, "")
	panicOnErr(err)

	// create worker pool
	// Start workers
	var wg sync.WaitGroup
	collectorChan := StartCollectorWorkerPool(100, &wg)

	// Feed the worker pool with entities to be collected
	go FeedWorkerPool(&args, client, collectorChan, i)

	// Wait for workers to finish
	wg.Wait()

	// collect cluster/nodes?
	//CollectCluster(i, client)
	// collect buckets?

	panicOnErr(i.Publish())
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
