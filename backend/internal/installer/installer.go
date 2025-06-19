package installer

import (
    "fmt"
    "log"
    "path/filepath"
    "sync"
    "time"

    "github.com/leococonut8585/dev-genesis/internal/powershell"
    "github.com/leococonut8585/dev-genesis/internal/websocket"
)

type Tool struct {
    Name       string
    ScriptName string
    Timeout    time.Duration
    Weight     int // Progress weight (0-100)
}

var Tools = []Tool{
    {Name: "Python 3.12", ScriptName: "python.ps1", Timeout: 5 * time.Minute, Weight: 15},
    {Name: "Node.js LTS", ScriptName: "nodejs.ps1", Timeout: 5 * time.Minute, Weight: 15},
    {Name: "Git", ScriptName: "git.ps1", Timeout: 3 * time.Minute, Weight: 10},
    {Name: "Visual Studio Code", ScriptName: "vscode.ps1", Timeout: 5 * time.Minute, Weight: 20},
    {Name: "WSL2 + Ubuntu", ScriptName: "wsl2.ps1", Timeout: 10 * time.Minute, Weight: 30},
    {Name: "Claude Code", ScriptName: "claude-code.ps1", Timeout: 3 * time.Minute, Weight: 10},
}

type Manager struct {
    ws           *websocket.Handler
    scriptDir    string
    mu           sync.Mutex
    totalWeight  int
    currentWeight int
    executors    map[string]*powershell.Executor
}

func NewManager(ws *websocket.Handler, scriptDir string) *Manager {
    totalWeight := 0
    for _, tool := range Tools {
        totalWeight += tool.Weight
    }
    
    return &Manager{
        ws:            ws,
        scriptDir:     scriptDir,
        totalWeight:   totalWeight,
        currentWeight: 0,
        executors:     make(map[string]*powershell.Executor),
    }
}

func (m *Manager) InstallAll() error {
    m.ws.SendStatus("[START] Starting Dev Genesis installation...")
    
    for i, tool := range Tools {
        if err := m.installTool(tool, i); err != nil {
            m.ws.SendError(fmt.Sprintf("Failed to install %s: %v", tool.Name, err))
            return err
        }
    }
    
    m.ws.SendComplete("[SUCCESS] All tools installed successfully!")
    return nil
}

func (m *Manager) installTool(tool Tool, index int) error {
    m.ws.SendStatus(fmt.Sprintf("Installing %s...", tool.Name))
    
    scriptPath := filepath.Join(m.scriptDir, tool.ScriptName)
    executor := powershell.NewExecutor(scriptPath, nil, tool.Timeout)
    
    m.mu.Lock()
    m.executors[tool.Name] = executor
    m.mu.Unlock()
    
    if err := executor.Execute(); err != nil {
        return err
    }
    
    baseProgress := m.calculateBaseProgress()
    
    // Monitor output
    go func() {
        for line := range executor.GetOutput() {
            log.Printf("[%s] %s", tool.Name, line)
            
            if percent, ok := powershell.ParseProgress(line); ok {
                toolProgress := (percent * tool.Weight) / 100
                totalProgress := baseProgress + toolProgress
                m.ws.SendProgress(totalProgress, fmt.Sprintf("Installing %s... (%d%%)", tool.Name, percent))
            }
        }
    }()
    
    // Monitor errors
    go func() {
        for err := range executor.GetErrors() {
            log.Printf("[%s] Error: %v", tool.Name, err)
        }
    }()
    
    // Wait for completion
    time.Sleep(1 * time.Second) // Simulate work
    
    m.mu.Lock()
    m.currentWeight += tool.Weight
    delete(m.executors, tool.Name)
    m.mu.Unlock()
    
    finalProgress := m.calculateBaseProgress()
    m.ws.SendProgress(finalProgress, fmt.Sprintf("[OK] %s installed", tool.Name))
    
    return nil
}

func (m *Manager) calculateBaseProgress() int {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    if m.totalWeight == 0 {
        return 0
    }
    
    return (m.currentWeight * 100) / m.totalWeight
}

func (m *Manager) Stop() {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    for name, executor := range m.executors {
        log.Printf("Stopping installation of %s", name)
        executor.Stop()
    }
}