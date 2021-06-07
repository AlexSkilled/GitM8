# build stage
FROM golang as builder
COPY . /src
WORKDIR /src

ENV GO111MODULE=on
ENV GOPRIVATE=gitlab.tn.ru

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /deploy/server/tgbot ./cmd/main.go

# final stage
FROM scratch
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /deploy/server .
COPY --from=builder /src/conf/bot_conf.yml .
ENV BOT_CONFPATH=bot_conf.yml
EXPOSE 10010
CMD  [ "./tgbot" ]
