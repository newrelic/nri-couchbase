package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	nrHttp "github.com/newrelic/infra-integrations-sdk/http"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-couchbase/src/arguments"
	"github.com/newrelic/nri-couchbase/src/definition"
)

// HTTPClient represents a single connection to an Elasticsearch host
type HTTPClient struct {
	baseURL      string
	baseQueryURL string
	username     string
	password     string
	client       *http.Client
	Hostname     string
	Port         int
	QueryPort    int
}

// CreateClient creates a new http client for Couchbase.
// The hostnameOverride parameter specifies a hostname that the client should connect to.
// Passing in an empty string causes the client to use the hostname specified in the command-line args. (default behavior)
func CreateClient(args *arguments.ArgumentList, hostnameOverride string) (*HTTPClient, error) {
	hostname := args.Hostname
	if hostnameOverride != "" {
		hostname = hostnameOverride
	}

	httpClient, err := nrHttp.New(args.CABundleFile, args.CABundleDir, time.Duration(args.Timeout)*time.Second)
	if err != nil {
		return nil, err
	}

	return &HTTPClient{
		client:       httpClient,
		username:     args.Username,
		password:     args.Password,
		baseURL:      getBaseURL(args.UseSSL, hostname, args.Port),
		baseQueryURL: getBaseURL(args.UseSSL, hostname, args.QueryPort),
		Hostname:     hostname,
		Port:         args.Port,
		QueryPort:    args.QueryPort,
	}, nil
}

func getBaseURL(useSSL bool, hostname string, port int) string {
	protocol := "http"
	if useSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s:%d", protocol, hostname, port)
}

// RequestAllBuckets takes a listing of bucket names (retrieved from the buckets endpoint) and returns each bucket's stats
func (c *HTTPClient) RequestAllBuckets(bucketList []string) map[string]*definition.BucketStats {
	var bucketStats = map[string]*definition.BucketStats{}
	for _, bucketName := range bucketList {
		bucketStatResponse := &definition.BucketStats{}
		endpoint := fmt.Sprintf("/pools/default/buckets/%s/stats", bucketName)
		err := c.Request(endpoint, bucketStatResponse)
		if err != nil {
			log.Error("Could not retrieve stats for bucket '%s': %v", bucketName, err)
			continue
		}
		bucketStats[bucketName] = bucketStatResponse
	}
	return bucketStats
}

// Request attempts to make a request to the Couchbase API, storing the result in the given model if successful.
// Returns an error if the request cannot be completed or a non-200 status code is returned.
func (c *HTTPClient) Request(endpoint string, model interface{}) error {
	// make sure we point to the query engine port for query engine metrics
	url := c.baseURL + endpoint
	if strings.HasPrefix(endpoint, "/admin/") {
		url = c.baseQueryURL + endpoint
	}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("could not create request for endpoint '%s': %v", endpoint, err)
	}
	request.SetBasicAuth(c.username, c.password)

	response, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("could not complete request for endpoint '%s': %v", endpoint, err)
	}

	err = c.checkStatusCode(response)
	if err != nil {
		return fmt.Errorf("could not complete request for endpoint '%s': %v", endpoint, err)
	}

	// decode json response
	err = json.NewDecoder(response.Body).Decode(model)
	if err != nil {
		return fmt.Errorf("received an unexpected response from endpoint '%s': %v", endpoint, err)
	}

	err = response.Body.Close()
	if err != nil {
		log.Error("Could not close response body: %v", err)
	}

	return nil
}

func (c *HTTPClient) checkStatusCode(response *http.Response) error {
	if response.StatusCode != 200 {
		return fmt.Errorf("received non-200 status code '%v'", response.StatusCode)
	}
	return nil
}
