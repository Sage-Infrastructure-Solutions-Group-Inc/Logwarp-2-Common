package ingest_input_plugin

import (
	"fmt"
	logwarp_common "github.com/Sage-Infrastructure-Solutions-Group-Inc/Logwarp-2-Common"
	"strings"
)

/*
The code that follows is in relation to the mapping of Queues to inputs in the configuration YAML file. There are a few
options when you map these queues back to your Inputs:
input:
  - ... (Truncated for Brevity)
    plugin-flush-timeout: 5 # In seconds
    queues:
      - e38c5275-61d1-11f0-8000-ede2451e963c  # Basic => The input will send to this Queue AND any others listed
      - load-balanced: ["e38c5275-61d1-11f0-8000-ede2451e963c", "e38c5275-61d1-11f0-8000-ede2451e963b"] # LoadBalanced the Input will load-balance between the defined queues AND send to any other Queues
      - fail-over: ["e38c5275-61d1-11f0-8000-ede2451e963c", "e38c5275-61d1-11f0-8000-ede2451e963b"]] # FailOver the Input will send to the first queue and only use the subsequent queues when a failure occurs
---
Where supported the FailOver configuration will likely make the most sense, though you may need to also send to other
queues to meet compliance and audit requirements. We could not possibly preordain all customer needs and thus elected
to simply make the system as flexible as possible.
*/

// The PluginQueueConfigurationTypeCode is our enum that we use to identify the various formats that can be used to
// determine the type of Queue mapping a configuration supplied.
type PluginQueueConfigurationTypeCode int

// The Basic PluginQueueConfigurationTypeCode is the configuration format that all Input plugins must inherently support.
// The FailOver and LoadBalanced PluginQueueConfigurationTypeCode are checked against the Input plugin interface to ensure
// that the configuration type is supported before it is supplied.
const (
	Basic PluginQueueConfigurationTypeCode = iota
	FailOver
	LoadBalanced
)

type PluginQueueConfiguration struct {
	Type         PluginQueueConfigurationTypeCode
	QueueIds     []string
	QueueConfigs []logwarp_common.IngestQueueConfig
}

// The NewPluginQueueConfiguration constructor makes it simple to map a Queue assignment to a plugin to the various types
// of assignments a user may want. The PluginQueueConfiguration struct helps organize this information to make sure that
// the software can easily tell which Input plugins support the configuration passed by checking the PluginQueueConfigurationTypeCode.
func NewPluginQueueConfiguration(configItem any, queues []logwarp_common.IngestQueueConfig) (PluginQueueConfiguration, error) {
	queueConfig := PluginQueueConfiguration{}
	var err error
	switch configItemType := configItem.(type) {
	case string:
		queueConfig.Type = Basic
		queueConfig.QueueIds = []string{strings.TrimSpace(configItemType)}
	case map[string]interface{}:
		val, ok := configItemType["load-balanced"]
		if ok {
			queueConfig.Type = LoadBalanced
			switch v := val.(type) {
			case []string:
				for i, s := range v {
					queueConfig.QueueIds[i] = strings.TrimSpace(s) // Cleanup any whitespace
				}
				queueConfig.QueueIds = v
			default:
				err = fmt.Errorf("invalid load-balanced configuration item type: %T (must be array of strings)", val)
			}
		} else {
			// Check for fail-over type
			val, ok = configItemType["fail-over"]
			if ok {
				queueConfig.Type = FailOver
				switch v := val.(type) {
				case []string:
					for i, s := range v {
						queueConfig.QueueIds[i] = strings.TrimSpace(s) // Cleanup any whitespace
					}
					queueConfig.QueueIds = v
				default:
					err = fmt.Errorf("invalid load-balanced configuration item type: %T (must be array of strings)", val)
				}
			}
		}
	default:
		err = fmt.Errorf("invalid configuration item type: %T (please reference the documentation for guiddance)", configItem)
	}
	for _, queueId := range queueConfig.QueueIds {
		for _, ingestQueueConfig := range queues {
			if queueId == ingestQueueConfig.Id {
				queueConfig.QueueConfigs = append(queueConfig.QueueConfigs, ingestQueueConfig)
			}
		}
	}
	return queueConfig, err
}
