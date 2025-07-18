package logwarp_common

import "time"

type CompressionMethod int

const CompressionLZ4 = CompressionMethod(1)
const CompressionSnappy = CompressionMethod(2)
const CompressionNone = CompressionMethod(0)

type Batch struct {
	Records     []Record
	Timestamp   time.Time
	Compression CompressionMethod
}

type Record struct {
	Timestamp time.Time
	Content   []byte
}

type IngestQueueConfig struct {
	Plugin       string                 `yaml:"plugin"`
	PluginConfig map[string]interface{} `yaml:"plugin-config"` // Use a map to store varying options
	Compression  string                 `yaml:"compression"`
	Id           string                 `yaml:"id"`
}
