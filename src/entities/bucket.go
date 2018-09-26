package entities

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-couchbase/src/definition"
)

type bucketCollector struct {
	defaultCollector
	bucketResponse         *definition.PoolsDefaultBucket
	collectExtendedMetrics bool
}

func (b *bucketCollector) GetEntity() (*integration.Entity, error) {
	return b.GetIntegration().Entity(b.GetName(), "bucket")
}

func (b *bucketCollector) Collect(collectInventory, collectMetrics bool) error {
	bucketEntity, err := b.GetEntity()
	if err != nil {
		return err
	}

	if collectMetrics {
		bucketMetricSet := collectBucketMetrics(bucketEntity, b.bucketResponse)

		if b.collectExtendedMetrics {
			var bucketStats definition.BucketStats
			endpoint := fmt.Sprintf("/pools/default/buckets/%s/stats", b.GetName())
			err := b.GetClient().Request(endpoint, &bucketStats)
			if err != nil {
				return err
			}

			collectExtendedBucketMetrics(bucketMetricSet, &bucketStats, b.GetName())
		}
	}

	if collectInventory {
		collectBucketInventory(bucketEntity, b.bucketResponse)
	}

	return nil
}

func collectBucketMetrics(bucketEntity *integration.Entity, baseBucketResponse *definition.PoolsDefaultBucket) *metric.Set {
	bucketMetricSet := bucketEntity.NewMetricSet("CouchbaseBucketSample",
		metric.Attribute{Key: "displayName", Value: bucketEntity.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: bucketEntity.Metadata.Namespace + ":" + bucketEntity.Metadata.Name},
	)

	err := bucketMetricSet.MarshalMetrics(baseBucketResponse)
	if err != nil {
		log.Error("Could not marshal metrics from bucket struct: %v", err)
	}

	return bucketMetricSet
}

func collectExtendedBucketMetrics(metricSet *metric.Set, bucketStats *definition.BucketStats, bucketName string) error {
	structType := reflect.TypeOf(*bucketStats.Op.Samples)
	structValue := reflect.ValueOf(*bucketStats.Op.Samples)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.FieldByName(field.Name)
		metricValue, err := getValueFromArray(fieldValue)
		if err != nil {
			log.Error("Could not get metric value for %s: %v", field.Name, err)
		}

		metricName := field.Tag.Get("metric_name")
		metricType := field.Tag.Get("source_type")
		sourceType := getSourceTypeFromTag(metricType)
		err = metricSet.SetMetric(metricName, metricValue, sourceType)
		if err != nil {
			log.Error("Could not set metric '%s': %v", metricName, err)
		}
	}
	return nil
}

func getValueFromArray(values reflect.Value) (float64, error) {
	array := reflect.Indirect(values)
	if array.Kind() != reflect.Slice {
		return 0, errors.New("value is not an array")
	}

	return array.Index(0).Float(), nil
}

func getSourceTypeFromTag(metricType string) metric.SourceType {
	switch metricType {
	case "gauge":
		return metric.GAUGE
	case "rate":
		return metric.RATE
	default:
		return metric.ATTRIBUTE
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
