package main

import (
	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	Hostname              string `default:"localhost" help:"The hostname or IP of the Couchbase node being monitored"`
	Port                  int    `default:"8091" help:"The port used to connect to the Couchbase API"`
	QueryPort			  int    `default:"8093" help:"The port used to connect to the N1QL service"`
	Username              string `default:"" help:"The username used to connect to the Couchbase API"`
	Password              string `default:"" help:"The password used to connect to the Couchbase API"`
	UseSSL                bool   `default:"false" help:"Signals whether to use SSL or not. Certificate bundle must be supplied"`
	CABundleFile          string `default:"" help:"Alternative Certificate Authority bundle file"`
	CABundleDir           string `default:"" help:"Alternative Certificate Authority bundle directory"`
	EnableClusterAndNodes bool   `default:"true" help:"If true, collects cluster and node resources"`
	EnableBuckets         bool   `default:"true" help:"If true, collects bucket resources"`
	EnableBucketStats     bool   `default:"true" help:"If true, collects additional bucket statistics"`
	Timeout               int    `default:"30" help:"Timeout for an API call"`
}

const (
	integrationName    = "com.newrelic.couchbase"
	integrationVersion = "0.1.0"
)

var (
	args argumentList
)

func main() {
	// Create Integration
	i, err := integration.New(integrationName, integrationVersion, integration.Args(&args))
	panicOnErr(err)

	client, err := CreateClient()
	panicOnErr(err)

	// collect cluster/nodes?
	CollectCluster(i, client)
	// collect buckets?

	panicOnErr(i.Publish())
}

// checkErr logs an error if it exists
func checkErr(f func() error) {
	if err := f(); err != nil {
		log.Error("%v", err)
	}
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
