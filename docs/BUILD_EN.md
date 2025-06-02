# Build & Release Guide

**📖 Languages**: [English](BUILD_EN.md) | [中文](BUILD.md)

This document explains how to build and release the NATS Client application.

## 🔧 Development Environment Setup

### Required Tools
- Go 1.21 or higher
- Git
- Fyne packaging tool

### Installation Steps
```bash
# Clone repository
git clone https://github.com/gogap/nats-app.git
cd nats-app

# Install dependencies
make deps

# Install Fyne packaging tool
make fyne-deps
```

## 🏗️ Local Building

### Quick Build (Current Platform)
```bash
make build
```

### Development Build (using go build)
```bash
make build-dev
```

### Cross-platform Build
```bash
make build-all
```

### Running the Application
```bash
# Development mode
make run

# Run packaged application
# macOS:
open nats-client.app

# Linux/Windows:
./nats-client
```

## 📦 Manual Packaging

### Basic Packaging
```bash
fyne package --name nats-client
```

### Packaging with Version Information
```bash
fyne package --name nats-client \
  --app-version 1.0.0 \
  --app-build 1 \
  --app-id io.github.gogap.nats-app
```

### Cross-platform Manual Packaging
```bash
# Windows
fyne package --os windows --name nats-client-windows

# macOS (must build on macOS)
fyne package --os darwin --name nats-client-darwin

# Linux (must build on Linux)
fyne package --os linux --name nats-client-linux
```

## 🚀 GitHub Actions Automated Release

### Release Process
When pushing a tag (e.g., `v1.0.0`), GitHub Actions will automatically:

1. **Build All Platforms**:
   - Linux (Ubuntu)
   - Windows (Windows Server)
   - macOS Intel (macOS)
   - macOS Apple Silicon (macOS)

2. **Generate Release Files**:
   - `nats-client-linux-amd64`
   - `nats-client-windows-amd64.exe`
   - `nats-client-darwin-amd64.zip`
   - `nats-client-darwin-arm64.zip`

3. **Create GitHub Release** and upload all files

### Creating a Release
```bash
# Create and push tag
git tag v1.0.0
git push origin v1.0.0
```

### Workflow Files
- `.github/workflows/ci.yml` - Continuous Integration Testing
- `.github/workflows/release.yml` - Automated Release

## 🛠️ Build Configuration

### Makefile Targets
- `deps` - Install Go dependencies
- `fyne-deps` - Install Fyne packaging tool
- `build` - Build using fyne package
- `build-all` - Build all platforms (note: some platforms require native builds)
- `build-dev` - Development build using go build
- `run` - Run the application
- `clean` - Clean build files
- `test` - Run tests
- `fmt` - Format code
- `lint` - Code linting

### Build Flags
Version information is injected during build:
- `Version` - Git tag or "dev"
- `BuildTime` - Build timestamp
- `GoVersion` - Go version

## 📁 Output Files

### macOS
- Generates `.app` application bundle
- Compressed to `.zip` file for release

### Windows
- Generates `.exe` executable file

### Linux
- Generates executable file (no extension)

## ⚠️ Important Notes

### Cross-platform Compilation Limitations
Due to Fyne using CGO, cross-platform compilation may encounter issues:
- Best to build on target platforms natively
- GitHub Actions builds on respective native platforms

### Dependency Requirements
- **Linux**: Requires `libgl1-mesa-dev` and `xorg-dev`
- **Windows**: No special requirements
- **macOS**: No special requirements

### File Permissions
Linux and macOS executables need execute permissions:
```bash
chmod +x nats-client-linux-amd64
```

## 🔍 Troubleshooting

### Common Issues

**Q: fyne package fails**
```bash
# Ensure latest version is installed
go install fyne.io/tools/cmd/fyne@latest
```

**Q: Linux dependencies missing**
```bash
sudo apt-get update
sudo apt-get install -y libgl1-mesa-dev xorg-dev
```

**Q: Cross-platform compilation fails**
- Use native platform builds
- Or use GitHub Actions automated builds

**Q: Version shows "dev"**
- Ensure repository has tags
- Use `git tag v1.0.0` to create tags

## 📚 Related Documentation

- [Fyne Packaging Documentation](https://developer.fyne.io/tutorial/packaging.html)
- [Go Cross Compilation](https://golang.org/doc/install/source#environment)
- [GitHub Actions Workflows](https://docs.github.com/en/actions/using-workflows) 