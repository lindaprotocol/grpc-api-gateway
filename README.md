# grpc-api-gateway
# Linda Protocol gRPC API Gateway

A high-performance gRPC to HTTP/JSON gateway for the Linda blockchain. This service exposes Linda's gRPC APIs as RESTful JSON endpoints, making it easy to integrate with web applications and services.

## Features

- **Complete Linda API Coverage**: FullNode, SolidityNode, and ScanService APIs
- **RESTful JSON Endpoints**: All gRPC methods exposed as HTTP endpoints
- **CORS Support**: Ready for browser-based applications
- **Hex Encoding**: All binary data (hashes, signatures) returned as hex strings
- **Base58 Addresses**: Linda addresses in standard Base58 format
- **Extensible**: Easy to add custom endpoints and middleware

## Prerequisites

- Go 1.23 or higher
- Protocol Buffers compiler (protoc) v3.21+
- Linda Full Node (optional, for testing)

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/lindaprotocol/grpc-api-gateway.git
cd grpc-api-gateway

# Install dependencies
make deps

# Generate protocol buffers
make proto

# Build the gateway
make build

# Test API of linda http
curl -X POST -k http://localhost:8086/wallet/listwitnesses
```
If you get witness-list json data, congratulations

> Note: json to protobuf, bytes should be passed via base64 formate

# grpc-gateway

grpc-gateway is a plugin of [protoc](http://github.com/google/protobuf).
It reads [gRPC](http://github.com/grpc/grpc-common) service definition,
and generates a reverse-proxy server which translates a RESTful JSON API into gRPC.
This server is generated according to [custom options](https://cloud.google.com/service-management/reference/rpc/google.api#http) in your gRPC definition.

It helps you to provide your APIs in both gRPC and RESTful style at the same time.

![architecture introduction diagram](https://docs.google.com/drawings/d/12hp4CPqrNPFhattL_cIoJptFvlAqm5wLQ0ggqI5mkCg/pub?w=749&h=370)

## Background

gRPC is great -- it generates API clients and server stubs in many programming languages. It is fast, easy-to-use, bandwidth-efficient and its design is combat-proven by Google.
However, you might still want to provide a traditional RESTful API as well. Reasons can range from maintaining backwards-compatibility, supporting languages or clients not well supported by gRPC to simply maintaining the aesthetics and tooling involved with a RESTful architecture.

This project aims to provide that HTTP+JSON interface to your gRPC service. A small amount of configuration in your service to attach HTTP semantics is all that's needed to generate a reverse-proxy with this library.

## Parameters and flags
`protoc-gen-grpc-gateway` supports custom mapping from Protobuf `import` to Golang import path.
They are compatible to [the parameters with same names in `protoc-gen-go`](https://github.com/golang/protobuf#parameters).

In addition, we also support the `request_context` parameter in order to use the `http.Request`'s Context (only for Go 1.7 and above).
This parameter can be useful to pass request scoped context between the gateway and the gRPC service.

`protoc-gen-grpc-gateway` also supports some more command line flags to control logging. You can give these flags together with parameters above. Run `protoc-gen-grpc-gateway --help` for more details about the flags.



## Features
### Supported
* Generating JSON API handlers
* Method parameters in request body
* Method parameters in request path
* Method parameters in query string
* Enum fields in path parameter (including repeated enum fields).
* Mapping streaming APIs to newline-delimited JSON streams
* Mapping HTTP headers with `Grpc-Metadata-` prefix to gRPC metadata (prefixed with `grpcgateway-`)
* Optionally emitting API definition for [Swagger](http://swagger.io).
* Setting [gRPC timeouts](http://www.grpc.io/docs/guides/wire.html) through inbound HTTP `Grpc-Timeout` header.



## Mapping gRPC to HTTP

* [How gRPC error codes map to HTTP status codes in the response](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/runtime/errors.go#L15)
* HTTP request source IP is added as `X-Forwarded-For` gRPC request header
* HTTP request host is added as `X-Forwarded-Host` gRPC request header
* HTTP `Authorization` header is added as `authorization` gRPC request header 
* Remaining Permanent HTTP header keys (as specified by the IANA [here](http://www.iana.org/assignments/message-headers/message-headers.xhtml) are prefixed with `grpcgateway-` and added with their values to gRPC request header
* HTTP headers that start with 'Grpc-Metadata-' are mapped to gRPC metadata (prefixed with `grpcgateway-`)
* While configurable, the default {un,}marshaling uses [jsonpb](https://godoc.org/github.com/golang/protobuf/jsonpb) with `OrigName: true`.

## License
grpc-gateway is licensed under the BSD 3-Clause License.
See [LICENSE.txt](https://github.com/lindaprotocol/grpc-api-gateway/blob/main/LICENSE) for more details.
