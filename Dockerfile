FROM golang as builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /deploy/server/migrations/main ./migrations/main.go &&\
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /deploy/server/main ./cmd/main.go

COPY ./conf/stage/bot_conf.yml /deploy/server/bot_conf.yml
COPY ./migrations/ /deploy/server/migrations/

FROM alpine

RUN apk update && apk add ca-certificates

WORKDIR /app
COPY --from=builder ./deploy/server/ .
ENV TGBOT_CONFPATH=bot_conf.yml

EXPOSE 10010