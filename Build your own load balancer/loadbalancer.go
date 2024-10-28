package main

import (
    "bufio"
    "io"
    "log"
    "net"
    "net/http"
    "sync"
)

var (
    servers     = []string{"localhost:8080", "localhost:8081"}
    serverIndex = 0
    mu          sync.Mutex
)

func getNextServer() string {
    mu.Lock()
    defer mu.Unlock()
    server := servers[serverIndex]
    serverIndex = (serverIndex + 1) % len(servers)
    return server
}

func handleConnection(conn net.Conn) {
    defer conn.Close()

    // Read the incoming request
    req, err := http.ReadRequest(bufio.NewReader(conn))
    if err != nil {
        log.Println("Error reading request:", err)
        return
    }

    log.Printf("Received request from %s\n", req.RemoteAddr)

    // Forward the request to the backend server
    backendServer := getNextServer()
    client := &http.Client{}
    req.URL.Host = backendServer
    req.URL.Scheme = "http"
    req.RequestURI = ""

    resp, err := client.Do(req)
    if err != nil {
        log.Println("Error forwarding request:", err)
        return
    }
    defer resp.Body.Close()

    // Send the response back to the client
    conn.Write([]byte(resp.Proto + " " + resp.Status + "\r\n"))
    for key, value := range resp.Header {
        for _, v := range value {
            conn.Write([]byte(key + ": " + v + "\r\n"))
        }
    }
    conn.Write([]byte("\r\n"))
    io.Copy(conn, resp.Body)
}

func main() {
    listener, err := net.Listen("tcp", ":80")
    if err != nil {
        log.Fatal("Error starting load balancer:", err)
    }
    defer listener.Close()

    log.Println("Load balancer listening on port 80...")
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println("Error accepting connection:", err)
            continue
        }
        go handleConnection(conn)
    }
}
