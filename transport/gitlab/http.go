package gitlab

import (
	"gitlab-tg-bot/internal/interfaces"
	"net/http"
)

const (
	EventHeaderKey = "X-Gitlab-Event"
	TokenHeaderKey = "X-Gitlab-Token"
)

type Handler struct {
	announcer interfaces.AnnouncerService
}

func NewHandler(announcer interfaces.AnnouncerService) http.Handler {
	return &Handler{
		announcer: announcer,
	}
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	//event := strings.TrimSpace(req.Header.Get(EventHeaderKey))
	//processor, ok := h.processors[event]
	//if !ok {
	//	logrus.Tracef("Нет обработчика для заголовка %s", event)
	//	return
	//}
	//
	//b, err := ioutil.ReadAll(req.Body)
	//if err != nil {
	//	logrus.Errorf("Ошибка при чтении тела запроса: %v", err)
	//}
	//
	//msg, err := processor.Process(b)
	//if err != nil {
	//	logrus.Errorf("Ошибка при обработке запроса: %v", err)
	//	return
	//}
	//// TODO заполнять из сервиса
	//var chatIds []int32
	//
	//h.notifier.SendMessage(chatIds, msg)
}
