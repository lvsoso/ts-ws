.PHONY: gen

gen:
	protoc -I ./proto  -I ${GOPATH}/src/github.com/googleapis -I /usr/local/include \
   --go_out ./proto --go_opt paths=source_relative \
   --go-grpc_out ./proto --go-grpc_opt paths=source_relative \
   --grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative \
   ./proto/helloworld/hello_world.proto