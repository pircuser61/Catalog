package main

import (
	"context"
	_ "embed"
	"net"
	"net/http"

	"github.com/flowchartsman/swaggerui"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "gitlab.ozon.dev/pircuser61/catalog/api"
	apiPkg "gitlab.ozon.dev/pircuser61/catalog/internal/api"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grcpAddr = ":8081"
	httpAddr = ":8080"
)

//go:embed swagger/api.swagger.json
var spec []byte

func runGRPCServer(good goodPkg.Interface) {
	listener, err := net.Listen("tcp", grcpAddr)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCatalogServer(grpcServer, apiPkg.New(good))

	if err = grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

func runREST(ctx context.Context) {
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterCatalogHandlerFromEndpoint(ctx, gwmux, grcpAddr, opts); err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handler(spec)))
	mux.Handle("/", gwmux)
	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		panic(err)
	}
}
