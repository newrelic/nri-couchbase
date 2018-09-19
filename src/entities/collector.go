package entities

import (
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-couchbase/src/client"
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

func GetClusterCollector(i *integration.Integration, client *client.HTTPClient) Collector {
	return &clusterCollector{
		defaultCollector{
			name: "cluster-test",
			client: client,
			integration: i,
		},
	}
}