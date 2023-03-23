docker run --rm -v "$PWD":/usr/src/app -w /usr/src/app golang:1.20 go build -v
sudo chmown $USER:$USER capp_telegram_news_bot_golang