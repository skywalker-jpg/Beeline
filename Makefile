GO_BUILD_PATH ?= $(CURDIR)/bin
GO_BUILD_APP_PATH ?= $(GO_BUILD_PATH)/beeline/

GOOS ?= linux
GOARCH ?= $(shell go env GOARCH)
CGO ?= 0

# Цель для сборки
build:
	CGO_ENABLED=$(CGO) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(GO_BUILD_APP_PATH) ./cmd/beeline/

# Запуск docker compose
up:
	docker-compose up --build -d

# Остановка docker compose
down:
	@docker-compose down

# Очистка сгенерированных файлов
clean:
	@rm -rf $(GO_BUILD_PATH)