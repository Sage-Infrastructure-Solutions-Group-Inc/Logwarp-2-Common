package ingest_input_plugin

import (
	logwarp_common "github.com/Sage-Infrastructure-Solutions-Group-Inc/Logwarp-2-Common"
)

// InputPlugin defines the common interface that all input plugins must implement.
type InputPlugin interface {
	Start(config map[string]interface{}, queues []logwarp_common.IngestQueueConfig) error // Start listening on configured ports or connections
	Stop()                                                                                // Stop listening and clean up resources
	Flush() error                                                                         // Flush the current buffer in the input plugin
}
