#!/bin/bash

echo "安装 Playwright 依赖..."
go run github.com/playwright-community/playwright-go/cmd/playwright install --with-deps chromium

echo "依赖安装完成" 