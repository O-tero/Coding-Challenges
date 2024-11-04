package main

import (
    "flag"
    "log"
    "net/http"
    "strings"
    "time"

    "load-balancer/server"
)

func main() {
    // Define command-line flags
    serversFlag := flag.String("servers", "", "Comma-separated list of backend server URLs")
    healthCheckURL := flag.String("health-check-url", "/health", "URL to check server health")
    healthInterval := flag.Int("health-interval", 10, "Health check interval in seconds")
    port := flag.String("port", "8000", "Port for the load balancer")

    flag.Parse()

    // Parse the servers flag
    servers := strings.Split(*serversFlag, ",")
    if len(servers) == 0 || *serversFlag == "" {
        log.Fatal("Please provide at least one backend server URL using the --servers flag")
    }

    // Initialize the load balancer
    lb := server.NewLoadBalancer(servers, *healthCheckURL, time.Duration(*healthInterval)*time.Second)

    log.Printf("Starting load balancer on port %s...\n", *port)
    log.Fatal(http.ListenAndServe(":"+*port, lb))
}
