# Makefile для telegram-bot-go

APP_NAME=telegram-bot-go
BUILD_TARGET=cmd/main.go
OUTPUT_BIN=bin/$(APP_NAME)
REMOTE_PATH=/home/youruser/$(APP_NAME)
SYSTEMD_SERVICE=telegram-bot-go.service

.PHONY: build deploy restart

# Локальная сборка (под текущую ОС)
build:
	go build -o $(OUTPUT_BIN) $(BUILD_TARGET)

# Кросс-компиляция для VPS (Linux amd64)
build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_BIN) $(BUILD_TARGET)

# Отправка бинаря и .env на VPS (замени user@host)
deploy: build-linux
	scp $(OUTPUT_BIN) .env user@your-vps:/home/youruser/$(APP_NAME)
	ssh user@your-vps "sudo systemctl restart $(SYSTEMD_SERVICE)"

# Перезапуск systemd без пересборки
restart:
	ssh user@your-vps "sudo systemctl restart $(SYSTEMD_SERVICE)"
