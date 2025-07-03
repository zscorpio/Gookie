# Gookie - Cookie获取工具

一个简单的跨平台桌面应用，用于打开浏览器并获取Cookie。

## 功能特点

- 使用系统Chrome浏览器打开网页
- 一键获取当前页面的Cookie
- 支持一键复制Cookie
- 简洁直观的用户界面

## 安装

### 下载预编译版本

从[GitHub Releases](https://github.com/您的用户名/Gookie/releases)下载适合您操作系统的最新版本：

- Windows: `Gookie-windows-amd64.exe`
- macOS (Intel): `Gookie-darwin-amd64`
- macOS (Apple Silicon): `Gookie-darwin-arm64`
- Linux: `Gookie-linux-amd64`

### 使用要求

- 必须安装Chrome浏览器

## 使用方法

1. 运行应用程序
2. 点击"新建浏览器"按钮，系统将打开Chrome浏览器
3. 在Chrome中访问您需要获取Cookie的网站
4. 返回Gookie应用，点击"获取Cookie"按钮
5. Cookie将显示在文本框中，可以点击"复制"按钮复制到剪贴板

## 从源码构建

### 前提条件

- Go 1.21或更高版本
- Node.js 18或更高版本
- Wails CLI
- 系统Chrome浏览器

### 安装依赖

```bash
# 安装Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 安装Playwright依赖
./install.sh
```

### 开发模式

```bash
# 启动开发服务器
wails dev
```

### 构建可执行文件

```bash
# 构建当前平台的生产版本
wails build
```

## 使用GitHub Actions自动构建

本项目配置了GitHub Actions自动构建流程，可以自动为多个平台构建二进制文件。

### 触发自动构建

- **创建Tag**：只有当您创建以'v'开头的Tag时（例如`v1.0.0`），才会触发自动构建并创建Release

### 创建新版本发布

1. 为您的更改添加一个新的Tag：
   ```bash
   git tag v1.0.0
   git push --tags
   ```

2. GitHub Actions将自动构建所有平台的二进制文件并创建一个新的Release

3. 构建完成后，您可以在GitHub Releases页面编辑发布说明
