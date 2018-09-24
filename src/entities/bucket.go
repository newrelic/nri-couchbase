package entities

import (
	"github.com/newrelic/nri-couchbase/src/definition"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
)

type bucketCollector struct {
	defaultCollector
	bucketResponse *definition.PoolsDefaultBucket
}

func (b *bucketCollector) GetEntity() (*integration.Entity, error) {
	entity, err := b.GetIntegration().Entity(b.GetName(), "bucket")
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (b *bucketCollector) Collect(collectInventory, collectMetrics bool) error {
	log.Info("I AM A BUCKET NAMED %s", b.GetName())

	bucketEntity, err := b.GetEntity()
	if err != nil {
		return err
	}
	
	if collectMetrics {
		collectBucketMetrics(bucketEntity, b.bucketResponse)
	}

	if collectInventory {
		collectBucketInventory(bucketEntity, b.bucketResponse)
	}
	
	return nil
}

func collectBucketMetrics(bucketEntity *integration.Entity, baseBucketResponse *definition.PoolsDefaultBucket) {
	bucketMetricSet := bucketEntity.NewMetricSet("CouchbaseBucketSample",
		metric.Attribute{Key: "displayName", Value: bucketEntity.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: bucketEntity.Metadata.Namespace + ":" + bucketEntity.Metadata.Name},
	)

	err := bucketMetricSet.MarshalMetrics(baseBucketResponse)
	if err != nil {
		log.Error("Could not marshal metrics from bucket struct: %v", err)
	}
}

func collectBucketInventory(bucketEntity *integration.Entity, baseBucketResponse *definition.PoolsDefaultBucket) {
	items := []inventoryItem{
		{"nodeLocator", baseBucketResponse.NodeLocator},
		{"proxyPort", baseBucketResponse.ProxyPort},
		{"bucketType", baseBucketResponse.BucketType},
		{"uuid", baseBucketResponse.UUID},
	}

	for _, item := range items {
		if err := bucketEntity.SetInventoryItem("config/"+item.key, "value", item.value); err != nil {
			log.Error("Could not set inventory item '%s' on bucket entity: %v", item.key, err)
		}
	}
}