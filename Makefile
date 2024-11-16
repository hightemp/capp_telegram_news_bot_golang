PROJECT_NAME=capp_telegram_news_bot_golang

build:
	@echo "Building binary"
	@go build -o $(PROJECT_NAME)

build-static:
	@echo "Building static binary"
	@CGO_ENABLED=0 go build -ldflags="-s -w -extldflags=-static" -o $(PROJECT_NAME)_static

build-docker:
	@echo "Building docker image"
	@docker build -t $(PROJECT_NAME) .

run:
	@echo "Running binary"
	@go run main.go

run-docker:
	@echo "Running docker image"
	@docker run -it --rm $(PROJECT_NAME)