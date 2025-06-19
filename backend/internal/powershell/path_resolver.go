package powershell

import (
    "os"
    "path/filepath"
)

// GetExecutablePath returns the absolute path of the current executable
func GetExecutablePath() (string, error) {
    ex, err := os.Executable()
    if err != nil {
        return "", err
    }
    return filepath.Abs(ex)
}

// GetScriptDirectory returns the absolute path to the scripts directory
func GetScriptDirectory() (string, error) {
    exePath, err := GetExecutablePath()
    if err != nil {
        return "", err
    }
    
    // Get the directory containing the executable
    exeDir := filepath.Dir(exePath)
    
    // Scripts are in backend/scripts/install relative to project root
    // When running from backend directory, go up one level
    scriptDir := filepath.Join(exeDir, "..", "backend", "scripts", "install")
    
    // Clean and return absolute path
    return filepath.Abs(scriptDir)
}