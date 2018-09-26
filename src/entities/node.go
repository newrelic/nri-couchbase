package entities

import (
	"strings"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-couchbase/src/definition"
)

type nodeCollector struct {
	defaultCollector
	nodeDetails *definition.Node
}

func (n *nodeCollector) GetEntity() (*integration.Entity, error) {
	return n.GetIntegration().Entity(n.GetName(), "node")
}

func (n *nodeCollector) Collect(collectInventory, collectMetrics bool) error {
	nodeEntity, err := n.GetEntity()
	if err != nil {
		return err
	}

	if collectMetrics {
		collectNodeMetrics(nodeEntity, n.nodeDetails)
	}

	if collectInventory {
		collectNodeInventory(nodeEntity, n.nodeDetails)
	}

	return nil
}

func collectNodeMetrics(nodeEntity *integration.Entity, nodeResponse *definition.Node) {
	nodeMetricSet := nodeEntity.NewMetricSet("CouchbaseNodeSample",
		metric.Attribute{Key: "displayName", Value: nodeEntity.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: nodeEntity.Metadata.Namespace + ":" + nodeEntity.Metadata.Name},
	)

	err := nodeMetricSet.MarshalMetrics(nodeResponse)
	if err != nil {
		log.Error("Could not marshal base node metrics: %v", err)
	}
}

func collectNodeInventory(nodeEntity *integration.Entity, nodeResponse *definition.Node) {
	inventoryItems := []inventoryItem{
		{"clusterMembership", nodeResponse.ClusterMembership},
		{"os", nodeResponse.OS},
		{"recoveryType", nodeResponse.RecoveryType},
		{"services", strings.Join(*nodeResponse.Services, ", ")},
		{"version", nodeResponse.Version},
	}

	splitHostname := strings.Split(*nodeResponse.Hostname, ":")
	if len(splitHostname) == 2 {
		inventoryItems = append(inventoryItems, []inventoryItem{
			{"hostname", splitHostname[0]},
			{"port", splitHostname[1]},
		}...)
	} else {
		log.Error("Unexpected hostname format '%s', skipping hostname and port inventory", *nodeResponse.Hostname)
	}

	for _, inventoryItem := range inventoryItems {
		if err := nodeEntity.SetInventoryItem("config/"+inventoryItem.key, "value", inventoryItem.value); err != nil {
			log.Error("Could not set inventory item '%s' on node entity: %v", inventoryItem.key, err)
		}
	}
}
