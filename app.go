package main

import (
	"context"
	"log"

	"github.com/playwright-community/playwright-go"
)

// App struct
type App struct {
	ctx           context.Context
	browserActive bool
	pw            *playwright.Playwright
	browser       playwright.Browser
	page          playwright.Page
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		browserActive: false,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
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
		return "启动浏览器失败: " + err.Error()
	}
	a.pw = pw

	// 启动系统Chrome浏览器
	launchOptions := playwright.BrowserTypeLaunchOptions{
		Channel: playwright.String("chrome"), // 使用系统Chrome浏览器
		Headless: playwright.Bool(false),     // 非无头模式，可以看到浏览器界面
	}

	browser, err := pw.Chromium.Launch(launchOptions)
	if err != nil {
		log.Printf("启动Chrome浏览器失败: %v", err)
		pw.Stop()
		return "启动浏览器失败: " + err.Error()
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
