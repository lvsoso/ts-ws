
### grpc-gateway
```shell
go get github.com/golang/protobuf/protoc-gen-go

go get google.golang.org/grpc/cmd/protoc-gen-go-grpc

go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway

go get github.com/golang/protobuf/descriptor@v1.5.0
go get github.com/golang/protobuf/proto@v1.5.0
go get google.golang.org/genproto/googleapis/api/annotations@v0.0.0-20200526211855-cb27e3aa2013

cd ${GOPATH}/src/github.com/
git clone https://github.com/googleapis/googleapis.git


# https://github.com/grpc-ecosystem/grpc-gateway/issues/422
```

```shell
# ws://localhost:8080/ws
#ws client send
{
    "client_id":"hahahahha",
    "task_ids":[1],
    "op":"STATUS"
}

# response
{
    "task_id": 1,
    "status": 3
}
```

ssh -qfN -D localhost:1080 oversea-host
ss -tln | grep 1080
all_proxy=socks5://localhost:1080/ curl -kIsS http://www.google.com/

all_proxy=socks5h://localhost:1080/ curl -kIsS http://www.google.com/