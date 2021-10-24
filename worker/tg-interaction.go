package worker

import (
	"context"
	"log"

	config "gitlab-tg-bot/conf"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"
	"gitlab-tg-bot/utils"
	"gitlab-tg-bot/worker/commands"
	"gitlab-tg-bot/worker/processors"

	"github.com/go-pg/pg/v9"

	tg "github.com/AlexSkilled/go_tg/pkg"
	tgmodel "github.com/AlexSkilled/go_tg/pkg/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Worker struct {
	interfaces.ServiceStorage

	bt           *tg.Bot
	conf         interfaces.Configuration
	interceptors map[int64]interfaces.Interceptor
}

func NewTelegramWorker(conf interfaces.Configuration, services interfaces.ServiceStorage) interfaces.TelegramWorker {
	bot := tg.NewBot(conf.GetString(config.Token))

	log.Printf("Авторизация в боте %s", bot.Bot.Self.UserName)

	bot.AddCommandHandler(processors.NewStartProcessor(services), commands.Start)
	bot.AddCommandHandler(processors.NewRegisterProcessor(services), commands.Register)

	bot.AddCommandHandler(processors.NewSubscribeProcessor(services), commands.Subscribe)

	return &Worker{
		bt:             bot,
		ServiceStorage: services,
		conf:           conf,
		interceptors:   map[int64]interfaces.Interceptor{},
	}
}

func (t *Worker) Start() {
	t.bt.EnrichContext = t
	t.bt.Start()
}

// SendMessages TODO
func (t *Worker) SendMessages(messages []model.OutputMessage) {
	for _, msg := range messages {
		msgConf := tgbotapi.NewMessage(msg.ChatId, msg.Msg)
		msgConf.ParseMode = "Markdown"
		msgConf.DisableWebPagePreview = true

		msgConf.DisableNotification = msg.DisableNotification

		message, err := t.bt.Bot.Send(msgConf)
		if err != nil {
			logrus.Errorln(err)
			continue
		}

		logrus.Debugf("В чат ID=%d отправлено сообщение: \n%v ", msgConf.ChatID, message)
	}
}

func (t *Worker) GetContext(message *tgmodel.MessageIn) (context.Context, error) {
	user, err := t.User().GetWithGitlabUsersById(message.From.ID)
	if err != nil {
		if err == pg.ErrNoRows {
			user, err = t.User().Register(model.User{
				Id:         message.From.ID,
				Name:       message.From.FirstName,
				TgUsername: message.From.UserName,
			})
			if err != nil {
				logrus.Errorln(err)
				return nil, err
			}
			//if err == nil {
			//t.processors[CommandStart].Handle(context.Background(), message, t.bot)
			//}
		}
	}

	ctx := context.WithValue(context.Background(), utils.ContextKey_User, user)
	ctx = context.WithValue(ctx, utils.ContextKey_ChatId, message.Chat.ID)
	return ctx, nil
}
