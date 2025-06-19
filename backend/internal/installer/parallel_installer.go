package installer

import (
    "fmt"
    "sync"
    "time"
)

// ToolGroup represents tools that can be installed in parallel
type ToolGroup struct {
    Name  string
    Tools []Tool
}

var ParallelGroups = []ToolGroup{
    {
        Name: "Core Development Tools",
        Tools: []Tool{
            {Name: "Python 3.12", ScriptName: "python.ps1", Timeout: 5 * time.Minute, Weight: 15},
            {Name: "Node.js LTS", ScriptName: "nodejs.ps1", Timeout: 5 * time.Minute, Weight: 15},
            {Name: "Git", ScriptName: "git.ps1", Timeout: 3 * time.Minute, Weight: 10},
        },
    },
    {
        Name: "IDE and Extensions",
        Tools: []Tool{
            {Name: "Visual Studio Code", ScriptName: "vscode.ps1", Timeout: 5 * time.Minute, Weight: 20},
            {Name: "Claude Code", ScriptName: "claude-code.ps1", Timeout: 3 * time.Minute, Weight: 10},
        },
    },
    {
        Name: "WSL2 Environment",
        Tools: []Tool{
            {Name: "WSL2 + Ubuntu", ScriptName: "wsl2.ps1", Timeout: 10 * time.Minute, Weight: 30},
        },
    },
}

func (m *Manager) InstallAllParallel() error {
    m.ws.SendStatus("🚀 Starting parallel Dev Genesis installation...")
    
    for _, group := range ParallelGroups {
        m.ws.SendStatus(fmt.Sprintf("📦 Installing %s...", group.Name))
        
        var wg sync.WaitGroup
        errors := make(chan error, len(group.Tools))
        
        // Install all tools in the group in parallel
        for _, tool := range group.Tools {
            wg.Add(1)
            go func(t Tool) {
                defer wg.Done()
                if err := m.installTool(t, 0); err != nil {
                    errors <- fmt.Errorf("%s: %v", t.Name, err)
                }
            }(tool)
        }
        
        // Wait for group to complete
        wg.Wait()
        close(errors)
        
        // Check for errors
        for err := range errors {
            m.ws.SendError(err.Error())
            return err
        }
        
        m.ws.SendProgress(m.calculateBaseProgress(), 
            fmt.Sprintf("✅ %s installed successfully", group.Name))
    }
    
    m.ws.SendComplete("🎉 All tools installed successfully!")
    return nil
}