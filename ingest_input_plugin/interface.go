package ingest_input_plugin

import "logwarp-ingest/config"

// InputPlugin defines the common interface that all input plugins must implement.
type InputPlugin interface {
	Start(config map[string]interface{}, queues []config.QueueConfig) error // Start listening on configured ports or connections
	Stop()                                                                  // Stop listening and clean up resources
	Flush() error                                                           // Flush the current buffer in the input plugin
}
