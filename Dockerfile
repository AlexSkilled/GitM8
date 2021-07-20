FROM golang:1.16-alpine

WORKDIR /app

COPY . .

RUN go build -o main cmd/main.go

ENV BOT_CONFPATH=conf/stage/bot_conf.yml
EXPOSE 10010

CMD ["./main"]
