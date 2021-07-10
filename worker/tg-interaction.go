package worker

import (
	"context"
	"gitlab-tg-bot/internal"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/internal/model"
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

func (t *Worker) handleCommands(ctx context.Context, userId int64, update tgbotapi.Update) {
	if update.Message.Text == CommandExit {
		// МБ в будущем будет необходимость прерывать не только заполнение форм, так что да
		if interceptor, ok := t.interceptors[userId]; ok {
			interceptor.DumpChatSession(update.Message.Chat.ID)
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
		var userId int64
		var ctx context.Context
		var err error
		if update.Message != nil {
			userId = update.Message.From.ID

			logrus.Infof("Пользователь %d в чате %d написал %s", userId, update.Message.Chat.ID, update.Message.Text)

			ctx, err = t.fillContext(userId, update)
			if err != nil {
				logrus.Errorln(err)
				continue
			}

			if strings.HasPrefix(update.Message.Text, CommandPrefix) {
				t.handleCommands(ctx, userId, update)
				continue
			}

		} else if update.CallbackQuery == nil {

			userId = update.CallbackQuery.From.ID

			logrus.Infof("Пользователь %d в чате %d ответил %s со значением %s", userId, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.Text, update.CallbackQuery.Data)

			ctx, err = t.fillContext(userId, update)
			if err != nil {
				logrus.Errorln(err)
				continue
			}
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

func (t *Worker) fillContext(userId int64, update tgbotapi.Update) (context.Context, error) {

	user, err := t.User().GetWithGitlabUsersByTgId(userId)
	if err != nil {
		if err == pg.ErrNoRows {
			user, err = t.User().Register(model.User{
				Id:         userId,
				Name:       update.Message.From.FirstName,
				TgUsername: update.Message.From.UserName,
			})

			if err == nil {
				t.processors[CommandStart].Process(context.Background(), update, t.bot)
			}
		}

		if err != nil {
			// TODO логирование
			//continue
		}
	}

	ctx := context.WithValue(context.Background(), utils.ContextKey_User, user)

	return ctx, err
}
