package server

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type LoadBalancer struct {
	servers         []*url.URL
	currentIndex    int
	mutex           sync.Mutex
	healthCheckURL  string
	healthInterval  time.Duration
	serverStatus    map[int]bool
}

func NewLoadBalancer(serverList []string, healthCheckURL string, healthInterval time.Duration) *LoadBalancer {
	var servers []*url.URL
	for _, server := range serverList {
		parsedURL, err := url.Parse(server)
		if err != nil {
			log.Fatalf("Error parsing server URL %s: %v", server, err)
		}
		servers = append(servers, parsedURL)
	}

	lb := &LoadBalancer{
		servers:         servers,
		healthCheckURL:  healthCheckURL,
		healthInterval:  healthInterval,
		serverStatus:    make(map[int]bool),
	}

	// Initialize all servers as healthy and start health checks
	for i := range lb.servers {
		lb.serverStatus[i] = true
	}
	go lb.startHealthChecks()

	return lb
}

func (lb *LoadBalancer) startHealthChecks() {
	for {
		for i, server := range lb.servers {
			go func(index int, server *url.URL) {
				healthURL := server.ResolveReference(&url.URL{Path: lb.healthCheckURL})
				resp, err := http.Get(healthURL.String())
				if err != nil || resp.StatusCode != http.StatusOK {
					log.Printf("Server %s is down (Error: %v, StatusCode: %d)", server, err, resp.StatusCode)
					lb.serverStatus[index] = false
				} else {
					log.Printf("Server %s is healthy", server)
					lb.serverStatus[index] = true
				}
				if resp != nil {
					resp.Body.Close()
				}
			}(i, server)
		}
		time.Sleep(lb.healthInterval)
	}
}

func (lb *LoadBalancer) getNextAvailableServer() (*url.URL, int) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	for i := 0; i < len(lb.servers); i++ {
		lb.currentIndex = (lb.currentIndex + 1) % len(lb.servers)
		if lb.serverStatus[lb.currentIndex] {
			return lb.servers[lb.currentIndex], lb.currentIndex
		}
	}

	log.Println("No available servers.")
	return nil, -1
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the next available server using round-robin
	server, _ := lb.getNextAvailableServer()
	if server == nil {
		http.Error(w, "No available servers", http.StatusServiceUnavailable)
		return
	}

	// Log the request and selected server
	log.Printf("Received request from %s for %s, routing to server: %s", r.RemoteAddr, r.URL.Path, server)

	// Forward the request to the selected server
	proxyURL := server.ResolveReference(r.URL)
	proxyReq, _ := http.NewRequest(r.Method, proxyURL.String(), r.Body)
	proxyReq.Header = r.Header

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Printf("Error forwarding request to %s: %v", server, err)
		http.Error(w, "Error forwarding request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Measure response time
	duration := time.Since(startTime)
	log.Printf("Response from server %s took %v", server, duration)

	// Write response back to the client
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
