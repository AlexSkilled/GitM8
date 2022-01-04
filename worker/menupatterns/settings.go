package menupatterns

import (
	"errors"
	"fmt"

	"gitlab-tg-bot/internal/message-handling/langs"
	"gitlab-tg-bot/internal/message-handling/menu/settingsmenu"

	tgmodel "github.com/AlexSkilled/go_tg/pkg/model"
)

func NewSettingsMenu(locale string) (tgmodel.MenuPattern, error) {
	languagesMenu, ok := menus[langs.GetWithLocale(locale, settingsmenu.Language)][locale]
	if !ok {
		return tgmodel.MenuPattern{}, errors.New(fmt.Sprintf("Для локали %s отсутствует меню языка", locale))
	}

	settingsMenu := tgmodel.NewMenuPattern(langs.GetWithLocale(locale, settingsmenu.Name))
	settingsMenu.AddMenuButton(langs.GetWithLocale(locale, settingsmenu.Language), languagesMenu.GetCallCommand())

	addMenu(locale, settingsMenu)

	return settingsMenu, nil
}
