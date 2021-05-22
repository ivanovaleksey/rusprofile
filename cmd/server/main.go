package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ivanovaleksey/rusprofile/app/server"
	"github.com/ivanovaleksey/rusprofile/pkg/pb/rusprofile"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"net/http"
)

func main() {
	zapLogger, _ := zap.NewDevelopment()
	grpc_zap.ReplaceGrpcLoggerV2(zapLogger)

	handler := func(p interface{}) (err error) {
		zapLogger.Error("panic occurred", zap.Any("err", p))
		return status.Errorf(codes.Internal, "%v", p)
	}

	srv := grpc.NewServer(grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(handler)),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(zapLogger),
		),
	))
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
