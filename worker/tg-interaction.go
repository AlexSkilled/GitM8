package worker

import (
	"context"
	config "gitlab-tg-bot/conf"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"
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

	CommandStart     string = "/start"
	CommandRegister  string = "/register"
	CommandSubscribe string = "/subscribe"

	CommandExit string = "/exit"
)

type Worker struct {
	interfaces.ServiceStorage

	bot          *tgbotapi.BotAPI
	processors   map[string]interfaces.TgProcessor
	conf         interfaces.Configuration
	interceptors map[int64]interfaces.Interceptor
}

func NewTelegramWorker(conf interfaces.Configuration,
	services interfaces.ServiceStorage) interfaces.TelegramWorker {
	bot, err := tgbotapi.NewBotAPI(conf.GetString(config.Token))
	if err != nil {
		panic(err)
	}

	bot.Debug = conf.GetBool(config.Debug)

	log.Printf("Авторизация в боте %s", bot.Self.UserName)
	processorsMap := map[string]interfaces.TgProcessor{
		CommandStart:     processors.NewStartProcessor(services),
		CommandRegister:  processors.NewRegisterProcessor(services),
		CommandSubscribe: processors.NewSubscribeProcessor(services),
	}

	return &Worker{
		bot:            bot,
		ServiceStorage: services,
		processors:     processorsMap,
		conf:           conf,
		interceptors:   map[int64]interfaces.Interceptor{},
	}
}

func (t *Worker) handleCommands(ctx context.Context, message *tgbotapi.Message) {
	if message.Text == CommandExit {
		// МБ в будущем будет необходимость прерывать не только заполнение форм, так что да
		if interceptor, ok := t.interceptors[message.Chat.ID]; ok {
			interceptor.DumpChatSession(message.Chat.ID)
			delete(t.interceptors, message.Chat.ID)
		}

		_, _ = t.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Выполнение прервано"))
		return
	}

	if processor, ok := t.processors[message.Text]; ok {
		if processor.IsInterceptor() {
			if _, ok = t.interceptors[message.Chat.ID]; ok {
				t.interceptors[message.Chat.ID].DumpChatSession(message.Chat.ID)
			}
			t.interceptors[message.Chat.ID] = processor.(interfaces.Interceptor)
		}

		processor.Process(ctx, message, t.bot)
		return
	}

	logrus.Printf("Не знаю как обработать команду - `%s`", message.Text)
}

func (t *Worker) Start() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updChan := t.bot.GetUpdatesChan(updateConfig)

	for update := range updChan {

		var message *tgbotapi.Message

		if update.EditedMessage != nil {
			continue
		}

		if update.Message != nil {
			message = update.Message
		} else if update.CallbackQuery.Message != nil {
			message = update.CallbackQuery.Message
			message.Text = update.CallbackQuery.Data
			message.From = update.CallbackQuery.From
		} else {
			continue
		}

		var ctx context.Context
		var err error

		userId := message.From.ID

		chatId := message.Chat.ID
		logrus.Infof("Пользователь %d в чате %d написал %s", userId, chatId, message.Text)

		ctx, err = t.fillContext(message)
		if err != nil {
			logrus.Errorln(err)
			continue
		}

		if strings.HasPrefix(message.Text, CommandPrefix) {
			t.handleCommands(ctx, message)
			continue
		}

		interceptor, ok := t.interceptors[chatId]
		if ok {
			if interceptor.Process(ctx, message, t.bot) {
				delete(t.interceptors, userId)
			}
		}
	}
}

func (t *Worker) SendMessages(messages []model.OutputMessage) {
	for _, msg := range messages {
		msgConf := tgbotapi.NewMessage(msg.ChatId, msg.Msg)
		msgConf.ParseMode = "Markdown"
		msgConf.DisableWebPagePreview = true

		msgConf.DisableNotification = msg.DisableNotification

		message, err := t.bot.Send(msgConf)
		if err != nil {
			logrus.Errorln(err)
			continue
		}

		logrus.Debugf("В чат ID=%d отправлено сообщение: \n%v ", msgConf.ChatID, message)
	}
}

func (t *Worker) fillContext(message *tgbotapi.Message) (context.Context, error) {

	user, err := t.User().GetWithGitlabUsersById(message.From.ID)
	if err != nil {
		if err == pg.ErrNoRows {
			user, err = t.User().Register(model.User{
				Id:         message.From.ID,
				Name:       message.From.FirstName,
				TgUsername: message.From.UserName,
			})

			if err == nil {
				t.processors[CommandStart].Process(context.Background(), message, t.bot)
			}
		}

		if err != nil {
			logrus.Errorln(err)
			return nil, err
		}
	}

	ctx := context.WithValue(context.Background(), utils.ContextKey_User, user)
	ctx = context.WithValue(ctx, utils.ContextKey_ChatId, message.Chat.ID)
	return ctx, err
}
