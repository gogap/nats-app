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
# Clone project
git clone https://github.com/gogap/nats-app.git
cd nats-app

# Install dependencies
make deps
```

## 🏗️ Local Build

### Quick Build (Current Platform)
```bash
make build
```

### Development Build (using go build)
```bash
make build-dev
```

### Run Application
```bash
# Development mode
make run

# Run packaged application
# macOS:
open nats-client.app

# Linux/Windows:
./nats-client
```

## 🎨 Font Support

This application has built-in Chinese font support using Go embed technology to directly embed the `SourceHanSansCN-Medium.otf` font into the executable:

- **Font**: Source Han Sans CN Medium weight
- **Size**: ~8MB
- **Implementation**: Go embed (no additional packaging steps required)
- **Advantage**: Automatic Chinese display support without system font installation

The font is automatically embedded during build, no additional steps required.

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

# macOS (requires building on macOS)
fyne package --os darwin --name nats-client-darwin

# Linux (requires building on Linux)
fyne package --os linux --name nats-client-linux
```

## 🚀 GitHub Actions Auto Release

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

### Create Release
```bash
# Create and push tag
git tag v1.0.0
git push origin v1.0.0
```

### Workflow Files
- `.github/workflows/release.yml` - Auto release

## 🛠️ Build Configuration

### Makefile Targets
- `deps` - Install Go dependencies
- `build` - Build application (font auto-embedded)
- `build-release` - Build release version
- `run` - Run application
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

### File Size Estimates
- **Windows**: ~60MB (including 8MB font)
- **macOS**: ~55MB (including 8MB font)
- **Linux**: ~60MB (including 8MB font)

### Platform-specific Output
- **macOS**: Generates `.app` bundle, compressed to `.zip` for release
- **Windows**: Generates `.exe` executable
- **Linux**: Generates executable (no extension)

## ⚠️ Notes

### Cross-platform Compilation Limitations
Due to Fyne using CGO, cross-platform compilation may encounter issues:
- Best to build natively on target platform
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

**Q: Build fails with font file not found**
```bash
# Ensure font file exists
ls fonts/SourceHanSansCN-Medium.otf
# If missing, download Source Han Sans Medium from official site
```

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
- Or use GitHub Actions auto-build

**Q: Version shows "dev"**
- Ensure Git repository has tags
- Use `git tag v1.0.0` to create tag

**Q: Chinese characters display incorrectly**
- Font is auto-embedded, no additional configuration needed
- If issues persist, check for system-level font rendering problems

## 📚 Related Documentation

- [Fyne Packaging Documentation](https://developer.fyne.io/tutorial/packaging.html)
- [Go embed Documentation](https://pkg.go.dev/embed)
- [GitHub Actions Workflows](https://docs.github.com/en/actions/using-workflows)
- [Source Han Sans Official](https://source.typekit.com/source-han-sans/) 