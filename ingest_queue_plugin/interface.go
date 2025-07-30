package ingest_queue_plugin

import (
	"github.com/Sage-Infrastructure-Solutions-Group-Inc/Logwarp-2-Common/protobuf"
)

type QueuePlugin interface {
	Config(config map[string]interface{}) error
	Enqueue(records protobuf.RecordList, inputPluginName string) error
	Test() bool
}
