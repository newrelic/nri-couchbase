package definition

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
	MaxBucketCount         *int                    `json:"maxBucketCount" metric_name:"cluster.maximumBucketCount" source_type:"gauge"`
	ClusterName            *string                 `json:"clusterName"`
}

// AutoCompactionSettings struct for pools/default endpoint, autoCompactionSettings object
type AutoCompactionSettings struct {
	DatabaseFragmentationThreshold *DatabaseFragmentationThreshold `json:"databaseFragmentationThreshold"`
	IndexFragmentationThreshold    *IndexFragmentationThreshold    `json:"indexFragmentationThreshold"`
	ViewFragmentationThreshold     *ViewFragmentationThreshold     `json:"viewFragmentationThreshold"`
}

// DatabaseFragmentationThreshold struct for pools/default endpoint, autoCompactionSettings/databaseFragmentationThreshold object
type DatabaseFragmentationThreshold struct {
	Percentage *int `json:"percentage" metric_name:"cluster.databaseFragmentationThreshold" source_type:"gauge"`
}

// IndexFragmentationThreshold struct for pools/default endpoint, autoCompactionSettings/indexFragmentationThreshold object
type IndexFragmentationThreshold struct {
	Percentage *int `json:"percentage" metric_name:"cluster.indexFragmentationThreshold" source_type:"gauge"`
}

// ViewFragmentationThreshold struct for pools/default endpoint, autoCompactionSettings/viewFragmentationThreshold object
type ViewFragmentationThreshold struct {
	Percentage *int `json:"percentage" metric_name:"cluster.viewFragmentationThreshold" source_type:"gauge"`
}

// StorageTotals struct for pools/default endpoint, storageTotals object
type StorageTotals struct {
	HDD *StorageTotalsHDD `json:"hdd"`
	RAM *StorageTotalsRAM `json:"ram"`
}

// StorageTotalsRAM struct for pools/default endpoint, storageTotals/ram object
type StorageTotalsRAM struct {
	Total             *int64 `json:"total" metric_name:"cluster.memoryTotalInBytes" source_type:"gauge"`
	QuotaTotal        *int64 `json:"quotaTotal" metric_name:"cluster.memoryQuotaTotalInBytes" source_type:"gauge"`
	QuotaUsed         *int64 `json:"quotaUsed" metric_name:"cluster.memoryQuotaUsedInBytes" source_type:"gauge"`
	Used              *int64 `json:"used" metric_name:"cluster.memoryUsedInBytes" source_type:"gauge"`
	UsedByData        *int64 `json:"usedByData" metric_name:"cluster.memoryUsedByDataInBytes" source_type:"gauge"`
	QuotaUsedPerNode  *int64 `json:"quotaUsedPerNode" metric_name:"cluster.memoryQuotaUsedPerNodeInBytes" source_type:"gauge"`
	QuotaTotalPerNode *int64 `json:"quotaTotalPerNode" metric_name:"cluster.memoryQuotaTotalPerNodeInBytes" source_type:"gauge"`
}

// StorageTotalsHDD struct for pools/default endpoint, storageTotals/hdd object
type StorageTotalsHDD struct {
	Total      *int64 `json:"total" metric_name:"cluster.diskTotalInBytes" source_type:"gauge"`
	QuotaTotal *int64 `json:"quotaTotal" metric_name:"cluster.diskQuotaTotalInBytes" source_type:"gauge"`
	Used       *int64 `json:"used" metric_name:"cluster.diskUsedInBytes" source_type:"gauge"`
	UsedByData *int64 `json:"usedByData" metric_name:"cluster.diskUsedByDataInBytes" source_type:"gauge"`
	Free       *int64 `json:"free" metric_name:"cluster.diskFreeInBytes" source_type:"gauge"`
}

// Node struct for pools/default endpoint, nodes objects
type Node struct {
	SystemStats       *SystemStats `json:"systemStats"`
	RecoveryType      *string      `json:"recoveryType" metric_name:"node.recoveryType" source_type:"attribute"`
	Status            *string      `json:"status" metric_name:"node.status" source_type:"attribute"`
	Uptime            *string      `json:"uptime" metric_name:"node.uptime" source_type:"gauge"`
	Services          *[]string    `json:"services"`
	ClusterMembership *string      `json:"clusterMembership"`
	Hostname          *string      `json:"hostname"`
	OS                *string      `json:"os"`
	Version           *string      `json:"version"`
}

// SystemStats struct for pools/default endpoint, nodes/systemStats objects
type SystemStats struct {
	CPUUtilization *float64 `json:"cpu_utilization_rate" metric_name:"node.cpuUtilization" source_type:"gauge"`
	SwapTotal      *int64   `json:"swap_total" metric_name:"node.swapTotalInBytes" source_type:"gauge"`
	SwapUsed       *int64   `json:"swap_used" metric_name:"node.swapUsedInBytes" source_type:"gauge"`
	MemoryTotal    *int64   `json:"mem_total" metric_name:"node.memoryTotalInBytes" source_type:"gauge"`
	MemoryFree     *int64   `json:"mem_free" metric_name:"node.memoryFreeInBytes" source_type:"gauge"`
}

// =========

// PoolsDefaultBucket struct for pools/default/buckets endpoint
type PoolsDefaultBucket struct {
	BucketName     *string     `json:"name"`
	BasicStats     *BasicStats `json:"basicStats"`
	EvictionPolicy *string     `json:"evictionPolicy" metric_name:"bucket.evictionPolicy" source_type:"attribute"`
	NodeLocator    *string     `json:"nodeLocator" metric_name:"bucket.nodeLocator" source_type:"attribute"`
	ReplicaIndex   *bool       `json:"replicaIndex" metric_name:"bucket.replicaIndex" source_type:"gauge"`
	ReplicaNumber  *int        `json:"replicaNumber" metric_name:"bucket.replicaNumber" source_type:"gauge"`
	ThreadsNumber  *int        `json:"threadsNumber" metric_name:"bucket.threadsNumber" source_type:"gauge"`
	ProxyPort      *int        `json:"proxyPort"`
	BucketType     *string     `json:"bucketType"`
	UUID           *string     `json:"uuid"`
}

// BasicStats struct for pools/default/buckets endpoint, basicStats objects
type BasicStats struct {
	QuotaPercentUsed *float64 `json:"quotaPercentUsed" metric_name:"bucket.quotaUtilization" source_type:"gauge"`
	OpsPerSec        *int     `json:"opsPerSec" metric_name:"bucket.totalOperationsPerSecond" source_type:"gauge"`
	DiskFetches      *int     `json:"diskFetches" metric_name:"bucket.diskFetchesPerSecond" source_type:"gauge"`
	ItemCount        *int     `json:"itemCount" metric_name:"bucket.itemCount" source_type:"gauge"`
	DiskUsed         *int64   `json:"diskUsed" metric_name:"bucket.diskUsedInBytes" source_type:"gauge"`
	DataUsed         *int64   `json:"dataUsed" metric_name:"bucket.dataUsedInBytes" source_type:"gauge"`
	MemUsed          *int64   `json:"memUsed" metric_name:"bucket.memoryUsedInBytes" source_type:"gauge"`
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
	HitRatio                    *[]float64 `json:"hit_ratio" metric_name:"bucket.hitRatio" source_type:"gauge"`
	EpCacheMissRate             *[]float64 `json:"ep_cache_miss_rate" metric_name:"bucket.cacheMissRatio" source_type:"gauge"`
	EpResidentItemsRate         *[]float64 `json:"ep_resident_items_rate" metric_name:"bucket.residentItemsRatio" source_type:"gauge"`
	VbActiveResidentItemsRatio  *[]float64 `json:"vb_active_resident_items_ratio" metric_name:"bucket.activeResidentItemsRatio" source_type:"gauge"`
	VbReplicaResidentItemsRatio *[]float64 `json:"vb_replica_resident_items_ratio" metric_name:"bucket.replicaResidentItemsRatio" source_type:"gauge"`
	VbPendingResidentItemsRatio *[]float64 `json:"vb_pending_resident_items_ratio" metric_name:"bucket.pendingResidentItemsRatio" source_type:"gauge"`
	AvgDiskUpdateTime           *[]float64 `json:"avg_disk_update_time" metric_name:"bucket.averageDiskUpdateTimeInMilliseconds" source_type:"gauge"`
	AvgDiskCommitTime           *[]float64 `json:"avg_disk_commit_time" metric_name:"bucket.averageDiskCommitTimeInMilliseconds" source_type:"gauge"`
	BytesRead                   *[]float64 `json:"bytes_read" metric_name:"bucket.readRatePerSecond" source_type:"rate"`
	BytesWritten                *[]float64 `json:"bytes_written" metric_name:"bucket.writeRatePerSecond" source_type:"rate"`
	CmdGet                      *[]float64 `json:"cmd_get" metric_name:"bucket.readOperationsPerSecond" source_type:"rate"`
	CmdSet                      *[]float64 `json:"cmd_set" metric_name:"bucket.writeOperationsPerSecond" source_type:"rate"`
	CurrConnections             *[]float64 `json:"curr_connections" metric_name:"bucket.currentConnections" source_type:"gauge"`
	DecrHits                    *[]float64 `json:"decr_hits" metric_name:"bucket.decrementHitsPerSecond" source_type:"rate"`
	DecrMisses                  *[]float64 `json:"decr_misses" metric_name:"bucket.decrementMissesPerSecond" source_type:"rate"`
	DeleteHits                  *[]float64 `json:"delete_hits" metric_name:"bucket.deleteHitsPerSecond" source_type:"rate"`
	DeleteMisses                *[]float64 `json:"delete_misses" metric_name:"bucket.deleteMissesPerSecond" source_type:"rate"`
	DiskWriteQueue              *[]float64 `json:"disk_write_queue" metric_name:"bucket.diskWriteQueue" source_type:"gauge"`
	EpMemHighWat                *[]float64 `json:"ep_mem_high_wat" metric_name:"bucket.memoryHighWaterMarkInBytes" source_type:"gauge"`
	EpMemLowWat                 *[]float64 `json:"ep_mem_low_wat" metric_name:"bucket.memoryLowWaterMarkInBytes" source_type:"gauge"`
	EpMetaDataMemory            *[]float64 `json:"ep_meta_data_memory" metric_name:"bucket.metadataInRAMInBytes" source_type:"gauge"`
	EpNumValueEjects            *[]float64 `json:"ep_num_value_ejects" metric_name:"bucket.ejectionsPerSecond" source_type:"rate"`
	EpOomErrors                 *[]float64 `json:"ep_oom_errors" metric_name:"bucket.outOfMemoryErrorsPerSecond" source_type:"rate"`
	EpOpsCreate                 *[]float64 `json:"ep_ops_create" metric_name:"bucket.diskCreateOperationsPerSecond" source_type:"rate"`
	EpOpsUpdate                 *[]float64 `json:"ep_ops_update" metric_name:"bucket.diskUpdateOperationsPerSecond" source_type:"rate"`
	EpOverhead                  *[]float64 `json:"ep_overhead" metric_name:"bucket.overheadInBytes" source_type:"gauge"`
	EpTmpOomErrors              *[]float64 `json:"ep_tmp_oom_errors" metric_name:"bucket.temporaryOutOfMemoryErrorsPerSecond" source_type:"rate"`
	Evictions                   *[]float64 `json:"evictions" metric_name:"bucket.evictionsPerSecond" source_type:"rate"`
	GetHits                     *[]float64 `json:"get_hits" metric_name:"bucket.getHitsPerSecond" source_type:"rate"`
	GetMisses                   *[]float64 `json:"get_misses" metric_name:"bucket.getMissesPerSecond" source_type:"rate"`
	IncrHits                    *[]float64 `json:"incr_hits" metric_name:"bucket.incrementHitsPerSecond" source_type:"rate"`
	IncrMisses                  *[]float64 `json:"incr_misses" metric_name:"bucket.incrementMissesPerSecond" source_type:"rate"`
	Misses                      *[]float64 `json:"misses" metric_name:"bucket.missesPerSecond" source_type:"rate"`
}

// =========

// AutoFailover struct for settings/autoFailover endpoint
type AutoFailover struct {
	Count   *int  `json:"count" metric_name:"cluster.autoFailoverCount" source_type:"gauge"`
	Enabled *bool `json:"enabled" metric_name:"cluster.autoFailoverEnabled" source_type:"gauge"`
}

// =========

// AdminVitals struct for admin/vitals endpoint
type AdminVitals struct {
	Uptime                  *string  `json:"uptime"` // requires postprocessing
	Version                 *string  `json:"version"`
	TotalThreads            *int     `json:"total.threads" metric_name:"queryengine.totalThreads" source_type:"gauge"`
	Cores                   *int     `json:"cores" metric_name:"queryengine.cores" source_type:"gauge"`
	GCNum                   *int     `json:"gc.num" metric_name:"queryengine.garbageCollectionNumber" source_type:"gauge"`
	GCPauseTime             *string  `json:"gc.pause.time"` // requires postprocessing
	GCPausePercent          *int     `json:"gc.pause.percent" metric_name:"queryengine.garbageCollectionPaused" source_type:"gauge"`
	MemoryUsage             *int64   `json:"memory.usage" metric_name:"queryengine.usedMemoryInBytes" source_type:"gauge"`
	MemoryTotal             *int64   `json:"memory.total" metric_name:"queryengine.totalMemoryInBytes" source_type:"gauge"`
	MemorySystem            *int64   `json:"memory.system" metric_name:"queryengine.systemMemoryInBytes" source_type:"gauge"`
	CPUUserPercent          *float64 `json:"cpu.user.percent" metric_name:"queryengine.userCPUUtilization" source_type:"gauge"`
	CPUSystemPercent        *float64 `json:"cpu.sys.percent" metric_name:"queryengine.systemCPUUtilization" source_type:"gauge"`
	RequestCompletedCount   *int     `json:"request.completed.count" metric_name:"queryengine.completedRequests" source_type:"gauge"`
	RequestActiveCount      *int     `json:"request.active.count" metric_name:"queryengine.activeRequests" source_type:"gauge"`
	RequestPerSec1Min       *float64 `json:"request.per.sec.1min" metric_name:"queryengine.requestsLast1MinutesPerSecond" source_type:"rate"`
	RequestPerSec5Min       *float64 `json:"request.per.sec.5min" metric_name:"queryengine.requestsLast5MinutesPerSecond" source_type:"rate"`
	RequestPerSec15Min      *float64 `json:"request.per.sec.15min" metric_name:"queryengine.requestsLast15MinutesPerSecond" source_type:"rate"`
	RequestTimeMean         *string  `json:"request_time.mean"`         // requires postprocessing
	RequestTimeMedian       *string  `json:"request_time.median"`       // requires postprocessing
	RequestTime80Percentile *string  `json:"request_time.80percentile"` // requires postprocessing
	RequestTime95Percentile *string  `json:"request_time.95percentile"` // requires postprocessing
	RequestTime99Percentile *string  `json:"request_time.99percentile"` // requires postprocessing
	RequestPreparedPercent  *float64 `json:"request.prepared.percent" metric_name:"queryengine.preparedStatementUtilization" source_type:"gauge"`
}

// =========

// AdminSettings struct for admin/settings endpoint
type AdminSettings struct {
	CompletedLimit     *int `json:"completed-limit" metric_name:"queryengine.completedLimit" source_type:"gauge"`
	CompletedThreshold *int `json:"completed-threshold" metric_name:"queryengine.completedThresholdInMilliseconds" source_type:"gauge"`
}
