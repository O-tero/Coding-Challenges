package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("Received request from %s\n", r.RemoteAddr)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hello From Backend Server"))
}

func main() {
    http.HandleFunc("/", handler)
    fmt.Println("Backend server listening on port 8081...")
    if err := http.ListenAndServe(":8081", nil); err != nil {
        panic(err)
    }
}
