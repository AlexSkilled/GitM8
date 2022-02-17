package processors

import (
	"context"
	"strings"

	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/internal/message-handling/info"
	"gitlab-tg-bot/internal/message-handling/langs"
	"gitlab-tg-bot/internal/message-handling/managing"
	"gitlab-tg-bot/internal/message-handling/start"
	"gitlab-tg-bot/worker/commands"

	tg "github.com/AlexSkilled/go_tg/pkg"
	tgmodel "github.com/AlexSkilled/go_tg/pkg/model"
)

type ManageProcessor struct {
	subscription interfaces.SubscriptionService
}

func NewManageProcessor(services interfaces.ServiceStorage) *ManageProcessor {
	return &ManageProcessor{
		services.Subscription(),
	}
}

func (m *ManageProcessor) Handle(ctx context.Context, message *tgmodel.MessageIn) (out tg.TgMessage) {
	if len(message.Args) == 0 {
		return m.menu(ctx, message)
	}
	return nil
}

func (m *ManageProcessor) Dump(_ int64) {}

func (m *ManageProcessor) menu(ctx context.Context, message *tgmodel.MessageIn) *tgmodel.MessageOut {
	tickets, err := m.subscription.GetUserTickets(message.From.ID)
	if err != nil {
		return &tgmodel.MessageOut{
			Text: langs.Get(ctx, info.Error),
		}
	}

	if len(tickets) == 0 {
		registerTicket := &tgmodel.InlineKeyboard{}
		registerTicket.AddButton(langs.Get(ctx, start.Register), commands.Register)
		return &tgmodel.MessageOut{
			Text:          langs.Get(ctx, managing.NoTickets),
			InlineButtons: registerTicket,
		}
	}

	ticketsButtons := &tgmodel.InlineKeyboard{Columns: 3}

	for _, item := range tickets {
		ticketsButtons.AddButton(item.Name, strings.Join([]string{commands.Manage, item.Name}, ""))
	}
	return &tgmodel.MessageOut{
		Text:          langs.Get(ctx, managing.Tickets),
		InlineButtons: ticketsButtons,
	}
}
