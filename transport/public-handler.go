package transport

import (
	"gitlab-tg-bot/internal"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service"
	"gitlab-tg-bot/transport/processors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	EventHeaderKey = "X-Gitlab-Event"
	TokenHeaderKey = "X-Gitlab-Token"

	MergeRequestHeader = "Merge Request Hook"
	PipelineHeader     = "Pipeline Hook"
)

type PublicHandler struct {
	conf       internal.Configuration
	notifier   interfaces.TelegramWorker
	processors map[string]interfaces.PublicProcessor
	services   service.Storage
}

func NewPublicHandler(conf internal.Configuration, notifier interfaces.TelegramWorker) http.Handler {
	instance := &PublicHandler{
		conf:       conf,
		notifier:   notifier,
		processors: map[string]interfaces.PublicProcessor{},
	}
	// TODO В пакете сделать функцию, для получения мапы всех обработчиков
	instance.processors[MergeRequestHeader] = &processors.MergeRequest{}

	return instance
}

func (h *PublicHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	event := strings.TrimSpace(req.Header.Get(EventHeaderKey))
	processor, ok := h.processors[event]
	if !ok {
		logrus.Tracef("Нет обработчика для заголовка %s", event)
		return
	}

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logrus.Errorf("Ошибка при чтении тела запроса: %v", err)
	}

	msg, err := processor.Process(b)
	if err != nil {
		logrus.Errorf("Ошибка при обработке запроса: %v", err)
		return
	}
	// TODO заполнять из сервиса
	var chatIds []int32

	h.notifier.SendMessage(chatIds, msg)
}
