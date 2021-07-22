package old

import (
	configuration "gitlab-tg-bot/conf"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	api  *tgbotapi.BotAPI
	conf Configuration
}

var _ Notifier = (*Bot)(nil)
var _ UpdateWatcher = (*Bot)(nil)

func NewBot(conf Configuration) (Notifier, error) {
	bot, err := tgbotapi.NewBotAPI(conf.GetString(configuration.Token))
	return &Bot{api: bot, conf: conf}, err
}

func (b *Bot) Notify(payload string) error {
	msgConf := tgbotapi.NewMessage(b.conf.GetInt64(configuration.ChatId), payload)
	msgConf.ParseMode = "Markdown"
	msgConf.DisableNotification = true
	msg, err := b.api.Send(msgConf)
	logrus.Debugf("В чат ID=%d отправлено сообщение: \n%v ", msgConf.ChatID, msg)
	return err
}

func (b *Bot) Polling() {
	updates := b.api.GetUpdatesChan(tgbotapi.UpdateConfig{
		Offset:  0,
		Timeout: 60,
	})

	for update := range updates {
		if update.Message == nil {
			continue
		}
	}
}
