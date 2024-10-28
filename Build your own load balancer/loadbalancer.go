package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

const backendAddr = "localhost:8081" 

func handleConnection(conn net.Conn) {
    defer conn.Close()

    // Read the incoming request
    req, err := http.ReadRequest(bufio.NewReader(conn))
    if err != nil {
        log.Println("Error reading request:", err)
        return
    }

    log.Printf("Received request from %s\n", req.RemoteAddr)

    resp, err := http.DefaultTransport.RoundTrip(req)
    if err != nil {
        log.Println("Error forwarding request:", err)
        return
    }
    defer resp.Body.Close()

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
        log.Fatal("Error starting server:", err)
        os.Exit(1)
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