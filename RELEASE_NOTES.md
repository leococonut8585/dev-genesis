# Dev Genesis v1.0.0 Release Notes

## 🎉 Initial Release - June 19, 2025

### ✨ Features

- **One-Click Installation**: Install Python, Node.js, Git, VS Code, WSL2, and Claude Code with a single click
- **Beautiful UI**: Modern, animated interface with real-time progress tracking
- **Smart Detection**: Automatically skips already installed tools
- **Parallel Installation**: Independent tools install simultaneously for faster setup
- **Retry Mechanism**: Automatic retry with exponential backoff for network errors
- **Cross-Platform**: Binaries available for Windows, Linux, and macOS

### 📦 What Gets Installed

| Tool | Version | Description |
|------|---------|-------------|
| Python | 3.12.x | Latest Python with pip |
| Node.js | 20.x LTS | JavaScript runtime with npm |
| Git | Latest | Version control system |
| Visual Studio Code | Latest | Code editor with extensions |
| WSL2 | Ubuntu 22.04 | Windows Subsystem for Linux |
| Claude Code | Latest | AI coding assistant |

### 🚀 Quick Start

1. Download `DevGenesisInstaller.exe` (recommended) or `dev-genesis-windows-amd64.exe`
2. Run as Administrator
3. Click the glowing "GENESIS" button
4. Wait for installation to complete (~4-6 minutes)
5. Start coding!

### 💻 System Requirements

- Windows 11 Pro (64-bit)
- 8GB RAM minimum (16GB recommended)
- 10GB free disk space
- Internet connection
- Administrator privileges

### 🔧 Advanced Usage

Run with custom port:
```bash
set PORT=9999
dev-genesis-windows-amd64.exe
```

### 🐛 Known Issues

- WSL2 installation may require a system restart
- Some antivirus software may flag the executable (false positive)

### 🙏 Acknowledgments

Created with ❤️ by Leo Sakaguchi
"Click Once, Code Forever"