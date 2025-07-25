package ingest_queue_plugin

import (
	"fmt"
	logwarp_common "github.com/Sage-Infrastructure-Solutions-Group-Inc/Logwarp-2-Common"
	"log/slog"
)

type DummyQueuePlugin struct {
	logGroup []any
}

func (d *DummyQueuePlugin) Config(config map[string]interface{}) error {
	d.logGroup = append(d.logGroup, slog.Group("source", "queue_plugin"))
	d.logGroup = append(d.logGroup, slog.Group("plugin_name", "Dummy Queue Plugin"))
	slog.Debug(fmt.Sprintf("Dummy Queue Plugin Config: %v", config), d.logGroup...)
	return nil
}

func (d *DummyQueuePlugin) Enqueue(batch logwarp_common.Batch) error {
	slog.Debug(fmt.Sprintf("Enqueue Batch: %v", batch), d.logGroup...)
	return nil
}

func (d *DummyQueuePlugin) Test() bool {
	return true
}

var interfaceTest QueuePlugin = &DummyQueuePlugin{}

func NewQueue() QueuePlugin {
	return &DummyQueuePlugin{}
}
