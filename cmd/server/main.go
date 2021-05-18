package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ivanovaleksey/rusprofile/pkg/pb/rusprofile"
	"github.com/ivanovaleksey/rusprofile/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
)

func main() {
	srv := grpc.NewServer()
	impl := server.NewServer()

	rusprofile.RegisterRusProfileServiceServer(srv, impl)
	reflection.Register(srv)

	go func() {
		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{
			grpc.WithInsecure(),
		}
		regErr := rusprofile.RegisterRusProfileServiceHandlerFromEndpoint(context.Background(), mux, "127.0.0.1:7001", opts)
		if regErr != nil {
			log.Fatal(regErr)
		}
		err := http.ListenAndServe(":7002", mux)
		fmt.Println(err)
	}()

	grpcListener, _ := net.Listen("tcp", ":7001")
	err := srv.Serve(grpcListener)
	fmt.Println(err)
}
