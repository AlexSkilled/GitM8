package transport

import (
	configuration "gitlab-tg-bot/conf"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"
	"gitlab-tg-bot/transport/gitlab"
	"net/http"
)

type PublicHandler struct {
	conf       interfaces.Configuration
	processors map[string]interfaces.PublicProcessor
	//services   service.Storage
}

func ServeHTTP(conf interfaces.Configuration, services interfaces.ServiceStorage,
	bot interfaces.TelegramMessageSender) {

	http.Handle(model.Gitlab.GetUri(), gitlab.NewHandler(services, bot))

	http.ListenAndServe(conf.GetString(configuration.ServerUrl), nil)
}
