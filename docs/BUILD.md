# 构建与发布指南

**📖 Languages**: [English](BUILD_EN.md) | [中文](BUILD.md)

本文档说明如何构建和发布NATS Client应用程序。

## 🔧 开发环境设置

### 必需工具
- Go 1.21 或更高版本
- Git
- Fyne打包工具

### 安装步骤
```bash
# 克隆项目
git clone https://github.com/gogap/nats-app.git
cd nats-app

# 安装依赖
make deps

# 安装Fyne打包工具
make fyne-deps
```

## 🏗️ 本地构建

### 快速构建（当前平台）
```bash
make build
```

### 开发构建（使用go build）
```bash
make build-dev
```

### 跨平台构建
```bash
make build-all
```

### 运行应用
```bash
# 开发模式运行
make run

# 运行打包的应用
# macOS:
open nats-client.app

# Linux/Windows:
./nats-client
```

## 📦 手动打包

### 基本打包
```bash
fyne package --name nats-client
```

### 带版本信息的打包
```bash
fyne package --name nats-client \
  --app-version 1.0.0 \
  --app-build 1 \
  --app-id io.github.gogap.nats-app
```

### 跨平台手动打包
```bash
# Windows
fyne package --os windows --name nats-client-windows

# macOS (需要在macOS上构建)
fyne package --os darwin --name nats-client-darwin

# Linux (需要在Linux上构建)
fyne package --os linux --name nats-client-linux
```

## 🚀 GitHub Actions 自动发布

### 发布流程
当推送tag时（如`v1.0.0`），GitHub Actions会自动：

1. **构建所有平台**：
   - Linux (Ubuntu)
   - Windows (Windows Server)
   - macOS Intel (macOS)
   - macOS Apple Silicon (macOS)

2. **生成发布文件**：
   - `nats-client-linux-amd64`
   - `nats-client-windows-amd64.exe`
   - `nats-client-darwin-amd64.zip`
   - `nats-client-darwin-arm64.zip`

3. **创建GitHub Release**并上传所有文件

### 创建发布
```bash
# 创建并推送tag
git tag v1.0.0
git push origin v1.0.0
```

### 工作流文件
- `.github/workflows/ci.yml` - 持续集成测试
- `.github/workflows/release.yml` - 自动发布

## 🛠️ 构建配置

### Makefile 目标
- `deps` - 安装Go依赖
- `fyne-deps` - 安装Fyne打包工具
- `build` - 使用fyne package构建
- `build-all` - 构建所有平台（注意：某些平台需要原生构建）
- `build-dev` - 使用go build进行开发构建
- `run` - 运行应用程序
- `clean` - 清理构建文件
- `test` - 运行测试
- `fmt` - 格式化代码
- `lint` - 代码检查

### 构建标志
应用程序在构建时会注入版本信息：
- `Version` - Git标签或"dev"
- `BuildTime` - 构建时间
- `GoVersion` - Go版本

## 📁 输出文件

### macOS
- 生成 `.app` 应用包
- 在发布时会压缩为 `.zip` 文件

### Windows
- 生成 `.exe` 可执行文件

### Linux
- 生成可执行文件（无扩展名）

## ⚠️ 注意事项

### 跨平台编译限制
由于Fyne使用CGO，跨平台编译可能遇到问题：
- 最好在目标平台上进行本地构建
- GitHub Actions在各自的原生平台上构建

### 依赖要求
- **Linux**: 需要 `libgl1-mesa-dev` 和 `xorg-dev`
- **Windows**: 无特殊要求
- **macOS**: 无特殊要求

### 文件权限
Linux和macOS的可执行文件需要执行权限：
```bash
chmod +x nats-client-linux-amd64
```

## 🔍 故障排除

### 常见问题

**Q: fyne package失败**
```bash
# 确保安装了最新版本
go install fyne.io/tools/cmd/fyne@latest
```

**Q: Linux依赖缺失**
```bash
sudo apt-get update
sudo apt-get install -y libgl1-mesa-dev xorg-dev
```

**Q: 跨平台编译失败**
- 使用原生平台构建
- 或者使用GitHub Actions自动构建

**Q: 版本信息显示"dev"**
- 确保在Git仓库中有标签
- 使用 `git tag v1.0.0` 创建标签

## 📚 相关文档

- [Fyne打包文档](https://developer.fyne.io/tutorial/packaging.html)
- [Go交叉编译](https://golang.org/doc/install/source#environment)
- [GitHub Actions工作流](https://docs.github.com/en/actions/using-workflows) 