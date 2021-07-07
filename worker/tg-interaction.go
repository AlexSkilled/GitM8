package worker

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gitlab-tg-bot/internal"
	"gitlab-tg-bot/internal/interfaces"
	"log"
)

const (
	CommandRegister  string = "/register"
	CommandStart     string = "/start"
	CommandSubscribe string = "/subscribe"
)

type Worker struct {
	interfaces.ServiceStorage

	bot         *tgbotapi.BotAPI
	processors  map[string]internal.TgProcessor
	conf        internal.Configuration
	interceptor internal.TgProcessor
}

func NewTelegramWorker(conf internal.Configuration,
	services interfaces.ServiceStorage) interfaces.TelegramWorker {
	bot, err := tgbotapi.NewBotAPI(conf.GetString(internal.Token))
	if err != nil {
		panic(err)
	}

	bot.Debug = conf.GetBool(internal.Debug)

	log.Printf("Авторизация в боте %s", bot.Self.UserName)
	processors := map[string]internal.TgProcessor{}
	processors[CommandRegister] = &Register{}

	return &Worker{
		bot:            bot,
		ServiceStorage: services,
		processors:     processors,
		conf:           conf,
	}
}

func (t *Worker) Start() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updChan := t.bot.GetUpdatesChan(updateConfig)

	for update := range updChan {
		if update.Message == nil {
			continue
		}

		if t.interceptor != nil {
			if t.interceptor.Process(update) {
				t.interceptor = nil
			}
		}

		if processor, ok := t.processors[update.Message.Text]; ok {
			if processor.IsInterceptor() {
				t.interceptor = processor
			}
			processor.Process(update)
			continue
		}
		logrus.Printf("Не знаю как обработать команду - %s", update.Message.Text)
	}
}

func (t *Worker) SendMessage(chatIds []int32, msg string) {
	for _, id := range chatIds {
		msgConf := tgbotapi.NewMessage(int64(id), msg)
		msgConf.ParseMode = "Markdown"
		msgConf.DisableNotification = true
		message, err := t.bot.Send(msgConf)
		if err != nil {
			logrus.Errorln(err)
			continue
		}

		logrus.Debugf("В чат ID=%d отправлено сообщение: \n%v ", msgConf.ChatID, message)
	}
}
