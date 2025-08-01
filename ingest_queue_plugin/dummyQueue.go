package ingest_queue_plugin

import (
	"fmt"
	"github.com/Sage-Infrastructure-Solutions-Group-Inc/Logwarp-2-Common/protobuf"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
	"os"
	"time"
)

type DummyQueuePlugin struct {
	logGroup    []any
	compression protobuf.CompressionMethod
	envVars     map[string]*structpb.Value
	ThrowError  bool
}

func (d *DummyQueuePlugin) Config(config map[string]interface{}) error {
	var err error
	d.logGroup = append(d.logGroup, slog.Group("source", "plugin_name", "queue_plugin"))
	slog.Debug(fmt.Sprintf("Dummy Queue Plugin Config: %v", config), d.logGroup...)
	// Configure the compression method
	d.envVars = make(map[string]*structpb.Value)
	switch tmpMetaEnvVars := config["resolve-meta-env-vars"].(type) {
	case []string:
		for _, v := range tmpMetaEnvVars {
			d.envVars[v] = structpb.NewStringValue(os.Getenv(v))
		}
		fmt.Printf("Environment variables are: %v", d.envVars)
		slog.Debug(fmt.Sprintf("Environment variables are: %v", d.envVars), d.logGroup...)
	default:
		slog.Warn(fmt.Sprintf("No value resolve-meta-env-vars defined on the queue plugin configuration!"), d.logGroup...)
	}
	switch c := config["compression"].(type) {
	case string:
		d.compression, err = strToCompression(c)
		if err != nil {
			slog.Warn(fmt.Sprintf("The provided compression method is not supported: %v (Defaulting to None)"))
		}
	case int:
		d.compression, err = numberToCompressionMethod(float64(c))
		if err != nil {
			slog.Warn(fmt.Sprintf("The provided compression method is not supported: %v (Defaulting to None)"))
		}
	case float64:
		d.compression, err = numberToCompressionMethod(c)
		if err != nil {
			slog.Warn(fmt.Sprintf("The provided compression method is not supported: %v (Defaulting to None)"))
		}
	case float32:
		d.compression, err = numberToCompressionMethod(float64(c))
		if err != nil {
			slog.Warn(fmt.Sprintf("The provided compression method is not supported: %v (Defaulting to None)"))
		}
	default:
		slog.Warn(fmt.Sprintf("The provided compression method is not supported: %v (Defaulting to None)"))
		d.compression = protobuf.CompressionMethod(0)
	}
	d.ThrowError = false // for testing error handling on queue failure.
	return nil
}

func (d *DummyQueuePlugin) Enqueue(records protobuf.RecordList, inputPluginName string) error {
	//slog.Debug(fmt.Sprintf("Enqueue Batch: %v", records), d.logGroup...)
	if len(records.Records) == 0 {
		slog.Debug(fmt.Sprintf("No records to process (The length of the submitted data is: %d", len(records.Records)), d.logGroup...)
		return nil
	}
	bytesEncodedRecords, err := proto.Marshal(&records)
	if err != nil {
		slog.Error(fmt.Sprintf("Enqueue Error: %v", err.Error()), d.logGroup...)
		return err
	}
	finalRecordBytes, err := performCompression(bytesEncodedRecords, d.compression, d.logGroup)
	if err != nil {
		slog.Error(fmt.Sprintf("Enqueue Error: %v", err.Error()), d.logGroup...)
	}
	if len(finalRecordBytes) == 0 {
		slog.Error(fmt.Sprintf("Empty message for queue due to previous failures."), d.logGroup...)
		return err
	}
	hostname, err := os.Hostname()
	if err != nil {
		slog.Warn(fmt.Sprintf("Enqueue Error: %v (Unable to resolve local hostname)", err.Error()), d.logGroup...)
	}
	batch := protobuf.Batch{
		Records:           finalRecordBytes,
		Timestamp:         timestamppb.New(time.Now()),
		Compression:       d.compression,
		InputPlugin:       inputPluginName,
		QueuePlugin:       "DummyQueuePlugin",
		SubmitterHostname: hostname,
	}
	batchBytes, err := proto.Marshal(&batch)
	if err != nil {
		slog.Error(fmt.Sprintf("Enqueue Error: %v (Could not create Batch)", err.Error()), d.logGroup...)
		return err
	}
	if d.ThrowError {
		slog.Debug(fmt.Sprintf("Enqueue Error: %v (For Testing)", d.ThrowError), d.logGroup...)
		return fmt.Errorf("Enqueue Error: (For Testing)")
	}
	slog.Debug(fmt.Sprintf("Enqueue Batch with Length: %d Values: %v", len(batchBytes), batchBytes), d.logGroup...)
	slog.Warn(fmt.Sprintf("This is a dummy queue and as such no data is stored"), d.logGroup...)
	return nil
}

func (d *DummyQueuePlugin) Test() bool {
	return true
}

var interfaceTest QueuePlugin = &DummyQueuePlugin{}

func NewQueue() QueuePlugin {
	return &DummyQueuePlugin{}
}
