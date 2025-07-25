package ingest_queue_plugin

import (
	logwarp_common "github.com/Sage-Infrastructure-Solutions-Group-Inc/Logwarp-2-Common"
)

type QueuePlugin interface {
	Config(config map[string]interface{}) error
	Enqueue(batch logwarp_common.Batch) error
	Test() bool
}
