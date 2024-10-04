//go:generate goversioninfo
package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/newrelic/infra-integrations-sdk/v3/integration"
	"github.com/newrelic/infra-integrations-sdk/v3/log"
	"github.com/newrelic/nri-couchbase/src/arguments"
	"github.com/newrelic/nri-couchbase/src/client"
	"github.com/newrelic/nri-couchbase/src/entities"
)

const (
	integrationName = "com.newrelic.couchbase"
)

var (
	args               arguments.ArgumentList
	integrationVersion = "0.0.0"
	gitCommit          = ""
	buildDate          = ""
)

func main() {
	// Create Integration
	i, err := integration.New(integrationName, integrationVersion, integration.Args(&args))
	exitOnError(err)

	if args.ShowVersion {
		fmt.Printf(
			"New Relic %s integration Version: %s, Platform: %s, GoVersion: %s, GitCommit: %s, BuildDate: %s\n",
			strings.Title(strings.Replace(integrationName, "com.newrelic.", "", 1)),
			integrationVersion,
			fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
			runtime.Version(),
			gitCommit,
			buildDate)
		os.Exit(0)
	}

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
