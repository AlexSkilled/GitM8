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
	notifier   interfaces.TelegramWorker
	processors map[string]interfaces.PublicProcessor
	//services   service.Storage
}

func ServeHTTP(conf interfaces.Configuration, services interfaces.ServiceStorage,
	bot interfaces.TelegramWorker) {

	http.Handle(model.Gitlab.GetUri(), gitlab.NewHandler())

	http.ListenAndServe(conf.GetString(configuration.ServerUrl), nil)
}
