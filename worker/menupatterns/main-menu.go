package menupatterns

import (
	"errors"
	"fmt"

	"gitlab-tg-bot/internal/message-handling/langs"
	"gitlab-tg-bot/internal/message-handling/menu/mainmenu"
	"gitlab-tg-bot/internal/message-handling/menu/settingsmenu"
	"gitlab-tg-bot/worker/commands"

	tgmodel "github.com/AlexSkilled/go_tg/pkg/model"
)

// Menu type (setting) to locale (ru_RU) to actual pattern
var menus = make(map[string]map[string]tgmodel.MenuPattern)

func addMenu(locale string, pattern tgmodel.MenuPattern) {
	if _, ok := menus[pattern.Name]; !ok {
		menus[pattern.Name] = make(map[string]tgmodel.MenuPattern)
	}
	menus[pattern.Name][locale] = pattern
}

func NewMainMenu(locale string) (tgmodel.MenuPattern, error) {
	mainMenu := tgmodel.NewMenuPattern(langs.GetWithLocale(locale, mainmenu.Name))

	// TODO остальное запихать
	settingsName := langs.GetWithLocale(locale, settingsmenu.Name)

	settings, ok := menus[settingsName][locale]
	if !ok {
		return tgmodel.MenuPattern{}, errors.New(fmt.Sprintf("Ошибка при создании главного меню. Для локали %s отсутствует меню настроек ", locale))
	}
	mainMenu.AddMenuButton(settingsName, settings.GetCallCommand())
	mainMenu.AddEntryPoint(commands.Start)

	addMenu(locale, mainMenu)

	return mainMenu, nil
}
