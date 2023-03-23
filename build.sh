mv .git _git
docker run --rm -v "$PWD":/usr/src/app -w /usr/src/app golang:1.20 go build -v
mv _git .git
sudo chown $USER:$USER capp_telegram_news_bot_golang