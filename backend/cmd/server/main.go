package main

import (
    "context"
    "embed"
    "encoding/json"
    "io/fs"
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

//go:embed all:../../web/static all:../../web/templates
var webUI embed.FS

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
    
    staticFS, err := fs.Sub(webUI, "web/static")
    if err != nil {
        log.Fatal(err)
    }
    router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))
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
    indexHTML, err := webUI.ReadFile("web/templates/index.html")
    if err != nil {
        http.Error(w, "Failed to load index.html", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.Write(indexHTML)
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
    
    // Create installer manager
    scriptDir := filepath.Join(filepath.Dir(os.Args[0]), "scripts", "install")
    manager := installer.NewManager(handler, scriptDir)
    
    // Wait for install command instead of auto-starting
    go func() {
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
    }()
    
    // Wait for handler completion
    select {}
}

func handleInstall(w http.ResponseWriter, r *http.Request) {
    // TODO: Trigger installation via WebSocket
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