package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "strings"

    "github.com/O-tero/internal/resp" 
)

func main() {
    listener, err := net.Listen("tcp", ":6379")
    if err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
    defer listener.Close()

    log.Println("Redis Lite server is running on port 6379...")
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("Failed to accept connection: %v", err)
            continue
        }
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    reader := bufio.NewReader(conn)

    for {
        data, err := reader.ReadBytes('\n')
        if err != nil {
            log.Printf("Connection closed: %v", err)
            return
        }

        // Deserialize the RESP message
        command, err := resp.Deserialize(data)
        if err != nil {
            sendError(conn, "ERR Invalid RESP message")
            continue
        }

        // Process the command
        response := processCommand(command)
        if response == nil {
            continue
        }

        // Serialize the response
        serialized, err := resp.Serialize(response)
        if err != nil {
            log.Printf("Failed to serialize response: %v", err)
            continue
        }

        // Send the response to the client
        conn.Write([]byte(serialized))
    }
}

func processCommand(command interface{}) interface{} {
    array, ok := command.([]interface{})
    if !ok || len(array) == 0 {
        return "ERR Invalid command"
    }

    cmd, ok := array[0].(string)
    if !ok {
        return "ERR Invalid command format"
    }

    cmd = strings.ToUpper(cmd)
    switch cmd {
    case "PING":
        if len(array) == 1 {
            return "PONG"
        } else if len(array) == 2 {
            return array[1]
        }
        return "ERR PING takes at most one argument"
    case "ECHO":
        if len(array) != 2 {
            return "ERR ECHO takes exactly one argument"
        }
        return array[1]
    default:
        return fmt.Sprintf("ERR unknown command '%s'", cmd)
    }
}

func sendError(conn net.Conn, message string) {
    serialized, _ := resp.Serialize(message)
    conn.Write([]byte(serialized))
}
