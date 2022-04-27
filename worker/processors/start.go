package processors

import (
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
		SubscriptionService: storage.Subscription(),
	}
}

func (s *Start) Handle(in *tgmodel.MessageIn, out tg.Sender) {
	locale, err := utils.ExtractLocale(in.Ctx)
	if err != nil {
		locale = langs.GetDefaultLocale()
	}
	buttons := map[string]string{}

	tickets, err := s.SubscriptionService.GetUserTickets(in.From.ID)
	if err != nil {
		in.Response(tgmodel.MessageOut{
			Text:   fmt.Sprintf("Internal server error %s", err),
			ChatId: in.Chat.ID,
		})
		return
	}

	if tickets != nil {
		buttons[start.Management] = commands.Manage
	}

	menu, err := menupatterns.NewMainMenu(locale, buttons)
	if err != nil {
		out.Send(&tgmodel.MessageOut{
			Text: fmt.Sprintf("Internal server error %s", err),
		})
		return
	}

	out.Send(&tgmodel.Callback{
		Text: langs.GetWithLocale(locale, start.MainMenu),
		Type: tgmodel.Callback_Type_OpenMenu,
		Menu: menu,
	})
}

func (s *Start) Dump(_ int64) {}
