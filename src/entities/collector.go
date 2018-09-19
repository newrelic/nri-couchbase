package entities

import (
	"github.com/newrelic/nri-couchbase/src/definition"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-couchbase/src/client"
	"github.com/newrelic/infra-integrations-sdk/log"
)

// Collector is an interface which represents an entity.
// A Collector knows how to collect itself through the CollectMetrics
// and CollectInventory methods.
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
func GetClusterCollectors(i *integration.Integration, client *client.HTTPClient) ([]Collector, error) {

	var clusterDetails definition.PoolsDefaultResponse
	err := client.Request("/pools/default", &clusterDetails)
	if err != nil {
		return nil, err
	}

	collectors := make([]Collector, 0, 10)
	clusterCollector := &clusterCollector{
		defaultCollector{
			name: "cluster-test",
			client: client,
			integration: i,
		},
		&clusterDetails,
	}

	collectors = append(collectors, clusterCollector)

	for _, node := range *clusterDetails.Nodes {
		log.Info("Creating node... %s", *node.Hostname)
		nodeCollector := &nodeCollector{
			defaultCollector{
				name: *node.Hostname,
				client: client,
				integration: i,
			},
			&node,
		}

		collectors = append(collectors, nodeCollector)
	}

	return collectors, nil
}