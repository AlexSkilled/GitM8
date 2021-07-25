package gitlab

import (
	"encoding/json"
	"fmt"
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
	events map[string]interfaces.HttpProcessor
}

func NewHandler() http.Handler {
	return &Handler{
		events: map[string]interfaces.HttpProcessor{
			events.PushHeader:         &events.Push{},
			events.MergeRequestHeader: &events.MergeRequest{},
			events.PipelineHeader:     &events.Pipeline{},
		},
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
	}

	err = json.Unmarshal(b, &event)
	if err != nil {
		logrus.Errorf("Ошибка при маршалинге тела запроса: %v", err)
	}
	msg, err := event.Process() // передать модель в сервис публикации
	if err != nil {
		logrus.Errorf("Ошибка при маршалинге тела запроса: %v", err)
	}
	fmt.Println(msg)
}
