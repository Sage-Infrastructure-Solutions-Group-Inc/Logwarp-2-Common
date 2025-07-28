package ingest_queue_plugin

import (
	"bytes"
	"fmt"
	"github.com/Sage-Infrastructure-Solutions-Group-Inc/Logwarp-2-Common/protobuf"
	"github.com/golang/snappy"
	"github.com/pierrec/lz4/v4"
	"log/slog"
)

func strToCompression(comp string) (protobuf.CompressionMethod, error) {
	if comp == "lz4" {
		return protobuf.CompressionMethod_LZ4, nil
	} else if comp == "snappy" {
		return protobuf.CompressionMethod_SNAPPY, nil
	} else if comp == "none" || comp == "" {
		return protobuf.CompressionMethod_NONE, nil
	}
	return protobuf.CompressionMethod_NONE, fmt.Errorf("unknown compression method: %s", comp)
}
func numberToCompressionMethod(number float64) (protobuf.CompressionMethod, error) {
	if int(number) == 0 {
		return protobuf.CompressionMethod_NONE, nil
	} else if int(number) == 1 {
		return protobuf.CompressionMethod_LZ4, nil
	} else if int(number) == 2 {
		return protobuf.CompressionMethod_SNAPPY, nil
	}
	return protobuf.CompressionMethod_NONE, fmt.Errorf("unknown compression method: %f", number)
}

func performCompression(encodedRecords []byte, compressionMethod protobuf.CompressionMethod, logGroup []any) ([]byte, error) {
	finalBytes := make([]byte, 0)
	var err error
	switch compressionMethod {
	case protobuf.CompressionMethod_LZ4:
		writeBuffer := bytes.Buffer{}
		writer := lz4.NewWriter(&writeBuffer)
		defer writer.Close()
		_, err = writer.Write(encodedRecords)
		if err != nil {
			slog.Error(fmt.Sprintf("Compression Error: %v", err.Error()), logGroup...)
		}
		err = writer.Flush()
		if err != nil {
			slog.Error(fmt.Sprintf("Compression Error: %v (Unable to Flush Writer)", err.Error()), logGroup...)
		}
		finalBytes = writeBuffer.Bytes()
	case protobuf.CompressionMethod_SNAPPY:
		writeBuffer := bytes.Buffer{}
		writer := snappy.NewBufferedWriter(&writeBuffer)
		defer writer.Close()
		_, err = writer.Write(encodedRecords)
		if err != nil {
			slog.Error(fmt.Sprintf("Compression Error: %v", err.Error()), logGroup...)
		}
		err = writer.Flush()
		if err != nil {
			slog.Error(fmt.Sprintf("Compression Error: %v (Unable to Flush Writer)", err.Error()), logGroup...)
		}
		finalBytes = writeBuffer.Bytes()
	default:
		finalBytes = append(finalBytes, encodedRecords...)
	}
	return finalBytes, err

}
