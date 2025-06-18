package websocket

import (
    "encoding/json"
    "log"
    "sync"
    "time"

    "github.com/gorilla/websocket"
)

type MessageType string

const (
    TypeProgress MessageType = "progress"
    TypeStatus   MessageType = "status"
    TypeError    MessageType = "error"
    TypeComplete MessageType = "complete"
    TypePing     MessageType = "ping"
    TypePong     MessageType = "pong"
)

type Message struct {
    Type       MessageType `json:"type"`
    Percentage int         `json:"percentage,omitempty"`
    Message    string      `json:"message,omitempty"`
    Error      string      `json:"error,omitempty"`
    Timestamp  int64       `json:"timestamp"`
}

type Handler struct {
    conn      *websocket.Conn
    send      chan Message
    mu        sync.Mutex
    done      chan struct{}
    pingTicker *time.Ticker
}

func NewHandler(conn *websocket.Conn) *Handler {
    h := &Handler{
        conn:      conn,
        send:      make(chan Message, 256),
        done:      make(chan struct{}),
        pingTicker: time.NewTicker(30 * time.Second),
    }
    
    go h.writePump()
    go h.readPump()
    go h.pingPump()
    
    return h
}

func (h *Handler) Close() {
    h.pingTicker.Stop()
    close(h.done)
    h.conn.Close()
}

func (h *Handler) SendProgress(percentage int, message string) {
    h.send <- Message{
        Type:       TypeProgress,
        Percentage: percentage,
        Message:    message,
        Timestamp:  time.Now().Unix(),
    }
}

func (h *Handler) SendStatus(message string) {
    h.send <- Message{
        Type:      TypeStatus,
        Message:   message,
        Timestamp: time.Now().Unix(),
    }
}

func (h *Handler) SendError(err string) {
    h.send <- Message{
        Type:      TypeError,
        Error:     err,
        Timestamp: time.Now().Unix(),
    }
}

func (h *Handler) SendComplete(message string) {
    h.send <- Message{
        Type:      TypeComplete,
        Message:   message,
        Timestamp: time.Now().Unix(),
    }
}

func (h *Handler) readPump() {
    defer h.Close()
    
    h.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
    h.conn.SetPongHandler(func(string) error {
        h.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        return nil
    })
    
    for {
        var msg Message
        err := h.conn.ReadJSON(&msg)
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("WebSocket error: %v", err)
            }
            break
        }
        
        // Handle incoming messages
        switch msg.Type {
        case TypePing:
            h.send <- Message{Type: TypePong, Timestamp: time.Now().Unix()}
        default:
            log.Printf("Received message: %+v", msg)
        }
    }
}

func (h *Handler) writePump() {
    ticker := time.NewTicker(54 * time.Second)
    defer func() {
        ticker.Stop()
        h.Close()
    }()
    
    for {
        select {
        case message, ok := <-h.send:
            h.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if !ok {
                h.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            
            if err := h.conn.WriteJSON(message); err != nil {
                return
            }
            
        case <-ticker.C:
            h.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if err := h.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
            
        case <-h.done:
            return
        }
    }
}

func (h *Handler) pingPump() {
    for {
        select {
        case <-h.pingTicker.C:
            h.send <- Message{Type: TypePing, Timestamp: time.Now().Unix()}
        case <-h.done:
            return
        }
    }
}