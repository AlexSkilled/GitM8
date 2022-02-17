package processors

import (
	"context"
	"fmt"

	"gitlab-tg-bot/internal/message-handling/help"
	helpGitlab "gitlab-tg-bot/internal/message-handling/help/gitlab"
	"gitlab-tg-bot/internal/message-handling/info"
	"gitlab-tg-bot/internal/message-handling/langs"
	"gitlab-tg-bot/worker/commands"

	tg "github.com/AlexSkilled/go_tg/pkg"
	tgmodel "github.com/AlexSkilled/go_tg/pkg/model"
)

type HelpProcessor struct{}

const (
	help_GetGitlabToken      = "how_to_get_gitlab_token"
	help_SetupGitlab_Webhook = "how_to_setup_gitlab_webhook"

	Help_GetGitlabToken      = commands.Help + " " + help_GetGitlabToken
	Help_SetupGitlab_Webhook = commands.Help + " " + help_SetupGitlab_Webhook
)

func (h *HelpProcessor) Handle(ctx context.Context, message *tgmodel.MessageIn) (out tg.TgMessage) {
	if len(message.Args) == 0 {
		return &tgmodel.MessageOut{Text: langs.Get(ctx, help.Defautlmessage)}
	}

	switch message.Args[0] {
	case help_GetGitlabToken:
		if len(message.Args) < 2 {
			return &tgmodel.MessageOut{
				Text: langs.Get(ctx, info.ErrorNotEnoughArguments) + " " + langs.Get(ctx, info.SpecifyGitSource),
			}
		}
		return &tgmodel.MessageOut{
			Text: fmt.Sprintf(langs.Get(ctx, helpGitlab.CreateToken), message.Args[0]),
		}
	case help_SetupGitlab_Webhook:
		return &tgmodel.MessageOut{
			Text: fmt.Sprintf(langs.Get(ctx, helpGitlab.SetupWebhookInstruction)),
		}
	}
	return &tgmodel.MessageOut{Text: langs.Get(ctx, help.Defautlmessage)}
}

func (h *HelpProcessor) Dump(_ int64) {}
