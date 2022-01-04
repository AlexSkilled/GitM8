package worker

import (
	"context"
	"log"

	config "gitlab-tg-bot/conf"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/internal/message-handling/langs"
	"gitlab-tg-bot/internal/message-handling/mainm"
	"gitlab-tg-bot/service/model"
	"gitlab-tg-bot/utils"
	"gitlab-tg-bot/worker/commands"
	"gitlab-tg-bot/worker/menupatterns"
	"gitlab-tg-bot/worker/processors"

	tg "github.com/AlexSkilled/go_tg/pkg"
	tgmodel "github.com/AlexSkilled/go_tg/pkg/model"
	"github.com/go-pg/pg/v9"
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

	bot.AddCommandHandler(processors.NewRegisterProcessor(services), commands.Register)

	bot.AddCommandHandler(processors.NewSubscribeProcessor(services), commands.Subscribe)

	bot.AddCommandHandler(processors.NewSettingsProcessor(services), commands.Settings)

	return &Worker{
		bt:             bot,
		ServiceStorage: services,
		conf:           conf,
		interceptors:   map[int64]interfaces.Interceptor{},
	}
}

func (t *Worker) addMenus(bot *tg.Bot) {
	menuPattern, err := menupatterns.NewLanguagesMenu()
	if err != nil {
		logrus.Error(err)
	} else {
		bot.AddMenu(menuPattern)
	}
	menuPattern, err = menupatterns.NewSettingsMenu()
	if err != nil {
		logrus.Error(err)
	} else {
		bot.AddMenu(menuPattern)
	}
	menuPattern, err = menupatterns.NewMainMenu()
	if err != nil {
		logrus.Error(err)
	} else {
		bot.AddMenu(menuPattern)
	}
}

func (t *Worker) Start() {
	t.bt.EnrichContext = t
	t.addMenus(t.bt)
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

		ctx := context.WithValue(context.Background(), utils.ContextKey_Locale, msg.Lang)

		logrus.Debugf(langs.Get(ctx, mainm.MessageSend), msgConf.ChatID, message)
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
		}
	}

	if len(user.Locale) == 0 {
		user.Locale = langs.GetDefaultLocale()
	}

	ctx := context.WithValue(context.Background(), utils.ContextKey_User, user)
	ctx = context.WithValue(ctx, utils.ContextKey_ChatId, message.Chat.ID)
	ctx = context.WithValue(ctx, utils.ContextKey_Locale, user.Locale)

	ctx = context.WithValue(ctx, tgmodel.LocaleContextKey, user.Locale)
	return ctx, nil
}
