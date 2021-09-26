FROM golang as builder

# RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /deploy/server/migration/main ./migrations/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /deploy/server/main ./cmd/main.go

COPY ./conf/stage/bot_conf.yml /deploy/server/bot_conf.yml
COPY ./migrations/ /deploy/server/migration/

FROM scratch

WORKDIR /
COPY --from=builder ./deploy/server/ .
ENV TGBOT_CONFPATH=bot_conf.yml

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 10010