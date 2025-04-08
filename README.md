# gRPC CRUD with Go

This project demonstrates a simple gRPC-based CRUD system in Golang, including both server and client implementations.

## Prerequisites

Make sure the following tools are installed on your system:

- Go 1.20+
- `protoc` (Protocol Buffers compiler)
- `protoc-gen-go` and `protoc-gen-go-grpc` plugins

## Installation

### Install Protocol Buffers Compiler

```bash
sudo apt update
sudo apt install protobuf-compiler
protoc --version  # Check version
```

### Setup Go Environment

```bash
export PATH=$PATH:$(go env GOPATH)/bin
source ~/.bashrc
```

### Install Go Plugins for Protobuf

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Check installed versions
protoc-gen-go --version
protoc-gen-go-grpc --version
```

## Project Setup

### Initialize Go Module

```bash
go mod init mubashir-crud
```

### Generate Go Code from .proto File

```bash
protoc --go_out=. --go-grpc_out=. proto/user.proto
```

### Install Required Dependencies

```bash
go get google.golang.org/grpc
go get google.golang.org/protobuf
go get github.com/google/uuid
```

### Sync Vendor Directory

```bash
go mod vendor
```

## Running the Project

### Run the gRPC Server

```bash
go run server/main.go
```

### Run the gRPC Client

```bash
go run client/main.go
```

## Project Structure

```
.
├── client/
│   └── main.go
├── proto/
│   └── user.proto
├── server/
│   └── main.go
├── go.mod
├── go.sum
└── README.md
```

## License

This project is licensed under the MIT License.

---

Happy Coding! 🚀