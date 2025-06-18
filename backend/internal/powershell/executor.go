package powershell

import (
    "bufio"
    "context"
    "fmt"
    "log"
    "os/exec"
    "strings"
    "sync"
    "time"
)

type Executor struct {
    scriptPath string
    args       []string
    timeout    time.Duration
    output     chan string
    errors     chan error
    mu         sync.Mutex
    cmd        *exec.Cmd
    cancel     context.CancelFunc
}

func NewExecutor(scriptPath string, args []string, timeout time.Duration) *Executor {
    return &Executor{
        scriptPath: scriptPath,
        args:       args,
        timeout:    timeout,
        output:     make(chan string, 100),
        errors:     make(chan error, 10),
    }
}

func (e *Executor) Execute() error {
    ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
    e.cancel = cancel
    defer cancel()
    
    // Build PowerShell command
    psArgs := []string{
        "-NoProfile",
        "-ExecutionPolicy", "Bypass",
        "-File", e.scriptPath,
    }
    psArgs = append(psArgs, e.args...)
    
    e.cmd = exec.CommandContext(ctx, "powershell.exe", psArgs...)
    
    // Get stdout pipe
    stdout, err := e.cmd.StdoutPipe()
    if err != nil {
        return fmt.Errorf("failed to get stdout pipe: %w", err)
    }
    
    // Get stderr pipe
    stderr, err := e.cmd.StderrPipe()
    if err != nil {
        return fmt.Errorf("failed to get stderr pipe: %w", err)
    }
    
    // Start the command
    if err := e.cmd.Start(); err != nil {
        return fmt.Errorf("failed to start PowerShell: %w", err)
    }
    
    // Read stdout
    go func() {
        scanner := bufio.NewScanner(stdout)
        for scanner.Scan() {
            line := scanner.Text()
            e.output <- line
            
            // Parse progress if line contains percentage
            if strings.Contains(line, "%") {
                log.Printf("Progress: %s", line)
            }
        }
        if err := scanner.Err(); err != nil {
            e.errors <- fmt.Errorf("stdout scanner error: %w", err)
        }
    }()
    
    // Read stderr
    go func() {
        scanner := bufio.NewScanner(stderr)
        for scanner.Scan() {
            line := scanner.Text()
            e.errors <- fmt.Errorf("PowerShell error: %s", line)
        }
    }()
    
    // Wait for completion
    go func() {
        if err := e.cmd.Wait(); err != nil {
            if ctx.Err() == context.DeadlineExceeded {
                e.errors <- fmt.Errorf("script timed out after %v", e.timeout)
            } else {
                e.errors <- fmt.Errorf("PowerShell execution failed: %w", err)
            }
        }
        close(e.output)
        close(e.errors)
    }()
    
    return nil
}

func (e *Executor) Stop() error {
    e.mu.Lock()
    defer e.mu.Unlock()
    
    if e.cancel != nil {
        e.cancel()
    }
    
    if e.cmd != nil && e.cmd.Process != nil {
        return e.cmd.Process.Kill()
    }
    
    return nil
}

func (e *Executor) GetOutput() <-chan string {
    return e.output
}

func (e *Executor) GetErrors() <-chan error {
    return e.errors
}

// ParseProgress attempts to extract percentage from PowerShell output
func ParseProgress(line string) (int, bool) {
    // Look for patterns like "45%" or "Progress: 45%"
    parts := strings.Fields(line)
    for _, part := range parts {
        if strings.HasSuffix(part, "%") {
            percentStr := strings.TrimSuffix(part, "%")
            var percent int
            if _, err := fmt.Sscanf(percentStr, "%d", &percent); err == nil {
                return percent, true
            }
        }
    }
    return 0, false
}