package menupatterns

import (
	"errors"
	"fmt"

	"gitlab-tg-bot/internal/message-handling/langs"
	"gitlab-tg-bot/internal/message-handling/menu/settingsmenu"
	"gitlab-tg-bot/internal/message-handling/start"
	"gitlab-tg-bot/worker/commands"

	tgmodel "github.com/AlexSkilled/go_tg/pkg/model"
)

func NewSettingsMenu() (tgmodel.MenuPattern, error) {
	languagesMenu, ok := menus[commands.ChangeLanguage]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Отсутствует меню выбора языка"))
	}

	settingsMenu := tgmodel.NewLocalizedMenuPattern(commands.Settings)

	for _, locale := range langs.AvailableLangs {
		settingsMenu.AddMenu(locale, langs.GetWithLocale(locale, settingsmenu.Name))
		settingsMenu.AddMenuButton(locale,
			langs.GetWithLocale(locale, settingsmenu.Language),
			languagesMenu.GetTransitionCommand())
		settingsMenu.AddMenuButton(locale, langs.GetWithLocale(locale, start.MainMenu), commands.Start)
	}
	addMenu(commands.Settings, settingsMenu)

	return settingsMenu, nil
}
