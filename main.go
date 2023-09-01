package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	healthcheck "github.com/Junkes887/load-balancer/health_check"
	loadbalancer "github.com/Junkes887/load-balancer/load_balancer"
	serverpool "github.com/Junkes887/load-balancer/server_pool"
	"github.com/joho/godotenv"
)

const (
	Attempts int = iota
	retries      = 0
)

func main() {
	godotenv.Load()
	serverPool := &serverpool.ServerPool{}

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	serverList := os.Getenv("SERVER_LIST")

	urls := strings.Split(serverList, ",")

	lb := loadbalancer.LB{SP: serverPool}

	lb.AddBackends(urls)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(lb.Proxy),
	}

	go healthcheck.LauchHealthCheck(serverPool)

	log.Printf("Load Balancer started at :%d\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
