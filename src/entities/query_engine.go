package entities

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/newrelic/infra-integrations-sdk/data/attribute"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-couchbase/src/client"
	"github.com/newrelic/nri-couchbase/src/definition"
)

type queryEngineCollector struct {
	defaultCollector
	clusterName string
}

func (qe *queryEngineCollector) GetEntity() (*integration.Entity, error) {
	clusterNameID := integration.IDAttribute{Key: "clusterName", Value: ClusterName}
	return qe.GetIntegration().Entity(qe.GetName(), "cb-queryEngine", clusterNameID)
}

func (qe *queryEngineCollector) Collect(collectInventory, collectMetrics bool) error {
	queryEngineEntity, err := qe.GetEntity()
	if err != nil {
		return err
	}

	qeClient := qe.GetClient()
	settings, vitals := getQueryEngineResponses(qeClient)

	if collectMetrics {
		collectQueryEngineMetrics(queryEngineEntity, settings, vitals, qe.clusterName, qeClient.Hostname)
	}

	if collectInventory && vitals != nil {
		collectQueryEngineInventory(queryEngineEntity, vitals, qeClient.Hostname, qeClient.QueryPort)
	}

	return nil
}

func collectQueryEngineMetrics(queryEngineEntity *integration.Entity, settingsResponse *definition.AdminSettings, vitalsResponse *definition.AdminVitals, clusterName, hostname string) {
	queryEngineMetricSet := queryEngineEntity.NewMetricSet("CouchbaseQueryEngineSample",
		attribute.Attribute{Key: "displayName", Value: queryEngineEntity.Metadata.Name},
		attribute.Attribute{Key: "entityName", Value: queryEngineEntity.Metadata.Namespace + ":" + queryEngineEntity.Metadata.Name},
		attribute.Attribute{Key: "cluster", Value: clusterName},
		attribute.Attribute{Key: "hostname", Value: hostname},
	)

	// marshal metrics from /admin/settings
	if settingsResponse != nil {
		tryMarshal(queryEngineMetricSet, settingsResponse)
	}

	// marshal metrics from /admin/vitals
	if vitalsResponse == nil {
		return
	}
	tryMarshal(queryEngineMetricSet, vitalsResponse)
	for _, metric := range []struct {
		metricName  string
		sourceType  metric.SourceType
		metricValue *string
	}{
		{"queryengine.averageRequestTimeInMilliseconds", metric.GAUGE, vitalsResponse.RequestTimeMean},
		{"queryengine.garbageCollectionTimePausedInMilliseconds", metric.GAUGE, vitalsResponse.GCPauseTime},
		{"queryengine.medianRequestTimeInMilliseconds", metric.GAUGE, vitalsResponse.RequestTimeMedian},
		{"queryengine.requestTime80thPercentileInMilliseconds", metric.GAUGE, vitalsResponse.RequestTime80Percentile},
		{"queryengine.requestTime95thPercentileInMilliseconds", metric.GAUGE, vitalsResponse.RequestTime95Percentile},
		{"queryengine.requestTime99thPercentileInMilliseconds", metric.GAUGE, vitalsResponse.RequestTime99Percentile},
		{"queryengine.uptimeInMilliseconds", metric.GAUGE, vitalsResponse.Uptime},
	} {
		if metric.metricValue == nil {
			log.Info("Metric '%s' was not returned in the JSON, skipping", metric.metricName)
			continue
		}
		convertedValue, err := convertTimeUnits(*metric.metricValue)
		if err != nil {
			log.Error("Could not convert time string '%s' to milliseconds: %v", metric.metricValue, err)
			continue
		}
		err = queryEngineMetricSet.SetMetric(metric.metricName, convertedValue, metric.sourceType)
		if err != nil {
			log.Error("Could not set time-converted metric '%s' on query engine: %v", metric.metricName, err)
		}
	}
}

func tryMarshal(metricSet *metric.Set, response interface{}) {
	err := metricSet.MarshalMetrics(response)
	if err != nil {
		log.Error("Could not marshal query engine metrics: %v", err)
	}
}

func collectQueryEngineInventory(queryEngineEntity *integration.Entity, vitalsResponse *definition.AdminVitals, hostname string, port int) {
	inventoryItems := []inventoryItem{
		{"version", vitalsResponse.Version},
		{"hostname", hostname},
		{"port", port},
	}

	for _, inventoryItem := range inventoryItems {
		if err := queryEngineEntity.SetInventoryItem(inventoryItem.key, "value", inventoryItem.value); err != nil {
			log.Error("Could not set inventory item '%s' on query engine entity: %v", inventoryItem.key, err)
		}
	}
}

func getQueryEngineResponses(client *client.HTTPClient) (*definition.AdminSettings, *definition.AdminVitals) {
	var settings *definition.AdminSettings
	var vitals *definition.AdminVitals

	err := client.Request("/admin/settings", &settings)
	if err != nil {
		log.Error("Couldn't retrieve query engine data from /admin/settings: %v", err)
	}

	err = client.Request("/admin/vitals", &vitals)
	if err != nil {
		log.Error("Couldn't retrieve query engine data from /admin/vitals: %v", err)
	}
	return settings, vitals
}

func convertTimeUnits(time string) (float64, error) {
	// go's regexp package does not support lookaround,
	// which would have been a lot cleaner.
	timeRegex, err := regexp.Compile(`[\d\.]+[a-zµ]+`)
	if err != nil {
		return 0, err
	}
	splitTime := timeRegex.FindAllString(time, -1)
	totalMilliseconds := 0.0
	for _, timeUnit := range splitTime {
		convertedUnit, err := convertTimeUnit(timeUnit)
		if err != nil {
			return 0, err
		}
		totalMilliseconds += convertedUnit
	}
	return totalMilliseconds, nil
}

func convertTimeUnit(time string) (float64, error) {
	timeRegex, err := regexp.Compile(`[\d\.]+`)
	if err != nil {
		return 0, err
	}
	timeUnitRegex, err := regexp.Compile("[a-zµ]+")
	if err != nil {
		return 0, err
	}

	timeValue, err := strconv.ParseFloat(timeRegex.FindString(time), 64)
	if err != nil {
		return 0, err
	}
	timeUnit := timeUnitRegex.FindString(time)

	return getMilliseconds(timeValue, timeUnit)
}

func getMilliseconds(timeValue float64, timeUnit string) (milliseconds float64, err error) {
	switch timeUnit {
	case "d":
		milliseconds = timeValue * 1000 * 60 * 60 * 24
	case "h":
		milliseconds = timeValue * 1000 * 60 * 60
	case "m":
		milliseconds = timeValue * 1000 * 60
	case "s":
		milliseconds = timeValue * 1000
	case "ms":
		milliseconds = timeValue
	case "us":
		milliseconds = timeValue / 1000
	case "µs":
		milliseconds = timeValue / 1000
	case "ns":
		milliseconds = timeValue / 1000 / 1000
	default:
		err = fmt.Errorf("unknown time unit '%s'", timeUnit)
	}

	return milliseconds, err
}
