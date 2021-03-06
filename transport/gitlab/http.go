package gitlab

import (
	"encoding/json"
	"gitlab-tg-bot/internal"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/transport/gitlab/events"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	EventHeaderKey = "X-Gitlab-Event"
	TokenHeaderKey = "X-Gitlab-Token"
)

type Handler struct {
	events         map[string]interfaces.GitMapper
	messageSender  interfaces.TelegramMessageSender
	messageService interfaces.MessageService
}

func NewHandler(storage interfaces.ServiceStorage, tg interfaces.TelegramMessageSender) http.Handler {
	return &Handler{
		events: map[string]interfaces.GitMapper{
			events.PushHeader:         &events.Push{},
			events.MergeRequestHeader: &events.MergeRequest{},
			events.PipelineHeader:     &events.Pipeline{},
		},
		messageService: storage.MessageHandler(),
		messageSender:  tg,
	}
}

func (h *Handler) ServeHTTP(_ http.ResponseWriter, req *http.Request) {
	eventType := strings.TrimSpace(req.Header.Get(EventHeaderKey))
	event, ok := h.events[eventType]
	if !ok {
		logrus.Errorf("Ошибка при попытке спарсить заголовок %s для Gitlab хука", eventType)
	}

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logrus.Errorf("Ошибка при чтении тела запроса: %v", err)
		return
	}

	err = json.Unmarshal(b, &event)
	if err != nil {
		logrus.Errorf("Ошибка при маршалинге тела запроса: %v", err)
		return
	}

	msg, err := h.messageService.ProcessMessage(event.ToModel())
	if err != nil {
		if err == internal.ErrorNoTickets {
			return
		}
		logrus.Errorf("Ошибка при маршалинге тела запроса: %v", err)
		return
	}

	h.messageSender.SendMessages(msg)
}
