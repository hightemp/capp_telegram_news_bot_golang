#!/bin/bash

CFILE="capp_telegram_news_bot_golang"

timestamp=$(date +%s)
VERSION=$(echo `cat VERSION`.$timestamp)

gh release create $VERSION -t $VERSION -n "" $CFILE ${CFILE}_static
