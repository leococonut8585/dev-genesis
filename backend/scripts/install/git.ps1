# Git Installation Script
param()

Write-Host "Starting Git installation..." -ForegroundColor Green
Write-Progress -Activity "Installing Git" -Status "0%" -PercentComplete 0

# Check if Git is already installed
$gitVersion = git --version 2>$null
if ($gitVersion) {
    Write-Host "Git is already installed" -ForegroundColor Yellow
    Write-Progress -Activity "Installing Git" -Status "100%" -PercentComplete 100
    exit 0
}

Write-Progress -Activity "Installing Git" -Status "25%" -PercentComplete 25

# Install using winget
try {
    winget install --id Git.Git --source winget --silent --accept-source-agreements --accept-package-agreements
    
    Write-Progress -Activity "Installing Git" -Status "75%" -PercentComplete 75
    
    # Verify installation
    Start-Sleep -Seconds 2
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
    
    if (git --version 2>$null) {
        Write-Host "Git installation completed successfully!" -ForegroundColor Green
        Write-Progress -Activity "Installing Git" -Status "100%" -PercentComplete 100
    } else {
        throw "Git installation verification failed"
    }
} catch {
    Write-Error "Failed to install Git: $_"
    exit 1
}