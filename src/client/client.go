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
)

// HTTPClient represents a single connection to an Elasticsearch host
type HTTPClient struct {
	BaseURL      string
	BaseQueryURL string
	Username     string
	Password     string
	Client       *http.Client
	Hostname     string
	Port         int
	QueryPort    int
}

// CreateClient creates a new http client for Couchbase.
// The hostnameOverride parameter specifies a hostname that the client should connect to.
// Passing in an empty string causes the client to use the hostname specified in the command-line args. (default behavior)
func CreateClient(args *arguments.ArgumentList, hostnameOverride string) (*HTTPClient, error) {
	hostname := func() string {
		if hostnameOverride != "" {
			return hostnameOverride
		}
		return args.Hostname
	}()

	options := []nrHttp.ClientOption{
		nrHttp.WithTimeout(time.Duration(args.Timeout) * time.Second),
	}
	if args.CABundleFile != "" {
		options = append(options, nrHttp.WithCABundleFile(args.CABundleFile))
	}

	if args.CABundleDir != "" {
		options = append(options, nrHttp.WithCABundleDir(args.CABundleDir))
	}

	httpClient, err := nrHttp.New(options...)
	if err != nil {
		return nil, err
	}

	return &HTTPClient{
		Client:       httpClient,
		Username:     args.Username,
		Password:     args.Password,
		BaseURL:      getBaseURL(args.UseSSL, hostname, args.Port),
		BaseQueryURL: getBaseURL(args.UseSSL, hostname, args.QueryPort),
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

// Request attempts to make a request to the Couchbase API, storing the result in the given model if successful.
// Returns an error if the request cannot be completed or a non-200 status code is returned.
func (c *HTTPClient) Request(endpoint string, model interface{}) error {
	// make sure we point to the query engine port for query engine metrics
	url := c.BaseURL + endpoint
	if strings.HasPrefix(endpoint, "/admin/") {
		url = c.BaseQueryURL + endpoint
	}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("could not create request for endpoint '%s': %v", endpoint, err)
	}
	request.SetBasicAuth(c.Username, c.Password)

	response, err := c.Client.Do(request)
	if err != nil {
		return fmt.Errorf("could not complete request for endpoint '%s': %v", endpoint, err)
	}

	err = checkStatusCode(response)
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

func checkStatusCode(response *http.Response) error {
	if response.StatusCode != 200 {
		return fmt.Errorf("received non-200 status code '%v'", response.StatusCode)
	}
	return nil
}
