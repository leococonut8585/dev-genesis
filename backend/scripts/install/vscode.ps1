# Visual Studio Code Installation Script
param()

Write-Host "Starting Visual Studio Code installation..." -ForegroundColor Green
Write-Progress -Activity "Installing VS Code" -Status "0%" -PercentComplete 0

# Check if VS Code is already installed
$codePath = Get-Command code -ErrorAction SilentlyContinue
if ($codePath) {
    Write-Host "VS Code is already installed" -ForegroundColor Yellow
    Write-Progress -Activity "Installing VS Code" -Status "100%" -PercentComplete 100
    exit 0
}

Write-Progress -Activity "Installing VS Code" -Status "25%" -PercentComplete 25

# Install using winget
try {
    winget install --id Microsoft.VisualStudioCode --source winget --silent --accept-source-agreements --accept-package-agreements
    
    Write-Progress -Activity "Installing VS Code" -Status "75%" -PercentComplete 75
    
    # Verify installation
    Start-Sleep -Seconds 3
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
    
    if (Get-Command code -ErrorAction SilentlyContinue) {
        Write-Host "VS Code installation completed successfully!" -ForegroundColor Green
        Write-Progress -Activity "Installing VS Code" -Status "100%" -PercentComplete 100
        
        # Install recommended extensions
        Write-Host "Installing recommended extensions..." -ForegroundColor Cyan
        code --install-extension ms-python.python
        code --install-extension dbaeumer.vscode-eslint
        code --install-extension esbenp.prettier-vscode
    } else {
        throw "VS Code installation verification failed"
    }
} catch {
    Write-Error "Failed to install VS Code: $_"
    exit 1
}