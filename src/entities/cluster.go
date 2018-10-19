package entities

import (
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-couchbase/src/client"
	"github.com/newrelic/nri-couchbase/src/definition"
)

type clusterCollector struct {
	defaultCollector
	collectionHostname string
	clusterDetails     *definition.PoolsDefaultResponse
}

func (c *clusterCollector) GetEntity() (*integration.Entity, error) {
	return c.GetIntegration().Entity(*c.clusterDetails.ClusterName, "cluster")
}

// CollectCluster creates entities for the cluster and its nodes,
// adding inventory and metrics according to flags
func (c *clusterCollector) Collect(collectInventory bool, collectMetrics bool) error {
	clusterResponse, autoFailoverResponse, err := getClusterResponses(c.GetClient())
	if err != nil {
		return err
	}

	clusterEntity, err := c.GetEntity()
	if err != nil {
		return err
	}

	if collectInventory {
		collectClusterInventory(clusterEntity, clusterResponse, c.collectionHostname)
	}

	if collectMetrics {
		collectClusterMetrics(clusterEntity, c.clusterDetails, autoFailoverResponse)
	}

	return nil
}

func getClusterResponses(client *client.HTTPClient) (clusterResponse *definition.PoolsResponse, autoFailoverResponse *definition.AutoFailover, err error) {
	clusterResponse = new(definition.PoolsResponse)
	autoFailoverResponse = new(definition.AutoFailover)

	err = client.Request("/pools", &clusterResponse)
	if err != nil {
		return
	}

	err = client.Request("/settings/autoFailover", &autoFailoverResponse)
	return
}

func collectClusterInventory(clusterEntity *integration.Entity, clusterResponse *definition.PoolsResponse, hostname string) {
	inventoryItems := []inventoryItem{
		{"couchbaseVersion", clusterResponse.Version},
		{"clusterUUID", clusterResponse.UUID},
		{"collectionNode", hostname},
	}

	for _, inventoryItem := range inventoryItems {
		if err := clusterEntity.SetInventoryItem(inventoryItem.key, "value", inventoryItem.value); err != nil {
			log.Error("Could not set inventory item '%s' on cluster entity: %v", inventoryItem.key, err)
		}
	}
}

func collectClusterMetrics(clusterEntity *integration.Entity, clusterDetailsResponse *definition.PoolsDefaultResponse, autoFailoverResponse *definition.AutoFailover) {
	clusterMetricSet := clusterEntity.NewMetricSet("CouchbaseClusterSample",
		metric.Attribute{Key: "displayName", Value: clusterEntity.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: clusterEntity.Metadata.Namespace + ":" + clusterEntity.Metadata.Name},
	)

	err := clusterMetricSet.MarshalMetrics(clusterDetailsResponse)
	if err != nil {
		log.Error("Could not marshal cluster metrics from endpoint '/pools/default': %v", err)
	}

	err = clusterMetricSet.MarshalMetrics(autoFailoverResponse)
	if err != nil {
		log.Error("Could not marshal cluster metrics from endpoint '/settings/autoFailover': %v", err)
	}
}
