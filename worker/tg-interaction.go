package worker

import (
	"context"
	"gitlab-tg-bot/internal"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/utils"
	processors "gitlab-tg-bot/worker/processors"
	"log"
	"strings"

	"github.com/go-pg/pg/v9"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

const (
	CommandPrefix string = "/"

	CommandRegister  string = "/register"
	CommandStart     string = "/start"
	CommandSubscribe string = "/subscribe"

	CommandExit string = "/exit"
)

type Worker struct {
	interfaces.ServiceStorage

	bot          *tgbotapi.BotAPI
	processors   map[string]interfaces.TgProcessor
	conf         internal.Configuration
	interceptors map[int64]interfaces.Interceptor
}

func NewTelegramWorker(conf internal.Configuration,
	services interfaces.ServiceStorage) interfaces.TelegramWorker {
	bot, err := tgbotapi.NewBotAPI(conf.GetString(internal.Token))
	if err != nil {
		panic(err)
	}

	bot.Debug = conf.GetBool(internal.Debug)

	log.Printf("Авторизация в боте %s", bot.Self.UserName)
	processorsMap := map[string]interfaces.TgProcessor{
		CommandStart:    processors.NewStartProcessor(services),
		CommandRegister: processors.NewRegisterProcessor(services),
	}

	return &Worker{
		bot:            bot,
		ServiceStorage: services,
		processors:     processorsMap,
		conf:           conf,
		interceptors:   map[int64]interfaces.Interceptor{},
	}
}

func (t *Worker) handleCommands(ctx context.Context, userId int64, update tgbotapi.Update) {
	if update.Message.Text == CommandExit {
		// МБ в будущем будет необходимость прерывать не только заполнение форм, так что да
		if interceptor, ok := t.interceptors[userId]; ok {
			interceptor.DumpUserSession(userId)
			delete(t.interceptors, userId)
		}

		_, _ = t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Выполнение прервано"))
		return
	}

	if processor, ok := t.processors[update.Message.Text]; ok {
		if processor.IsInterceptor() {
			t.interceptors[userId] = processor.(interfaces.Interceptor)
		}

		processor.Process(ctx, update, t.bot)
		return
	}

	logrus.Printf("Не знаю как обработать команду - `%s`", update.Message.Text)
}

func (t *Worker) Start() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updChan := t.bot.GetUpdatesChan(updateConfig)

	for update := range updChan {
		if update.Message == nil {
			continue
		}

		userId := update.Message.From.ID
		user, err := t.User().GetByTelegramId(userId)
		if err != nil {
			if err == pg.ErrNoRows {
				t.handleCommands(context.Background(), userId, update)
			}
			continue
		}

		ctx := context.Background()
		context.WithValue(ctx, utils.ContextKey_User, user)

		if strings.HasPrefix(update.Message.Text, CommandPrefix) {
			t.handleCommands(ctx, userId, update)
			continue
		}

		interceptor, ok := t.interceptors[userId]
		if ok {
			if interceptor.Process(ctx, update, t.bot) {
				delete(t.interceptors, userId)
			}
		}
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
