package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/consul"
	"movieexample.com/rating/internal/controller/rating"
	httphandler "movieexample.com/rating/internal/handler/http"
	"movieexample.com/rating/internal/repository/memory"
)

const serviceName = "rating"

func main() {
	var port int
	flag.IntVar(&port, "port", 8082, "API handler port")
	flag.Parse()

	log.Printf("Starting the rating service on port: %d", port)

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

	repo := memory.New()
	svc := rating.New(repo)
	h := httphandler.New(svc)
	http.Handle("/rating", http.HandlerFunc(h.Handle))
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			panic(err)
		}
	}()

	<-signalChan
	log.Println("Received shutdown signal, stopping service...")
}
