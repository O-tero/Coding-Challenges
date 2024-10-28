package main

import (
    "flag"
    "log"
    "net/http"
    "time"

    "load-balancer/server" 
)

func main() {
    servers := []string{"http://localhost:8080", "http://localhost:8081", "http://localhost:8082"}

    healthCheckURL := flag.String("health-check-url", "/health", "URL to check server health")
    healthInterval := flag.Int("health-interval", 10, "Health check interval in seconds")
    port := flag.String("port", "8000", "Port for the load balancer")

    flag.Parse()

    lb := loadbalancer.NewLoadBalancer(servers, *healthCheckURL, time.Duration(*healthInterval)*time.Second)

    log.Printf("Starting load balancer on port %s...\n", *port)
    log.Fatal(http.ListenAndServe(":"+*port, lb))
}
