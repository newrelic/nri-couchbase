package main

// PoolsResponse struct for pools endpoint
type PoolsResponse struct {
	UUID    *string `json:"uuid"`
	Version *string `json:"implementationVersion"`
}

// =========

// PoolsDefaultResponse struct for pools/default endpoint
type PoolsDefaultResponse struct {
	AutoCompactionSettings *AutoCompactionSettings `json:"autoCompactionSettings"`
	StorageTotals          *StorageTotals          `json:"storageTotals"`
	Nodes                  *[]Node                 `json:"nodes"`
	MaxBucketCount         *int                    `json:"maxBucketCount" metric_name:"cluster.maximumBucketCount" metric_type:"gauge"`
}

// AutoCompactionSettings struct for pools/default endpoint, autoCompactionSettings object
type AutoCompactionSettings struct {
	DatabaseFragmentationThreshold *DatabaseFragmentationThreshold `json:"databaseFragmentationThreshold"`
	IndexFragmentationThreshold    *IndexFragmentationThreshold    `json:"indexFragmentationThreshold"`
	ViewFragmentationThreshold     *ViewFragmentationThreshold     `json:"viewFragmentationThreshold"`
}

// DatabaseFragmentationThreshold struct for pools/default endpoint, autoCompactionSettings/databaseFragmentationThreshold object
type DatabaseFragmentationThreshold struct {
	Percentage *int `json:"percentage" metric_name:"cluster.databaseFragmentationThreshold" metric_type:"gauge"`
}

// IndexFragmentationThreshold struct for pools/default endpoint, autoCompactionSettings/indexFragmentationThreshold object
type IndexFragmentationThreshold struct {
	Percentage *int `json:"percentage" metric_name:"cluster.indexFragmentationThreshold" metric_type:"gauge"`
}

// ViewFragmentationThreshold struct for pools/default endpoint, autoCompactionSettings/viewFragmentationThreshold object
type ViewFragmentationThreshold struct {
	Percentage *int `json:"percentage" metric_name:"cluster.viewFragmentationThreshold" metric_type:"gauge"`
}

// StorageTotals struct for pools/default endpoint, storageTotals object
type StorageTotals struct {
	HDD *StorageTotalsHDD `json:"hdd"`
	RAM *StorageTotalsRAM `json:"ram"`
}

// StorageTotalsRAM struct for pools/default endpoint, storageTotals/ram object
type StorageTotalsRAM struct {
	Total             *int64 `json:"total" metric_name:"cluster.memoryTotalInBytes" metric_type:"gauge"`
	QuotaTotal        *int64 `json:"quotaTotal" metric_name:"cluster.memoryQuotaTotalInBytes" metric_type:"gauge"`
	QuotaUsed         *int64 `json:"quotaUsed" metric_name:"cluster.memoryQuotaUsedInBytes" metric_type:"gauge"`
	Used              *int64 `json:"used" metric_name:"cluster.memoryUsedInBytes" metric_type:"gauge"`
	UsedByData        *int64 `json:"usedByData" metric_name:"cluster.memoryUsedByDataInBytes" metric_type:"gauge"`
	QuotaUsedPerNode  *int64 `json:"quotaUsedPerNode" metric_name:"cluster.memoryQuotaUsedPerNodeInBytes" metric_type:"gauge"`
	QuotaTotalPerNode *int64 `json:"quotaTotalPerNode" metric_name:"cluster.memoryQuotaTotalPerNodeInBytes" metric_type:"gauge"`
}

// StorageTotalsHDD struct for pools/default endpoint, storageTotals/hdd object
type StorageTotalsHDD struct {
	Total      *int64 `json:"total" metric_name:"cluster.diskTotalInBytes" metric_type:"gauge"`
	QuotaTotal *int64 `json:"quotaTotal" metric_name:"cluster.diskQuotaTotalInBytes" metric_type:"gauge"`
	Used       *int64 `json:"used" metric_name:"cluster.diskUsedInBytes" metric_type:"gauge"`
	UsedByData *int64 `json:"usedByData" metric_name:"cluster.diskUsedByDataInBytes" metric_type:"gauge"`
	Free       *int64 `json:"free" metric_name:"cluster.diskFreeInBytes" metric_type:"gauge"`
}

// Node struct for pools/default endpoint, nodes objects
type Node struct {
	SystemStats       *SystemStats `json:"systemStats"`
	RecoveryType      *string      `json:"recoveryType" metric_name:"node.recoveryType" metric_type:"attribute"`
	Services          *[]string    `json:"services"` // will require postprocessing
	Status            *string      `json:"status" metric_name:"node.status" metric_type:"attribute"`
	Uptime            *string      `json:"uptime" metric_name:"node.uptime" metric_type:"gauge"`
	ClusterMembership *string      `json:"clusterMembership"`
	Hostname          *string      `json:"hostname"`
	OS                *string      `json:"os"`
	Port              *int         // postprocess
	Version           *string      `json:"version"`
}

// SystemStats struct for pools/default endpoint, nodes/systemStats objects
type SystemStats struct {
	CPUUtilization *float64 `json:"cpu_utilization_rate" metric_name:"node.cpuUtilization" metric_type:"gauge"`
	SwapTotal      *int64   `json:"swap_total" metric_name:"node.swapTotalInBytes" metric_type:"gauge"`
	SwapUsed       *int64   `json:"swap_used" metric_name:"node.swapUsedInBytes" metric_type:"gauge"`
	MemoryTotal    *int64   `json:"mem_total" metric_name:"node.memoryTotalInBytes" metric_type:"gauge"`
	MemoryFree     *int64   `json:"mem_free" metric_name:"node.memoryUsedInBytes" metric_type:"gauge"`
}

// =========

// PoolsDefaultBucket struct for pools/default/buckets endpoint
type PoolsDefaultBucket struct {
	BucketName     *string     `json:"name"`
	BasicStats     *BasicStats `json:"basicStats"`
	EvictionPolicy *string     `json:"evictionPolicy" metric_name:"bucket.evictionPolicy" metric_type:"attribute"`
	NodeLocator    *string     `json:"nodeLocator" metric_name:"bucket.nodeLocator" metric_type:"attribute"`
	ReplicaIndex   *bool       `json:"replicaIndex" metric_name:"bucket.replicaIndex" metric_type:"attribute"`
	ReplicaNumber  *int        `json:"replicaNumber" metric_name:"bucket.replicaNumber" metric_type:"gauge"`
	ThreadsNumber  *int        `json:"threadsNumber" metric_name:"bucket.threadsNumber" metric_type:"gauge"`
	ProxyPort      *int        `json:"proxyPort"`
	BucketType     *string     `json:"bucketType"`
	UUID           *string     `json:"uuid"`
}

// BasicStats struct for pools/default/buckets endpoint, basicStats objects
type BasicStats struct {
	QuotaPercentUsed *float64 `json:"quotaPercentUsed" metric_name:"bucket.quotaUtilization" metric_type:"gauge"`
	OpsPerSec        *int     `json:"opsPerSec" metric_name:"bucket.totalOperationsPerSecond" metric_type:"gauge"`
	DiskFetches      *int     `json:"diskFetches" metric_name:"bucket.diskFetchesPerSecond" metric_type:"gauge"`
	ItemCount        *int     `json:"itemCount" metric_name:"bucket.itemCount" metric_type:"gauge"`
	DiskUsed         *int64   `json:"diskUsed" metric_name:"bucket.diskUsedInBytes" metric_type:"gauge"`
	DataUsed         *int64   `json:"dataUsed" metric_name:"bucket.dataUsedInBytes" metric_type:"gauge"`
	MemUsed          *int64   `json:"memUsed" metric_name:"bucket.memoryUsedInBytes" metric_type:"gauge"`
}

// =========

// BucketStats struct for pools/default/buckets/<bucket>/stats endpoint
type BucketStats struct {
	Op *OpStats `json:"op"`
}

// OpStats struct for op object
type OpStats struct {
	Samples *SampleStats `json:"samples"`
}

// SampleStats struct for op/samples object
type SampleStats struct {
	HitRatio                    *[]float64 `json:"hit_ratio"`
	EpCacheMissRate             *[]float64 `json:"ep_cache_miss_rate"`
	EpResidentItemsRate         *[]float64 `json:"ep_resident_items_rate"`
	VbActiveResidentItemsRatio  *[]float64 `json:"vb_active_resident_items_ratio"`
	VbReplicaResidentItemsRatio *[]float64 `json:"vb_replica_resident_items_ratio"`
	VbPendingResidentItemsRatio *[]float64 `json:"vb_pending_resident_items_ratio"`
	AvgDiskUpdateTime           *[]float64 `json:"avg_disk_update_time"`
	AvgDiskCommitTime           *[]float64 `json:"avg_disk_commit_time"`
	BytesRead                   *[]float64 `json:"bytes_read"`
	BytesWritten                *[]float64 `json:"bytes_written"`
	CmdGet                      *[]float64 `json:"cmd_get"`
	CmdSet                      *[]float64 `json:"cmd_set"`
	CurrConnections             *[]float64 `json:"curr_connections"`
	DecrHits                    *[]float64 `json:"decr_hits"`
	DecrMisses                  *[]float64 `json:"decr_misses"`
	DeleteHits                  *[]float64 `json:"delete_hits"`
	DeleteMisses                *[]float64 `json:"delete_misses"`
	DiskWriteQueue              *[]float64 `json:"disk_write_queue"`
	EpMemHighWat                *[]float64 `json:"ep_mem_high_wat"`
	EpMemLowWat                 *[]float64 `json:"ep_mem_low_wat"`
	EpMetaDataMemory            *[]float64 `json:"ep_meta_data_memory"`
	EpNumValueEjects            *[]float64 `json:"ep_num_value_ejects"`
	EpOomErrors                 *[]float64 `json:"ep_oom_errors"`
	EpOpsCreate                 *[]float64 `json:"ep_ops_create"`
	EpOpsUpdate                 *[]float64 `json:"ep_ops_update"`
	EpOverhead                  *[]float64 `json:"ep_overhead"`
	EpTmpOomErrors              *[]float64 `json:"ep_tmp_oom_errors"`
	Evictions                   *[]float64 `json:"evictions"`
	GetHits                     *[]float64 `json:"get_hits"`
	GetMisses                   *[]float64 `json:"get_misses"`
	IncrHits                    *[]float64 `json:"incr_hits"`
	IncrMisses                  *[]float64 `json:"incr_misses"`
	Misses                      *[]float64 `json:"misses"`
}

// =========

// AutoFailover struct for settings/autoFailover endpoint
type AutoFailover struct {
	Count *int `json:"count" metric_name:"cluster.autoFailoverCount" metric_type:"gauge"`
	Enabled *bool `json:"enabled" metric_name:"cluster.autoFailoverEnabled" metric_type:"attribute"`
}

// =========

// AdminVitals struct for admin/vitals endpoint
type AdminVitals struct {
	Uptime *string `json:"uptime"` // requires postprocessing
	Version *string `json:"version"`
	TotalThreads *int `json:"total.threads" metric_name:"queryengine.totalThreads" metric_type:"gauge"`
	Cores *int `json:"cores" metric_name:"queryengine.cores" metric_type:"gauge"`
	GCNum *int `json:"gc.num" metric_name:"queryengine.garbageCollectionNumber" metric_type:"attribute"`
	GCPauseTime *string `json:"gc.pause.time"` // requires postprocessing
	GCPausePercent *int `json:"gc.pause.percent" metric_name:"queryengine.garbageCollectionPaused" metric_type:"gauge"`
	MemoryUsage *int64 `json:"memory.usage" metric_name:"queryengine.usedMemoryInBytes" metric_type:"gauge"`
	MemoryTotal *int64 `json:"memory.total" metric_name:"queryengine.totalMemoryInBytes" metric_type:"gauge"`
	MemorySystem *int64 `json:"memory.system" metric_name:"queryengine.systemMemoryInBytes" metric_type:"gauge"`
	CPUUserPercent *float64 `json:"cpu.user.percent" metric_name:"queryengine.userCPUUtilization" metric_type:"gauge"`
	CPUSystemPercent *float64 `json:"cpu.system.percent" metric_name:"queryengine.systemCPUUtilization" metric_type:"gauge"`
	RequestCompletedCount *int `json:"request.completed.count" metric_name:"queryengine.completedRequests" metric_type:"gauge"`
	RequestActiveCount *int `json:"request.active.count" metric_name:"queryengine.activeRequests" metric_type:"gauge"`
	RequestPerSec1Min *float64 `json:"request.per.sec.1min" metric_name:"queryengine.requestsLast1MinutesPerSecond" metric_type:"rate"`
	RequestPerSec5Min *float64 `json:"request.per.sec.5min" metric_name:"queryengine.requestsLast5MinutesPerSecond" metric_type:"rate"`
	RequestPerSec15Min *float64 `json:"request.per.sec.15min" metric_name:"queryengine.requestsLast15MinutesPerSecond" metric_type:"rate"`
	RequestTimeMean *string `json:"request_time.mean"` // requires postprocessing
	RequestTimeMedian *string `json:"request_time.median"` // requires postprocessing
	RequestTime80Percentile *string `json:"request_time.80percentile"` // requires postprocessing
	RequestTime95Percentile *string `json:"request_time.95percentile"` // requires postprocessing
	RequestTime99Percentile *string `json:"request_time.99percentile"` // requires postprocessing
	RequestPreparedPercent *float64 `json:"request.prepared.percent" metric_name:"queryengine.preparedStatementUtilization" metric_type:"gauge"`
}

// =========

// AdminSettings struct for admin/settings endpoint
type AdminSettings struct {
	CompletedLimit *int `json:"completed-limit" metric_name:"queryengine.completedLimit" metric_type:"gauge"`
	CompletedThreshold *int `json:"completed-threshold" metric_name:"queryengine.completedThresholdInMilliseconds" metric_type:"gauge"`
}