package processors

import (
	"context"
	"fmt"

	"gitlab-tg-bot/internal/message-handling/help"
	helpGitlab "gitlab-tg-bot/internal/message-handling/help/gitlab"
	"gitlab-tg-bot/internal/message-handling/info"
	"gitlab-tg-bot/internal/message-handling/langs"

	tg "github.com/AlexSkilled/go_tg/pkg"
	tgmodel "github.com/AlexSkilled/go_tg/pkg/model"
)

type HelpProcessor struct {
}

const (
	Help_GetGitlabToken = "how_to_get_gitlab_token"
)

func (h *HelpProcessor) Handle(ctx context.Context, message *tgmodel.MessageIn) (out tg.TgMessage) {
	if len(message.Args) == 0 {
		return &tgmodel.MessageOut{Text: langs.Get(ctx, help.Defautlmessage)}
	}

	switch message.Args[0] {
	case Help_GetGitlabToken:
		if len(message.Args) < 2 {
			return &tgmodel.MessageOut{
				Text: langs.Get(ctx, info.ErrorNotEnoughArguments) + " " + langs.Get(ctx, info.SpecifyGitSource),
			}
		}
		return &tgmodel.MessageOut{
			Text: fmt.Sprintf(langs.Get(ctx, helpGitlab.CreateToken), message.Args[0]),
		}
	}
	return &tgmodel.MessageOut{Text: langs.Get(ctx, help.Defautlmessage)}
}

func (h *HelpProcessor) Dump(_ int64) {}
