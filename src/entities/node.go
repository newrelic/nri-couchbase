package entities

import (
	"github.com/newrelic/nri-couchbase/src/definition"
	"github.com/newrelic/infra-integrations-sdk/integration"
)

type nodeCollector struct {
	defaultCollector
	nodeDetails *definition.Node
}

func (n *nodeCollector) GetEntity() (*integration.Entity, error) {
	e, err := n.GetIntegration().Entity(n.GetName(), "node")
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (n *nodeCollector) Collect(collectInventory, collectMetrics bool) error {
	return nil
}