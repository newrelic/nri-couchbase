package entities

import (
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-couchbase/src/client"
	"github.com/newrelic/nri-couchbase/src/definition"
)

type clusterCollector struct {
	defaultCollector
}

// CollectCluster creates entities for the cluster and its nodes,
// adding inventory and metrics according to flags
func (c *clusterCollector) Collect(collectInventory bool, collectMetrics bool) error {
	clusterResponse, clusterDetails, err := getClusterResponses(c.GetClient())
	if err != nil {
		return err
	}	

	clusterEntity, err := c.GetIntegration().Entity(*clusterDetails.ClusterName, "cluster")
	if err != nil {
		return err
	}

	if collectInventory {
		collectClusterInventory(clusterEntity, clusterResponse)
	}

	if collectMetrics {
		collectClusterMetrics(clusterEntity, clusterDetails)
	}

	return nil
}

func getClusterResponses(client *client.HTTPClient) (clusterResponse *definition.PoolsResponse, clusterDetailsResponse *definition.PoolsDefaultResponse, err error) {
	clusterResponse = new(definition.PoolsResponse)
	clusterDetailsResponse = new(definition.PoolsDefaultResponse)

	err = client.Request("/pools", &clusterResponse)
	if err != nil {
		return
	}

	err = client.Request("/pools/default", &clusterDetailsResponse)
	return
}

func collectClusterInventory(clusterEntity *integration.Entity, clusterResponse *definition.PoolsResponse) {
	inventoryItems := []struct{
		key string
		value interface{}
	}{
		{"couchbaseVersion", clusterResponse.Version},
		{"clusterUUID", clusterResponse.UUID},
	}

	for _, inventoryItem := range inventoryItems {
		if err := clusterEntity.SetInventoryItem("config/"+inventoryItem.key, "value", inventoryItem.value); err != nil {
			log.Error("Could not set inventory item '%s' on cluster entity: %v", inventoryItem.key, err)
		}
	}
}

func collectClusterMetrics(clusterEntity *integration.Entity, clusterDetailsResponse *definition.PoolsDefaultResponse) {
	clusterMetricSet := clusterEntity.NewMetricSet("CouchbaseClusterSample",
		metric.Attribute{Key: "displayName", Value: clusterEntity.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: clusterEntity.Metadata.Namespace + ":" + clusterEntity.Metadata.Name},
	)

	err := clusterMetricSet.MarshalMetrics(clusterDetailsResponse)
	if err != nil {
		log.Error("Could not marshal cluster metrics")
	}
}

func (c *clusterCollector) GetEntity() (*integration.Entity, error) {
	e, err := c.GetIntegration().Entity(c.name, "cluster")
	if err != nil {
		return nil, err
	}

	return e, nil
}