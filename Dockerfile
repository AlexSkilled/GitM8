FROM golang as builder

# RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY . .

RUN apt install ca-certificates && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /deploy/server/migrations/main ./migrations/main.go &&\
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /deploy/server/main ./cmd/main.go

COPY ./conf/stage/bot_conf.yml /deploy/server/bot_conf.yml
COPY ./migrations/ /deploy/server/migrations/

FROM scratch

WORKDIR /app
COPY --from=builder ./deploy/server/ .
ENV TGBOT_CONFPATH=bot_conf.yml

EXPOSE 10010