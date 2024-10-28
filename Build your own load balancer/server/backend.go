package main

import (
    "fmt"
    "log"
    "net/http"
)

//initializes a separate server on the given port with a unique message.
func startBackendServer(port string, message string) {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "<html><body>%s</body></html>", message)
        log.Printf("Handled request on port %s\n", port)
    })

    server := &http.Server{
        Addr:    ":" + port,
        Handler: mux,
    }

    log.Printf("Starting backend server on port %s...\n", port)
    log.Fatal(server.ListenAndServe())
}

func main() {
    go startBackendServer("8080", "Hello from NBA server 8080")
    startBackendServer("8081", "Hello from basketball server 8081") 
}
