# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## 2.6.0 (2023-06-06)
# Changed
- Update Go version to 1.20

## 2.5.2 (2022-06-20)
### Changed
- Update Go version to 1.18
- Bump dependencies
### Added
Added support for more distributions:
- RHEL(EL) 9
- Ubuntu 22.04

## 2.5.1 (2021-10-20)
### Added
Added support for more distributions:
- Debian 11
- Ubuntu 20.10
- Ubuntu 21.04
- SUSE 12.15
- SUSE 15.1
- SUSE 15.2
- SUSE 15.3
- Oracle Linux 7
- Oracle Linux 8

## 2.5.0 (2021-08-30)
### Added

Moved default config.sample to [V4](https://docs.newrelic.com/docs/create-integrations/infrastructure-integrations-sdk/specifications/host-integrations-newer-configuration-format/), added a dependency for infra-agent version 1.20.0

Please notice that old [V3](https://docs.newrelic.com/docs/create-integrations/infrastructure-integrations-sdk/specifications/host-integrations-standard-configuration-format/) configuration format is deprecated, but still supported.

## 2.4.1 (2021-06-10)
### Changed
- Add ARM support.

## 2.4.0 (2021-05-05)
### Changed
- Update Go to v1.16.
- Migrate to Go Modules
- Update Infrastracture SDK to v3.6.7.
- Update other dependecies.
## 2.3.8 (2021-03-24)
### Fixed
- Adds arm packages and binaries

## 2.3.7 (2020-07-07)
### Fixed
- Crash when additional fields may come back as floats

## 2.3.6 (2020-03-25)
### Fixed
- Crash when API returns float values for fields that look like integers

## 2.3.5 (2020-02-03)
### Fixed
- Crash when API returns string values for fields that are normally integers

## 2.3.4 (2020-01-27)
### Fixed
- Added safer defaults for ClusterName
### Added
- Collect node's `interestingMetrics`

## 2.2.2 (2020-01-17)
### Fixed
-  Issue causing duplicate metrics to show up for buckets

## 2.2.1 (2020-01-14)
### Added
-  A number of new metrics.

## 2.2.0 (2020-01-13)
### Changed
- BROKEN

## 2.1.0 (2019-11-18)
### Changed
- Renamed the integration executable from nr-couchbase to nri-couchbase in order to be consistent with the package naming. **Important Note:** if you have any security module rules (eg. SELinux), alerts or automation that depends on the name of this binary, these will have to be updated.

## 2.0.4 - 2019-10-23
### Fixed
- Unique GUIDS for windows components

## 2.0.3 - 2019-10-23
### Added
- Windows installer packaging
-
## 2.0.2 - 2019-07-26
### Added
- provided sidecar Kubernetes containers

## 2.0.0 - 2019-04-18
### Changes
- Renamed entity namespaces to scope them to couchbase
- Updated to v3 SDK
- Added clusterName as an identity attribute
### Fixes
- Bug where integration would collect all nodes as the same node

## 1.0.2 - 2019-03-19
### Fixes
- Remove unused dependency for nri-jmx

## 1.0.1 - 2019-02-04
### Fixes
- `EnableClusterAndNodes` and `EnableBuckets` actually worked now
- Fixed issues where unchecked pointer dereferences could happen
### Changes
- Bumped to protocol version 2 in definition file

## 1.0.0 - 2018-11-29
### Changes
- Bumped version for GA release

## 0.1.1 - 2018-11-15
### Added
- Cluster name and hostname to Node and Query Engine entities

## 0.1.0 - 2018-09-14
### Added
- Initial version: Includes Metrics and Inventory data
