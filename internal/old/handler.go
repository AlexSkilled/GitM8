package old

import (
	"fmt"
	"gitlab-tg-bot/conf"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"github.com/sirupsen/logrus"
)

const (
	EventHeaderKey     = "X-Gitlab-Event"
	TokenHeaderKey     = "X-Gitlab-Token"
	MergeRequestHeader = "Merge Request Hook"
	PipelineHeader     = "Pipeline Hook"
)

type Handler struct {
	conf       Configuration
	notifier   Notifier
	processors map[string]Processor
}

func NewHandler(conf Configuration, notifier Notifier) http.Handler {
	instance := Handler{conf: conf, notifier: notifier, processors: map[string]Processor{}}
	instance.processors[MergeRequestHeader] = MergeRequestProcessor{}
	instance.processors[PipelineHeader] = PipelineProcessor{}
	return instance
}

func (h Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	authenticated := h.authenticate(req.Header)
	if !authenticated {
		rw.WriteHeader(http.StatusForbidden)
		return
	}
	event := strings.TrimSpace(req.Header.Get(EventHeaderKey))
	p, ok := h.processors[event]
	if !ok {
		logrus.Tracef("Нет обработчика для заголовка %s", event)
		return
	}
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logrus.Errorf("Ошибка при чтении тела запроса: %v", err)
	}
	msg, skip, err := p.Process(b)
	if err != nil {
		logrus.Errorf("Ошибка при обработке запроса: %v", err)
		return
	}
	if skip {
		logrus.Infof("Пропускаем обработку запроса для заголовка %s", event)
		return
	}
	err = h.notifier.Notify(msg)
	if err != nil {
		logrus.Errorf("Ошибка при отправке сообщения: %v", err)
		return
	}
}

func (h Handler) authenticate(header http.Header) bool {
	if h.conf.GetBool(conf.NoAuth) {
		return true
	}
	token := header.Get(TokenHeaderKey)
	parsedT, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("mailformed token")
		}
		return []byte(h.conf.GetString(conf.SecretKey)), nil
	})
	if err != nil {
		logrus.Errorf("Попытка доступа с неправильно сформированным токеном: %s", token)
		return false
	}
	if !parsedT.Valid {
		logrus.Errorf("Попытка доступа с невалидным токеном: %s", token)
		return false
	}
	return true
}

var _ http.Handler = (*Handler)(nil)
