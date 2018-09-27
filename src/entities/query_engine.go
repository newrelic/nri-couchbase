package entities

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-couchbase/src/client"
	"github.com/newrelic/nri-couchbase/src/definition"
)

type queryEngineCollector struct {
	defaultCollector
}

func (qe *queryEngineCollector) GetEntity() (*integration.Entity, error) {
	return qe.GetIntegration().Entity(qe.GetName(), "queryEngine")
}

func (qe *queryEngineCollector) Collect(collectInventory, collectMetrics bool) error {
	queryEngineEntity, err := qe.GetEntity()
	if err != nil {
		return err
	}

	qeClient := qe.GetClient()
	settings, vitals, err := getQueryEngineResponses(qeClient)

	if collectMetrics {
		collectQueryEngineMetrics(queryEngineEntity, settings, vitals)
	}

	if collectInventory {
		collectQueryEngineInventory(queryEngineEntity, vitals, qeClient.Hostname, qeClient.QueryPort)
	}

	return nil
}

func collectQueryEngineMetrics(queryEngineEntity *integration.Entity, settingsResponse *definition.AdminSettings, vitalsResponse *definition.AdminVitals) {
	queryEngineMetricSet := queryEngineEntity.NewMetricSet("CouchbaseQueryEngineSample",
		metric.Attribute{Key: "displayName", Value: queryEngineEntity.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: queryEngineEntity.Metadata.Namespace + ":" + queryEngineEntity.Metadata.Name},
	)

	for _, response := range []interface{}{settingsResponse, vitalsResponse} {
		err := queryEngineMetricSet.MarshalMetrics(response)
		if err != nil {
			log.Error("Could not marshal query engine metrics from : %v", err)
		}
	}

	// metrics that need time conversion
	for _, metric := range []struct {
		metricName  string
		sourceType  metric.SourceType
		metricValue string
	}{
		{"queryengine.averageRequestTimeInMilliseconds", metric.GAUGE, *vitalsResponse.RequestTimeMean},
		{"queryengine.garbageCollectionTimePausedInMilliseconds", metric.GAUGE, *vitalsResponse.GCPauseTime},
		{"queryengine.medianRequestTimeInMilliseconds", metric.GAUGE, *vitalsResponse.RequestTimeMedian},
		{"queryengine.requestTime80thPercentileInMilliseconds", metric.GAUGE, *vitalsResponse.RequestTime80Percentile},
		{"queryengine.requestTime95thPercentileInMilliseconds", metric.GAUGE, *vitalsResponse.RequestTime95Percentile},
		{"queryengine.requestTime99thPercentileInMilliseconds", metric.GAUGE, *vitalsResponse.RequestTime99Percentile},
		{"queryengine.uptimeInMilliseconds", metric.GAUGE, *vitalsResponse.Uptime},
	} {
		convertedValue, err := convertTimeUnits(metric.metricValue)
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

func collectQueryEngineInventory(queryEngineEntity *integration.Entity, vitalsResponse *definition.AdminVitals, hostname string, port int) {
	inventoryItems := []inventoryItem{
		{"version", vitalsResponse.Version},
		{"hostname", hostname},
		{"port", port},
	}

	for _, inventoryItem := range inventoryItems {
		if err := queryEngineEntity.SetInventoryItem("config/"+inventoryItem.key, "value", inventoryItem.value); err != nil {
			log.Error("Could not set inventory item '%s' on query engine entity: %v", inventoryItem.key, err)
		}
	}
}

func getQueryEngineResponses(client *client.HTTPClient) (settings *definition.AdminSettings, vitals *definition.AdminVitals, err error) {
	settings = new(definition.AdminSettings)
	vitals = new(definition.AdminVitals)

	err = client.Request("/admin/settings", &settings)
	if err != nil {
		return
	}

	err = client.Request("/admin/vitals", &vitals)
	return
}

func convertTimeUnits(time string) (float64, error) {
	// go's regexp package does not support lookaround,
	// which would have been a lot cleaner.
	timeRegex, err := regexp.Compile("[\\d\\.]+[a-zµ]+")
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
	timeRegex, err := regexp.Compile("[\\d\\.]+")
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