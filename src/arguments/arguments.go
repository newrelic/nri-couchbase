package arguments

import (
	"errors"
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/args"
)

// ArgumentList struct to hold all arguments needed to connect to a Couchbase environment
type ArgumentList struct {
	args.DefaultArgumentList
	Hostname              string `default:"localhost" help:"The hostname or IP of the Couchbase node being monitored"`
	Port                  int    `default:"8091" help:"The port used to connect to the Couchbase API"`
	QueryPort             int    `default:"8093" help:"The port used to connect to the N1QL service"`
	Username              string `default:"" help:"The username used to connect to the Couchbase API"`
	Password              string `default:"" help:"The password used to connect to the Couchbase API"`
	UseSSL                bool   `default:"false" help:"Signals whether to use SSL or not. Certificate bundle must be supplied"`
	CABundleFile          string `default:"" help:"Alternative Certificate Authority bundle file"`
	CABundleDir           string `default:"" help:"Alternative Certificate Authority bundle directory"`
	EnableClusterAndNodes bool   `default:"true" help:"If true, collects cluster and node resources"`
	EnableBuckets         bool   `default:"true" help:"If true, collects bucket resources"`
	EnableBucketStats     bool   `default:"true" help:"If true, collects additional bucket statistics"`
	Timeout               int    `default:"30" help:"Timeout for an API call"`
}

// Validate validates an argument list and returns an error if something is wrong
func (args *ArgumentList) Validate() error {
	if args.Username == "" {
		return errors.New("must provide a username argument")
	}

	if args.Password == "" {
		return errors.New("must provide a password argument")
	}

	if args.Hostname == "" {
		return errors.New("must provide a host argument")
	}

	if args.Port < 0 || args.Port > 65535 {
		return fmt.Errorf("invalid port %v", args.Port)
	}

	return nil
}