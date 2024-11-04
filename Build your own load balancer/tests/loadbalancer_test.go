package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
	"errors"
	"net/url"
	"net/http/httputil"
)

type LoadBalancer struct {
	servers []*url.URL
	current int
	mu      sync.Mutex
}

func NewLoadBalancer(serverURLs []string, healthCheckPath string, healthCheckInterval time.Duration) (*LoadBalancer, error) {
	var servers []*url.URL
	for _, serverURL := range serverURLs {
		parsedURL, err := url.Parse(serverURL)
		if err != nil {
			return nil, err
		}
		servers = append(servers, parsedURL)
	}
	if len(servers) == 0 {
		return nil, errors.New("no servers provided")
	}
	lb := &LoadBalancer{servers: servers}
	go lb.healthCheck(healthCheckPath, healthCheckInterval)
	return lb, nil
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server := lb.nextServer()
	proxy := httputil.NewSingleHostReverseProxy(server)
	proxy.ServeHTTP(w, r)
}

func (lb *LoadBalancer) nextServer() *url.URL {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	server := lb.servers[lb.current]
	lb.current = (lb.current + 1) % len(lb.servers)
	return server
}

func (lb *LoadBalancer) healthCheck(path string, interval time.Duration) {
	for {
		time.Sleep(interval)
		for _, server := range lb.servers {
			resp, err := http.Get(server.String() + path)
			if err != nil || resp.StatusCode != http.StatusOK {
				lb.removeServer(server)
			}
		}
	}
}

func (lb *LoadBalancer) removeServer(server *url.URL) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	for i, s := range lb.servers {
		if s == server {
			lb.servers = append(lb.servers[:i], lb.servers[i+1:]...)
			break
		}
	}
}

func TestLoadBalancerWithMultipleClients(t *testing.T) {
	// Create mock backend servers
	serverA := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Response from Server A"))
	}))
	defer serverA.Close()

	serverB := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Response from Server B"))
	}))
	defer serverB.Close()

	// Extract server URLs
	serverURLs := []string{serverA.URL, serverB.URL}

	// Create a LoadBalancer instance
	lb, err := NewLoadBalancer(serverURLs, "/health", time.Second*2)
	if err != nil {
		t.Fatalf("Failed to create load balancer: %v", err)
	}

	// Create an HTTP test server to act as the load balancer's entry point
	loadBalancerServer := httptest.NewServer(lb)
	defer loadBalancerServer.Close()

	// Simulate multiple clients accessing the load balancer
	numClients := 10
	var wg sync.WaitGroup
	wg.Add(numClients)

	clientResponses := make([]string, numClients)
	clientErrors := make([]error, numClients)

	for i := 0; i < numClients; i++ {
		go func(index int) {
			defer wg.Done()
			resp, err := http.Get(loadBalancerServer.URL)
			if err != nil {
				clientErrors[index] = err
				return
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			clientResponses[index] = string(body)
		}(i)
	}

	// Wait for all clients to complete their requests
	wg.Wait()

	// Verify the responses are distributed among backends
	serverAResponses := 0
	serverBResponses := 0

	for _, response := range clientResponses {
		if response == "Response from Server A" {
			serverAResponses++
		} else if response == "Response from Server B" {
			serverBResponses++
		} else {
			t.Errorf("Unexpected response: %s", response)
		}
	}

	t.Logf("Server A handled %d requests; Server B handled %d requests", serverAResponses, serverBResponses)

	// Check for errors from clients
	for i, err := range clientErrors {
		if err != nil {
			t.Errorf("Client %d encountered an error: %v", i, err)
		}
	}

	// Ensure load balancing occurred (approximately even distribution is expected)
	if serverAResponses == 0 || serverBResponses == 0 {
		t.Errorf("One of the servers received no requests; serverA: %d, serverB: %d", serverAResponses, serverBResponses)
	}
}
