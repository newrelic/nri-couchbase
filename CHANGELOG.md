# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## 2.2.0 (2020-01-13)
### Added
-  A number of new metrics.

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
