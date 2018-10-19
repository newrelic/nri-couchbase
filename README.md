# New Relic Infrastructure Integration for Couchbase

The New Relic Infrastructure Integration for Couchbase captures critical performance metrics and inventory reported by Couchbase clusters. Data on the cluster, node, and bucket levels is collected.

All data is retrieved using Couchbase's REST API.

## Requirements

No additional requirements

## Installation

- download an archive file for the `Couchbase` Integration
- extract `couchbase-definition.yml` and `/bin` directory into `/var/db/newrelic-infra/newrelic-integrations`
- add execute permissions for the binary file `nr-couchbase` (if required)
- extract `couchbase-config.yml.sample` into `/etc/newrelic-infra/integrations.d`

## Usage

This is the description about how to run the Couchbase Integration with New Relic Infrastructure agent, so it is required to have the agent installed (see [agent installation](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/installation/install-infrastructure-linux)).

In order to use the Couchbase Integration it is required to configure `couchbase-config.yml.sample` file. Firstly, rename the file to `couchbase-config.yml`. Then, depending on your needs, specify all instances that you want to monitor. Once this is done, restart the Infrastructure agent.

You can view your data in Insights by creating your own custom NRQL queries. To do so use the **CouchbaseClusterSample**, **CouchbaseNodeSample**, **CouchbaseQueryEngineSample**, and **CouchbaseBucketSample** event types.

## Compatibility

* Supported OS: No limitations
* Couchbase versions: 4.0+

## Integration Development usage

Assuming that you have source code you can build and run the Couchbase Integration locally.

* Go to directory of the Couchbase Integration and build it
```bash
$ make
```
* The command above will execute tests for the Couchbase Integration and build an executable file called `nr-couchbase` in `bin` directory.
```bash
$ ./bin/nr-couchbase
```
* If you want to know more about usage of `./nr-couchbase` check
```bash
$ ./bin/nr-couchbase -help
```

For managing external dependencies [govendor tool](https://github.com/kardianos/govendor) is used. It is required to lock all external dependencies to specific version (if possible) into vendor directory.
