package processors

import (
	"context"

	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/worker/commands"

	tg "github.com/AlexSkilled/go_tg/pkg"
	tgmodel "github.com/AlexSkilled/go_tg/pkg/model"
	"github.com/sirupsen/logrus"
)

type Settings struct {
	interfaces.SettingsService
}

func NewSettingsProcessor(storage interfaces.ServiceStorage) *Settings {
	s := &Settings{storage.Settings()}

	return s
}

func (s *Settings) Handle(ctx context.Context, message *tgmodel.MessageIn) (out tg.TgMessage) {
	if len(message.Args) == 0 {
		return &tgmodel.Callback{Command: commands.Settings, Type: tgmodel.Callback_Type_TransitToMenu}
	}
	switch message.Args[0] {
	case commands.ChangeLanguage:
		if len(message.Args) > 1 {
			err := s.SettingsService.ChangeLanguage(message.From.ID, message.Args[1])
			if err != nil {
				logrus.Error(err)
			}
		}
	}
	return nil
}

func (s *Settings) Dump(_ int64) {}
