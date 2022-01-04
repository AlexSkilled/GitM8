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
var menus = make(map[string]tgmodel.MenuPattern)

func addMenu(command string, menu tgmodel.MenuPattern) {
	menus[command] = menu
}

func NewMainMenu() (tgmodel.MenuPattern, error) {

	mainMenu := tgmodel.NewLocalizedMenuPattern(commands.Start)

	for _, locale := range langs.AvailableLangs {
		mainMenuName := langs.GetWithLocale(locale, mainmenu.Name)
		mainMenu.AddMenu(locale, mainMenuName)

		// TODO остальное запихать
		settingsName := langs.GetWithLocale(locale, settingsmenu.Name)

		settings, ok := menus[commands.Settings]
		if !ok {
			return nil, errors.New(fmt.Sprintf("Ошибка при создании главного меню. Для локали %s отсутствует меню настроек ", locale))
		}
		mainMenu.AddMenuButton(locale, settingsName, settings.GetTransitionCommand())
	}

	addMenu(commands.Start, mainMenu)

	return mainMenu, nil
}
