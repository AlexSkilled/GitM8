package processors

import (
	"context"
	"fmt"
	"strings"

	"gitlab-tg-bot/internal/interfaces"
	helpGitlab "gitlab-tg-bot/internal/message-handling/help/gitlab"
	"gitlab-tg-bot/internal/message-handling/info"
	"gitlab-tg-bot/internal/message-handling/langs"
	"gitlab-tg-bot/internal/message-handling/register"
	"gitlab-tg-bot/internal/message-handling/start"
	"gitlab-tg-bot/service/model"
	"gitlab-tg-bot/utils"
	"gitlab-tg-bot/worker/commands"

	tg "github.com/AlexSkilled/go_tg/pkg"
	tgmodel "github.com/AlexSkilled/go_tg/pkg/model"
	"github.com/sirupsen/logrus"
)

const (
	StepChoseRegistrationType interfaces.StepName = iota
	StepDomain
	StepToken
)

const (
	RegisterGetWebhookURL = "webhook"
	RegisterToken         = "token"
)

type Register struct {
	services      interfaces.ServiceStorage
	dialogContext map[int64]*registrationProcess // [tgUserId]->regForm
}

type registrationProcess struct {
	Type        string
	GitlabName  string
	GitlabToken string
	Domain      string
	CurrentStep interfaces.StepName
}

func (r *registrationProcess) ToDto() model.GitUser {
	return model.GitUser{
		Username: r.GitlabName,
		Token:    r.GitlabToken,
		Domain:   r.Domain,
	}
}

func NewRegisterProcessor(services interfaces.ServiceStorage) *Register {
	return &Register{
		services:      services,
		dialogContext: map[int64]*registrationProcess{},
	}
}

func (r *Register) Handle(ctx context.Context, message *tgmodel.MessageIn) (out tg.TgMessage) {
	registration, ok := r.dialogContext[message.From.ID]
	if !ok {
		return r.initRegistration(ctx, message)
	}

	switch registration.CurrentStep {
	case StepDomain:
		registration.Domain = message.Text
		registration.CurrentStep = StepChoseRegistrationType

		locale, err := utils.ExtractLocale(ctx)
		if err != nil {
			locale = langs.GetDefaultLocale()
		}

		registerMenu := &tgmodel.InlineKeyboard{Columns: 2}

		registerMenu.AddButton(langs.GetWithLocale(locale, register.ButtonToken), RegisterToken)
		registerMenu.AddButton(langs.GetWithLocale(locale, register.ButtonUrl), RegisterGetWebhookURL)

		return &tgmodel.MessageOut{
			Text:          langs.Get(ctx, register.WebhookOrTokenMessage),
			InlineButtons: registerMenu,
		}
	case StepChoseRegistrationType:
		return r.selectRegistrationType(ctx, message, registration)
	case StepToken:
		registration.GitlabToken = message.Text
		registration.CurrentStep++
	}

	err := r.services.User().AddGitAccount(message.From.ID, registration.ToDto())
	var response string
	if err != nil {
		response = langs.Get(ctx, register.Error)
		if strings.Contains(err.Error(), "<401>") {
			response = fmt.Sprintf(response, langs.Get(ctx, register.ErrorInvalidToken))
		} else {
			response = fmt.Sprintf(response, fmt.Sprintf(langs.Get(ctx, info.Error), err.Error()))
			logrus.Errorln(err)
		}
	} else {
		response = langs.Get(ctx, register.Success)
	}

	delete(r.dialogContext, message.From.ID)
	return &tgmodel.Callback{
		Command: commands.Start,
		Text:    response,
	}
}

func (r *Register) initRegistration(ctx context.Context, message *tgmodel.MessageIn) (out tg.TgMessage) {
	r.dialogContext[message.From.ID] = &registrationProcess{
		CurrentStep: StepDomain,
	}

	gits := &tgmodel.InlineKeyboard{Columns: 1}
	gits.AddButton(string(model.Gitlab), string(model.Gitlab)+".com")
	gits.AddButton(langs.Get(ctx, start.MainMenu), commands.Start)
	return &tgmodel.MessageOut{
		Text:          langs.Get(ctx, register.AskDomain),
		InlineButtons: gits,
	}
}

func (r *Register) selectRegistrationType(ctx context.Context, message *tgmodel.MessageIn, registration *registrationProcess) (out tg.TgMessage) {
	switch message.Text {
	case RegisterToken:
		registration.CurrentStep = StepToken
		registration.Type = RegisterToken
		return &tg.MultipleMessage{
			message.Chat.ID: []tg.TgMessage{
				&tgmodel.MessageOut{
					Text: fmt.Sprintf(langs.Get(ctx, helpGitlab.CreateToken), registration.Domain),
				},
				&tgmodel.MessageOut{
					Text: langs.Get(ctx, register.EnterToken),
				},
			},
		}
	case RegisterGetWebhookURL:
		webhook, err := r.services.GitApi().GetWebhookUrl(registration.Domain, message.From.ID)
		if err != nil {
			return &tgmodel.MessageOut{
				Text: langs.Get(ctx, info.ErrorCouldNotFindDomain),
			}
		}
		webhookMenu := &tgmodel.InlineKeyboard{Columns: 1}

		webhookMenu.AddButton(langs.Get(ctx, register.ButtonSetupWebhook), Help_SetupGitlab_Webhook)
		webhookMenu.AddButton(langs.Get(ctx, start.MainMenu), commands.Start)

		delete(r.dialogContext, message.From.ID)

		return &tgmodel.MessageOut{
			Text:          webhook,
			InlineButtons: webhookMenu,
		}
	}

	return &tgmodel.MessageOut{
		Text: langs.Get(ctx, info.ErrorAnotherAnswerExpected),
	}
}

func (r *Register) Dump(userId int64) {
	delete(r.dialogContext, userId)
}
