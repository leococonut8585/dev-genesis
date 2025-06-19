package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "os/signal"
    "path/filepath"
    "strings"
    "syscall"
    "time"

    "github.com/gorilla/mux"
    "github.com/gorilla/websocket"
    
    wsHandler "github.com/leococonut8585/dev-genesis/internal/websocket"
    "github.com/leococonut8585/dev-genesis/internal/installer"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        origin := r.Header.Get("Origin")
        return strings.HasPrefix(origin, "http://localhost") || 
               strings.HasPrefix(origin, "http://127.0.0.1")
    },
}

func main() {
    port := "8888"
    if p := os.Getenv("PORT"); p != "" {
        port = p
    }

    router := mux.NewRouter()
    
    router.HandleFunc("/ws", handleWebSocket)
    router.HandleFunc("/api/install", handleInstall).Methods("POST")
    router.HandleFunc("/api/status", handleStatus).Methods("GET")
    
    // Serve static files from filesystem instead of embed
    router.PathPrefix("/static/").Handler(
        http.StripPrefix("/static/", 
            http.FileServer(http.Dir("./web/static"))))
    
    router.HandleFunc("/", serveIndex)
    
    srv := &http.Server{
        Addr:         ":" + port,
        Handler:      router,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
    
    go func() {
        log.Printf("🚀 Dev Genesis server starting on http://localhost:%s", port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed to start: %v", err)
        }
    }()
    
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    <-sigChan
    
    log.Println("⏹️  Shutting down server...")
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("Server forced to shutdown: %v", err)
    }
    
    log.Println("✅ Server stopped")
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./web/templates/index.html")
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket upgrade failed: %v", err)
        return
    }
    
    log.Println("✨ New WebSocket connection established")
    
    handler := wsHandler.NewHandler(conn)
    defer handler.Close()
    
    handler.SendStatus("Connected to Dev Genesis server")
    
    scriptDir := filepath.Join("scripts", "install")
    manager := installer.NewManager(handler, scriptDir)
    
    for {
        var msg wsHandler.Message
        if err := conn.ReadJSON(&msg); err != nil {
            log.Printf("WebSocket read error: %v", err)
            break
        }
        if msg.Type == wsHandler.TypeInstall {
            log.Println("🚀 Installation command received")
            go manager.InstallAll()
        }
    }
}

func handleInstall(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "status": "started",
        "message": "Installation triggered via WebSocket",
    })
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status": "ready",
        "version": "1.0.0",
        "tools": installer.Tools,
    })
}