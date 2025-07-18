package ingest_queue_plugin

import "logwarp_common"

type QueuePlugin interface {
	Config(config map[string]interface{}) error
	Enqueue(batch logwarp_common_old.Batch) error
}
