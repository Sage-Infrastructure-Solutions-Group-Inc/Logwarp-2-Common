# LogWarp Common
This repository contains common types, interfaces, and functions 
required to integrate your own custom plugins with LogWarp 2


## Generating the Protobuf Specification
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
export PATH="$PATH:$(go env GOPATH)/bin"
protoc protobuf/batch.proto --go_out=.
```