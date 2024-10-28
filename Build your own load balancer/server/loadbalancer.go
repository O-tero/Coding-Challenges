package loadbalancer

import (
    "io"
    "log"
    "net/http"
    "sync"
    "time"
)

type Server struct {
    URL     string
    Healthy bool
}

type LoadBalancer struct {
    servers       []*Server
    currentIndex  int
    mutex         sync.Mutex
    healthCheckURL string
    healthInterval time.Duration
}

// NewLoadBalancer initializes the load balancer with servers and a health check URL.
func NewLoadBalancer(servers []string, healthCheckURL string, healthInterval time.Duration) *LoadBalancer {
    srv := make([]*Server, len(servers))
    for i, url := range servers {
        srv[i] = &Server{URL: url, Healthy: true}
    }
    lb := &LoadBalancer{
        servers:       srv,
        healthCheckURL: healthCheckURL,
        healthInterval: healthInterval,
    }
    go lb.startHealthCheck()
    return lb
}

// startHealthCheck performs periodic health checks on backend servers.
func (lb *LoadBalancer) startHealthCheck() {
    for {
        time.Sleep(lb.healthInterval)
        for _, server := range lb.servers {
            lb.checkServerHealth(server)
        }
    }
}

// checkServerHealth checks the health of a single server.
func (lb *LoadBalancer) checkServerHealth(server *Server) {
    resp, err := http.Get(server.URL + lb.healthCheckURL)
    if err != nil || resp.StatusCode != http.StatusOK {
        server.Healthy = false
        log.Printf("Server %s is unhealthy\n", server.URL)
    } else {
        server.Healthy = true
        log.Printf("Server %s is healthy\n", server.URL)
    }
    if resp != nil {
        resp.Body.Close()
    }
}

// getNextServer uses round-robin scheduling among healthy servers.
func (lb *LoadBalancer) getNextServer() *Server {
    lb.mutex.Lock()
    defer lb.mutex.Unlock()

    for i := 0; i < len(lb.servers); i++ {
        server := lb.servers[lb.currentIndex]
        lb.currentIndex = (lb.currentIndex + 1) % len(lb.servers)
        if server.Healthy {
            return server
        }
    }
    return nil
}

// ServeHTTP handles incoming requests and forwards them to the next healthy server.
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    server := lb.getNextServer()
    if server == nil {
        http.Error(w, "No healthy servers available", http.StatusServiceUnavailable)
        return
    }

    resp, err := http.Get(server.URL + r.RequestURI)
    if err != nil {
        http.Error(w, "Failed to reach backend server", http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    for k, v := range resp.Header {
        for _, vv := range v {
            w.Header().Add(k, vv)
        }
    }
    w.WriteHeader(resp.StatusCode)
    _, _ = io.Copy(w, resp.Body)
}
