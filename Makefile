.PHONY: build run test clean swagger lint

APP_NAME := management-backend
ENTRY := cmd/system/main.go
BUILD_DIR := bin

## 构建
build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(ENTRY)

## 开发模式（依赖 air）
dev:
	air

## 直接运行
run:
	go run $(ENTRY)

## 全部测试
test:
	go test ./... -v -count=1

## 指定模块测试：MODULE=system
test-module:
	go test ./internal/module/$(MODULE)/... -v -count=1

## 单个测试函数：MODULE=system FUNC=TestCreateUser
test-one:
	go test ./internal/module/$(MODULE)/... -v -run $(FUNC) -count=1

## Swagger 文档生成
swagger:
	swag init -g $(ENTRY) -o api/swagger

## 格式化
fmt:
	goimports -w .
	gofmt -w .

## 静态检查
lint:
	golangci-lint run ./...

## 清理
clean:
	rm -rf $(BUILD_DIR)
	go clean -cache

## 下载依赖
deps:
	go mod download
	go mod tidy
