package ingest_input_plugin

import (
	"github.com/Sage-Infrastructure-Solutions-Group-Inc/Logwarp-2-Common/ingest_queue_plugin"
)

// InputPlugin defines the common interface that all input plugins must implement.
type InputPlugin interface {
	Configure(config map[string]interface{}, queues []ingest_queue_plugin.QueuePlugin) error // Configure the plugin
	Stop() error                                                                             // Stop listening and clean up resources
	Flush() error                                                                            // Flush the current buffer in the input plugin
	Run()                                                                                    // Run the plugin
	EnableTestMode()
	Cache() error
	ProcessCache() error
	SupportsLoadBalancedQueue() bool // Returns true if the Queue configuration type is supported
	SupportsFailOverQueue() bool     // Returns true if the Queue configuration type is supported
}
