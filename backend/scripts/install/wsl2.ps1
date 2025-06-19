# WSL2 + Ubuntu Installation Script
param()

Write-Host "Starting WSL2 + Ubuntu installation..." -ForegroundColor Green
Write-Progress -Activity "Installing WSL2" -Status "0%" -PercentComplete 0

# Check if WSL2 is already installed
$wslStatus = wsl --status 2>$null
if ($wslStatus) {
    Write-Host "WSL2 is already installed" -ForegroundColor Yellow
    Write-Progress -Activity "Installing WSL2" -Status "100%" -PercentComplete 100
    exit 0
}

Write-Progress -Activity "Installing WSL2" -Status "10%" -PercentComplete 10

try {
    # Enable WSL feature
    Write-Host "Enabling WSL feature..." -ForegroundColor Cyan
    dism.exe /online /enable-feature /featurename:Microsoft-Windows-Subsystem-Linux /all /norestart
    
    Write-Progress -Activity "Installing WSL2" -Status "30%" -PercentComplete 30
    
    # Enable Virtual Machine feature
    Write-Host "Enabling Virtual Machine feature..." -ForegroundColor Cyan
    dism.exe /online /enable-feature /featurename:VirtualMachinePlatform /all /norestart
    
    Write-Progress -Activity "Installing WSL2" -Status "50%" -PercentComplete 50
    
    # Install WSL2 with Ubuntu
    Write-Host "Installing WSL2 with Ubuntu..." -ForegroundColor Cyan
    wsl --install -d Ubuntu-22.04 --no-launch
    
    Write-Progress -Activity "Installing WSL2" -Status "70%" -PercentComplete 70
    
    # Set WSL2 as default
    wsl --set-default-version 2
    
    Write-Progress -Activity "Installing WSL2" -Status "90%" -PercentComplete 90
    
    Write-Host "WSL2 + Ubuntu installation completed successfully!" -ForegroundColor Green
    Write-Host "Note: A system restart may be required to complete the installation" -ForegroundColor Yellow
    
    Write-Progress -Activity "Installing WSL2" -Status "100%" -PercentComplete 100
} catch {
    Write-Error "Failed to install WSL2: $_"
    exit 1
}