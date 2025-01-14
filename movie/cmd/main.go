package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"movieexample.com/gen"
	"movieexample.com/movie/internal/controller/movie"
	metadatagateway "movieexample.com/movie/internal/gateway/metadata/http"
	ratinggateway "movieexample.com/movie/internal/gateway/rating/http"
	grpchandler "movieexample.com/movie/internal/handler/grpc"
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/consul"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const serviceName = "movie"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()

	log.Printf("Starting the movie service on port: %d", port)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer func() {
		log.Println("Deregistering service from Consul...")
		if err := registry.Deregister(ctx, instanceID, serviceName); err != nil {
			log.Println("Failed to deregister from Consul:", err)
		} else {
			log.Println("Successfully deregistered from Consul")
		}
	}()

	metadatagateway := metadatagateway.New(registry)
	ratinggateway := ratinggateway.New(registry)
	ctrl := movie.New(ratinggateway, metadatagateway)
	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterMovieServiceServer(srv, h)
	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-signalChan
	log.Println("Received shutdown signal, stopping service...")
	srv.GracefulStop()
	log.Println("gRPC server stopped")
}
