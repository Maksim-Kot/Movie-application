package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Maksim-Kot/Movie-application/gen"
	"github.com/Maksim-Kot/Movie-application/movie/internal/controller/movie"
	metadatagateway "github.com/Maksim-Kot/Movie-application/movie/internal/gateway/metadata/http"
	ratinggateway "github.com/Maksim-Kot/Movie-application/movie/internal/gateway/rating/http"
	grpchandler "github.com/Maksim-Kot/Movie-application/movie/internal/handler/grpc"
	"github.com/Maksim-Kot/Movie-application/pkg/discovery"
	"github.com/Maksim-Kot/Movie-application/pkg/discovery/consul"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
)

func main() {
	f, err := os.Open("configs/base.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var cfg config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}
	port := cfg.API.Port
	serviceName := cfg.API.Name
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
