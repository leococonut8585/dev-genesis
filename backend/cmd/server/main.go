package main

import (
    "embed"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

//go:embed all:../../web/static all:../../web/templates
var webUI embed.FS

func main() {
    port := "8888"
    if p := os.Getenv("PORT"); p != "" {
        port = p
    }

    log.Printf("Starting Dev Genesis server on port %s", port)
    
    // TODO: Implement server
    fmt.Println("Dev Genesis - Coming Soon!")
    
    // Graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    <-sigChan
    
    log.Println("Shutting down...")
}