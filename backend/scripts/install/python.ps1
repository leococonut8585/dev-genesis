# Python 3.12 Installation Script
param()

Write-Host "Starting Python 3.12 installation..." -ForegroundColor Green
Write-Progress -Activity "Installing Python" -Status "0%" -PercentComplete 0

# Check if Python is already installed
$pythonVersion = python --version 2>$null
if ($pythonVersion -match "Python 3\.12") {
    Write-Host "Python 3.12 is already installed" -ForegroundColor Yellow
    Write-Progress -Activity "Installing Python" -Status "100%" -PercentComplete 100
    exit 0
}

Write-Progress -Activity "Installing Python" -Status "25%" -PercentComplete 25

# Install using winget
try {
    winget install --id Python.Python.3.12 --source winget --silent --accept-source-agreements --accept-package-agreements
    
    Write-Progress -Activity "Installing Python" -Status "75%" -PercentComplete 75
    
    # Verify installation
    Start-Sleep -Seconds 2
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
    
    if (python --version 2>$null) {
        Write-Host "Python installation completed successfully!" -ForegroundColor Green
        Write-Progress -Activity "Installing Python" -Status "100%" -PercentComplete 100
    } else {
        throw "Python installation verification failed"
    }
} catch {
    Write-Error "Failed to install Python: $_"
    exit 1
}