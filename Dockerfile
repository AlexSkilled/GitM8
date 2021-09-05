FROM golang

WORKDIR /app

COPY . .
ENV BOT_CONFPATH=conf/stage/bot_conf.yml

RUN go build -o main cmd/main.go

EXPOSE 10010
