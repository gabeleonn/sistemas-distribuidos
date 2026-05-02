# Raft

## Client

## Cluster

## Setup

### Protoc
Geração de stubs deve ser feita através dos comandos:

```bash
cd proto
protoc \
  --go_out=../autogen --go_opt=paths=source_relative \
  --go-grpc_out=../autogen --go-grpc_opt=paths=source_relative \
  node.proto
```
