FROM golang:1.16-alpine

WORKDIR /app

COPY . .

RUN go build -o migration migrations/main.go

RUN go build -o main cmd/main.go

ENV BOT_CONFPATH=conf/stage/bot_conf.yml
EXPOSE 10010

CMD ["./migration"]
CMD ["./main"]

## build stage
#FROM golang as builder
#COPY . /src
#WORKDIR /src
#
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /deploy/server/tgbot ./cmd/main.go
#
## final stage
#FROM scratch
#WORKDIR /
#COPY --from=builder /deploy/server .
#COPY --from=builder /src/conf/bot_conf.yml .
#ENV BOT_CONFPATH=bot_conf.yml
#EXPOSE 10010
#CMD  [ "./tgbot" ]
