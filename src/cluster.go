package main

import (
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
)

// CollectCluster creates entities for the cluster and its nodes,
// adding inventory and metrics according to flags
func CollectCluster(i *integration.Integration, client *HTTPClient) error {
	clusterResponse, clusterDetails, err := getClusterResponses(client)
	if err != nil {
		return err
	}	

	clusterEntity, err := i.Entity(*clusterDetails.ClusterName, "cluster")
	if err != nil {
		return err
	}

	if args.HasInventory() {
		collectClusterInventory(clusterEntity, clusterResponse)
	}

	if args.HasMetrics() {
		collectClusterMetrics(clusterEntity, clusterDetails)
	}

	return nil
}

func getClusterResponses(client *HTTPClient) (clusterResponse *PoolsResponse, clusterDetailsResponse *PoolsDefaultResponse, err error) {
	clusterResponse = new(PoolsResponse)
	clusterDetailsResponse = new(PoolsDefaultResponse)

	err = client.Request("/pools", &clusterResponse)
	if err != nil {
		return
	}

	err = client.Request("/pools/default", &clusterDetailsResponse)
	return
}

func collectClusterInventory(clusterEntity *integration.Entity, clusterResponse *PoolsResponse) {
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

func collectClusterMetrics(clusterEntity *integration.Entity, clusterDetailsResponse *PoolsDefaultResponse) {
	clusterMetricSet := clusterEntity.NewMetricSet("CouchbaseClusterSample",
		metric.Attribute{Key: "displayName", Value: clusterEntity.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: clusterEntity.Metadata.Namespace + ":" + clusterEntity.Metadata.Name},
	)

	err := clusterMetricSet.MarshalMetrics(clusterDetailsResponse)
	if err != nil {
		log.Error("Could not marshal cluster metrics")
	}
}