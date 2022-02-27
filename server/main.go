package main

import (
	"context"
	"log"
	"net"
	"net/http"

	helloworldpb "server/proto/helloworld"

	"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/httprule"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	helloworldpb.UnimplementedGreeterServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) SayHello(ctx context.Context, in *helloworldpb.HelloRequest) (*helloworldpb.HelloReply, error) {
	return &helloworldpb.HelloReply{Message: in.Name + " world"}, nil
}

func main() {

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8181")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// enable reflection
	reflection.Register(s)

	ctx := context.TODO()
	gmux := runtime.NewServeMux()
	// Attach the Greeter service to the server
	server := &server{}
	err = helloworldpb.RegisterGreeterHandlerServer(ctx, gmux, server)
	if err != nil {
		panic(err)
	}

	gmux.Handle("GET", makePattern("/ws"), handleWebSocket)

	mux := http.NewServeMux()
	mux.Handle("/", interceptor(gmux))
	// mux.Handle("/", gmux)

	// Serve gRPC Server
	log.Println("Serving gRPC on 0.0.0.0:8181")
	go s.Serve(lis)

	// Serve http Server
	log.Println("Serving http on 0.0.0.0:8080")
	http.ListenAndServe(":8080", mux)
}

func interceptor(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println(r.Header.Get("Autorization"))
		handler.ServeHTTP(w, r)
	})
}

func makePattern(rule string) runtime.Pattern {
	var tmpl httprule.Template
	if compiler, err := httprule.Parse(rule); err != nil {
		panic(err)
	} else {
		tmpl = compiler.Compile()
	}

	pattern, err := runtime.NewPattern(1, tmpl.OpCodes, tmpl.Pool, tmpl.Verb)
	if err != nil {
		panic(err)
	}
	return pattern
}
