package internal

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	api  *tgbotapi.BotAPI
	conf Configuration
}

var _ Notifier = (*Bot)(nil)

func NewNotifier(conf Configuration) (Notifier, error) {
	bot, err := tgbotapi.NewBotAPI(conf.GetString(Token))
	return &Bot{api: bot, conf: conf}, err
}

func (b *Bot) Notify(payload string) error {
	msgConf := tgbotapi.NewMessage(b.conf.GetInt64(ChatId), payload)
	msgConf.ParseMode = "Markdown"
	msg, err := b.api.Send(msgConf)
	logrus.Debugf("В чат ID=%d отправлено сообщение: \n%v ", msgConf.ChatID, msg)
	return err
}
