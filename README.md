# Prerequisites

- docker https://docs.docker.com/docker-for-mac/install/
- go `brew install go --cross-compile-common`
- protobufs `brew install protobuf`
- go dep `brew install dep`
- golang protobufs `go get -u github.com/golang/protobuf/protoc-gen-go`
- grpc gateway `go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway`
- grpc gateway swagger `go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger`
- golang grpc `go get -u google.golang.org/grpc`

Note: we are relying on grpc-gateway@master (see Gopkg.toml) to simplify setup. Ideally we'd depend on v1.2.2, which is what `dep init` determines to be the latest version, but `go get` always checks out master, which generates code that is not compatible with the v1.2.2 runtime. This means that these installation instructions will break the next time a backwards-incompatible change is committed in grpc-gateway. When that happens run `dep ensure -update github.com/grpc-ecosystem/grpc-gateway && go generate`.

# Building the server
- `dep ensure`
- `go build`
- Note that protos are generated with `go generate`
