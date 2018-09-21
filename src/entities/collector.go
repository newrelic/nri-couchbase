package entities

import (
	"strconv"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-couchbase/src/arguments"
	"github.com/newrelic/nri-couchbase/src/client"
	"github.com/newrelic/nri-couchbase/src/definition"
)

// inventoryItem is really equivalent to a map element but allows for easier appending and reflection-less iteration
type inventoryItem struct {
	key   string
	value interface{}
}

// Collector is an interface which represents an entity.
// A Collector knows how to collect itself through the Collect method,
// which takes in flags for collecting inventory and metrics.
type Collector interface {
	Collect(bool, bool) error
	GetName() string
	GetEntity() (*integration.Entity, error)
	GetIntegration() *integration.Integration
	GetClient() *client.HTTPClient
}

// defaultCollector is the most basic implementation of the
// Collector interface, and can be inherited to create a minimal
// running version which creates no metrics or inventory
type defaultCollector struct {
	name        string
	integration *integration.Integration
	client      *client.HTTPClient
}

func (d *defaultCollector) GetName() string {
	return d.name
}

// GetIntegration returns the integration associated with the collector
func (d *defaultCollector) GetIntegration() *integration.Integration {
	return d.integration
}

// GetSession returns the session associated with the collector
func (d *defaultCollector) GetClient() *client.HTTPClient {
	return d.client
}

// GetClusterCollectors returns a slice of collectors, one for the cluster and one for each node.
// Each collector collects metrics and inventory for its entity
func GetClusterCollectors(args *arguments.ArgumentList, i *integration.Integration, nodeClient *client.HTTPClient) ([]Collector, error) {

	var clusterDetails definition.PoolsDefaultResponse
	err := nodeClient.Request("/pools/default", &clusterDetails)
	if err != nil {
		return nil, err
	}

	collectors := make([]Collector, 0, 10)
	clusterCollector := &clusterCollector{
		defaultCollector{
			name:        "cluster-test",
			client:      nodeClient,
			integration: i,
		},
		args.Hostname,
		&clusterDetails,
	}

	collectors = append(collectors, clusterCollector)

	for _, node := range *clusterDetails.Nodes {
		nodeCollector := &nodeCollector{
			defaultCollector{
				name:        *node.Hostname,
				client:      nodeClient,
				integration: i,
			},
			&node,
		}

		collectors = append(collectors, nodeCollector)

		// check for query engine
		for _, service := range *node.Services {
			if service == "n1ql" {
				// create client for new host
				nodeHost := strings.Split(*node.Hostname, ":")[0]
				queryEngineClient, err := client.CreateClient(args, nodeHost)
				if err != nil {
					log.Error("Could not create client for query engine on node '%s': %v", *node.Hostname, err)
					continue
				}
				log.Info("Created new client for node %s", *node.Hostname)
				queryEngineCollector := &queryEngineCollector{
					defaultCollector{
						name:        nodeHost+":"+strconv.Itoa(args.QueryPort),
						client:      queryEngineClient,
						integration: i,
					},
				}
				collectors = append(collectors, queryEngineCollector)
			}
		}
	}

	return collectors, nil
}
