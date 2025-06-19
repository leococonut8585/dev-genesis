# Claude Code Installation Script
param()

Write-Host "Starting Claude Code installation..." -ForegroundColor Green
Write-Progress -Activity "Installing Claude Code" -Status "0%" -PercentComplete 0

try {
    # Check if npm is available
    $npmVersion = npm --version 2>$null
    if (-not $npmVersion) {
        throw "npm is not installed. Please install Node.js first."
    }
    
    Write-Progress -Activity "Installing Claude Code" -Status "25%" -PercentComplete 25
    
    # Install Claude Code globally
    Write-Host "Installing Claude Code via npm..." -ForegroundColor Cyan
    npm install -g @anthropic-ai/claude-code
    
    Write-Progress -Activity "Installing Claude Code" -Status "75%" -PercentComplete 75
    
    # Verify installation
    $claudeCodeVersion = claude-code --version 2>$null
    if ($claudeCodeVersion) {
        Write-Host "Claude Code installation completed successfully!" -ForegroundColor Green
        Write-Progress -Activity "Installing Claude Code" -Status "100%" -PercentComplete 100
    } else {
        throw "Claude Code installation verification failed"
    }
} catch {
    Write-Error "Failed to install Claude Code: $_"
    exit 1
}