package transport

import (
	configuration "gitlab-tg-bot/conf"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/transport/gitlab"
	"net/http"
)

const Gitlab = "gitlab" // на этот URI будут регистрироваться хуки для гита.
// т.е гит будет слать хуки на http://our.service.web/gitlab

type PublicHandler struct {
	conf       interfaces.Configuration
	notifier   interfaces.TelegramWorker
	processors map[string]interfaces.PublicProcessor
	//services   service.Storage
}

func ServeHTTP(conf interfaces.Configuration, services interfaces.ServiceStorage,
	notifier interfaces.TelegramWorker) {
	//instance := &PublicHandler{
	//	conf:       conf,
	//	notifier:   notifier,
	//	processors: map[string]interfaces.PublicProcessor{},
	//}
	// TODO В пакете сделать функцию, для получения мапы всех обработчиков
	//instance.processors[MergeRequestHeader] = &processors.MergeRequest{}

	http.Handle(Gitlab, gitlab.NewHandler(services.Announcer()))
	http.ListenAndServe("localhost"+conf.GetString(configuration.ServerUrl), nil)
}
