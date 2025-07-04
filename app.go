package main

import (
	"context"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/playwright-community/playwright-go"
)

// App struct
type App struct {
	ctx           context.Context
	browserActive bool
	pw            *playwright.Playwright
	browser       playwright.Browser
	page          playwright.Page
	chromePath    string // 存储Chrome浏览器路径
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		browserActive: false,
		chromePath:    getDefaultChromePath(), // 设置默认路径
	}
}

// 获取不同操作系统下Chrome浏览器的默认路径
func getDefaultChromePath() string {
	switch runtime.GOOS {
	case "windows":
		// Windows常见的Chrome路径
		paths := []string{
			"C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
			"C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe",
			"%LOCALAPPDATA%\\Google\\Chrome\\Application\\chrome.exe",
		}
		for _, path := range paths {
			// 简单替换环境变量
			if strings.Contains(path, "%LOCALAPPDATA%") {
				localAppData := os.Getenv("LOCALAPPDATA")
				path = strings.Replace(path, "%LOCALAPPDATA%", localAppData, 1)
			}
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	case "darwin":
		// macOS Chrome路径
		return "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	case "linux":
		// Linux常见的Chrome路径
		paths := []string{
			"/usr/bin/google-chrome",
			"/usr/bin/google-chrome-stable",
			"/usr/bin/chrome",
			"/usr/bin/chromium-browser",
		}
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	}
	return "" // 如果找不到默认路径，返回空字符串
}

// SetChromePath 设置Chrome浏览器的路径
func (a *App) SetChromePath(path string) string {
	if _, err := os.Stat(path); err != nil {
		return "无效的Chrome路径，文件不存在: " + path
	}
	
	a.chromePath = path
	return "已成功设置Chrome路径: " + path
}

// GetChromePath 获取当前设置的Chrome浏览器路径
func (a *App) GetChromePath() string {
	if a.chromePath == "" {
		return "未设置Chrome路径，系统将尝试使用默认路径"
	}
	return a.chromePath
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	
	// 记录当前的Chrome路径配置
	if a.chromePath != "" {
		log.Printf("应用启动：使用Chrome路径: %s", a.chromePath)
	} else {
		log.Println("应用启动：未指定Chrome路径，将使用系统默认Chrome或channel方式启动")
	}
}

// CreateBrowser 创建一个新的浏览器实例
func (a *App) CreateBrowser() string {
	// 如果已有浏览器实例，先关闭
	if a.browserActive && a.browser != nil {
		if a.page != nil {
			err := a.page.Close()
			if err != nil {
				log.Printf("关闭页面失败: %v", err)
			}
		}
		err := a.browser.Close()
		if err != nil {
			log.Printf("关闭浏览器失败: %v", err)
		}
		if a.pw != nil {
			err = a.pw.Stop()
			if err != nil {
				log.Printf("停止playwright失败: %v", err)
			}
		}
	}

	// 初始化Playwright
	pw, err := playwright.Run()
	if err != nil {
		log.Printf("初始化Playwright失败: %v", err)
		
		// 提供更友好的错误信息
		if runtime.GOOS == "windows" {
			return "启动浏览器失败：请确保已安装Playwright驱动。\n可以通过以下步骤手动安装：\n1. 打开命令提示符\n2. 运行命令：npx playwright install --with-deps chromium\n\n详细错误：" + err.Error()
		}
		return "启动浏览器失败: " + err.Error()
	}
	a.pw = pw

	// 设置启动选项
	launchOptions := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false), // 非无头模式，可以看到浏览器界面
	}
	
	// 如果指定了Chrome路径，则使用executablePath
	if a.chromePath != "" {
		log.Printf("使用指定的Chrome路径: %s", a.chromePath)
		launchOptions.ExecutablePath = playwright.String(a.chromePath)
	} else {
		// 如果没有指定路径，则尝试使用系统Chrome
		log.Println("未指定Chrome路径，尝试使用系统Chrome")
		launchOptions.Channel = playwright.String("chrome")
	}

	// 启动浏览器
	browser, err := pw.Chromium.Launch(launchOptions)
	if err != nil {
		log.Printf("启动Chrome浏览器失败: %v", err)
		pw.Stop()
		
		if a.chromePath != "" {
			return "启动Chrome浏览器失败: 指定的路径可能不正确: " + a.chromePath + "\n错误详情: " + err.Error()
		}
		return "启动Chrome浏览器失败: " + err.Error() + "\n请尝试使用SetChromePath方法手动指定Chrome路径"
	}
	a.browser = browser

	// 创建新页面
	page, err := browser.NewPage()
	if err != nil {
		log.Printf("创建新页面失败: %v", err)
		browser.Close()
		pw.Stop()
		return "创建新页面失败: " + err.Error()
	}
	a.page = page

	// 导航到百度首页
	if _, err := page.Goto("https://www.baidu.com", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		log.Printf("导航到百度失败: %v", err)
	}

	a.browserActive = true
	return "浏览器创建成功"
}

// GetCookie 从当前浏览器获取cookie
func (a *App) GetCookie() string {
	if !a.browserActive || a.page == nil {
		return "请先创建浏览器"
	}

	// 获取当前页面的cookie
	cookies, err := a.page.Context().Cookies()
	if err != nil {
		log.Printf("获取Cookie失败: %v", err)
		return "获取Cookie失败: " + err.Error()
	}

	// 将cookie格式化为字符串
	cookieStr := ""
	for i, cookie := range cookies {
		if i > 0 {
			cookieStr += "; "
		}
		cookieStr += cookie.Name + "=" + cookie.Value
	}

	if cookieStr == "" {
		return "当前页面没有Cookie"
	}

	return cookieStr
}
