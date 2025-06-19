# Node.js LTS Installation Script
param()

Write-Host "Starting Node.js LTS installation..." -ForegroundColor Green
Write-Progress -Activity "Installing Node.js" -Status "0%" -PercentComplete 0

# Check if Node.js is already installed
$nodeVersion = node --version 2>$null
if ($nodeVersion -match "v20") {
    Write-Host "Node.js LTS is already installed" -ForegroundColor Yellow
    Write-Progress -Activity "Installing Node.js" -Status "100%" -PercentComplete 100
    exit 0
}

Write-Progress -Activity "Installing Node.js" -Status "25%" -PercentComplete 25

# Install using winget
try {
    winget install --id OpenJS.NodeJS.LTS --source winget --silent --accept-source-agreements --accept-package-agreements
    
    Write-Progress -Activity "Installing Node.js" -Status "75%" -PercentComplete 75
    
    # Verify installation
    Start-Sleep -Seconds 2
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
    
    if (node --version 2>$null) {
        Write-Host "Node.js installation completed successfully!" -ForegroundColor Green
        Write-Progress -Activity "Installing Node.js" -Status "100%" -PercentComplete 100
    } else {
        throw "Node.js installation verification failed"
    }
} catch {
    Write-Error "Failed to install Node.js: $_"
    exit 1
}