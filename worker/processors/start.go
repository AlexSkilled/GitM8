package processors

import (
	"context"
	"fmt"

	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/internal/message-handling/langs"
	"gitlab-tg-bot/internal/message-handling/start"
	"gitlab-tg-bot/utils"
	"gitlab-tg-bot/worker/commands"
	"gitlab-tg-bot/worker/menupatterns"

	tg "github.com/AlexSkilled/go_tg/pkg"
	tgmodel "github.com/AlexSkilled/go_tg/pkg/model"
)

type Start struct {
	interfaces.SubscriptionService
}

func NewStartProcessor(storage interfaces.ServiceStorage) *Start {
	return &Start{
		storage.Subscription(),
	}
}

func (s *Start) Handle(ctx context.Context, message *tgmodel.MessageIn) (out tg.TgMessage) {
	locale, err := utils.ExtractLocale(ctx)
	if err != nil {
		locale = langs.GetDefaultLocale()
	}
	buttons := map[string]string{}

	tickets, err := s.SubscriptionService.GetUserTickets(message.From.ID)
	if err != nil {
		return &tgmodel.MessageOut{
			Text: fmt.Sprintf("Internal server error %s", err),
		}
	}

	if tickets != nil {
		buttons[start.Management] = commands.Manage
	}

	menu, err := menupatterns.NewMainMenu(locale, buttons)
	if err != nil {
		return &tgmodel.MessageOut{
			Text: fmt.Sprintf("Internal server error %s", err),
		}
	}

	return &tgmodel.MessageOut{
		Text:          langs.GetWithLocale(locale, start.MainMenu),
		InlineButtons: menu,
	}
}

func (s *Start) Dump(_ int64) {}
