package main

import (
	"os"
	"sync"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-couchbase/src/arguments"
	"github.com/newrelic/nri-couchbase/src/client"
	"github.com/newrelic/nri-couchbase/src/entities"
)

const (
	integrationName    = "com.newrelic.couchbase"
	integrationVersion = "1.0.2"
)

var (
	args arguments.ArgumentList
)

func main() {
	// Create Integration
	i, err := integration.New(integrationName, integrationVersion, integration.Args(&args))
	exitOnError(err)

	log.SetupLogging(args.Verbose)

	client, err := client.CreateClient(&args, "")
	exitOnError(err)

	err = entities.SetClusterName(client)
	exitOnError(err)

	collect(i, client)

	exitOnError(i.Publish())
}

func collect(i *integration.Integration, client *client.HTTPClient) {
	// create worker pool
	// Start workers
	var wg sync.WaitGroup
	collectorChan := StartCollectorWorkerPool(10, &wg)

	// Feed the worker pool with entities to be collected
	go FeedWorkerPool(&args, client, collectorChan, i)

	// Wait for workers to finish
	wg.Wait()
}

func exitOnError(err error) {
	if err != nil {
		log.Error("Could not complete collection: %v", err)
		os.Exit(1)
	}
}
