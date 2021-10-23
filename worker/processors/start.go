package processors

import (
	"context"

	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/worker/commands"

	tg "github.com/AlexSkilled/go_tg/pkg"
	tgmodel "github.com/AlexSkilled/go_tg/pkg/model"
)

type Start struct {
	services interfaces.ServiceStorage
}

var _ tg.CommandHandler = (*Start)(nil)

func NewStartProcessor(services interfaces.ServiceStorage) *Start {
	return &Start{
		services: services,
	}
}

func (s *Start) IsInterceptor() bool {
	return false
}

func (s *Start) Handle(_ context.Context, _ *tgmodel.MessageIn) (out *tgmodel.MessageOut) {
	btns := &tgmodel.InlineKeyboard{}
	btns.AddButton("Регистрация", commands.Register)

	return &tgmodel.MessageOut{
		Text:          "Главное меню",
		InlineButtons: btns,
	}

}
