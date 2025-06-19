package powershell

import (
    "fmt"
    "log"
    "time"
)

type RetryConfig struct {
    MaxAttempts int
    InitialDelay time.Duration
    MaxDelay time.Duration
    Multiplier float64
}

var DefaultRetryConfig = RetryConfig{
    MaxAttempts: 3,
    InitialDelay: 5 * time.Second,
    MaxDelay: 30 * time.Second,
    Multiplier: 2.0,
}

func ExecuteWithRetry(executor *Executor, config RetryConfig) error {
    var lastErr error
    delay := config.InitialDelay
    
    for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
        log.Printf("Attempt %d/%d for %s", attempt, config.MaxAttempts, executor.scriptPath)
        
        err := executor.Execute()
        if err == nil {
            return nil
        }
        
        lastErr = err
        
        if attempt < config.MaxAttempts {
            log.Printf("Failed attempt %d: %v. Retrying in %v...", attempt, err, delay)
            time.Sleep(delay)
            
            // Exponential backoff
            delay = time.Duration(float64(delay) * config.Multiplier)
            if delay > config.MaxDelay {
                delay = config.MaxDelay
            }
        }
    }
    
    return fmt.Errorf("failed after %d attempts: %v", config.MaxAttempts, lastErr)
}