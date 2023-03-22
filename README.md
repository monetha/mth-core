# mth-core

This repository hosts common code for our services, which help with logging, authentication, API middlewares etc.

This project produces no runnable binary. However, you can check `test` target in Makefile if you're wondering how to run all unit tests.

## Generated files

This library uses .proto files to generate `.pb.go` files. For that - you need to have protoc compiler installed on your system. Installation instructions here: https://grpc.io/docs/protoc-installation/ .

After `protoc` compiler is installed - simply run `make generate-proto` to regenerate files
