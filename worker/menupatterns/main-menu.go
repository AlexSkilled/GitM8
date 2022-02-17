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

// Menu type (setting) to locale (ru_RU) to actual pattern
var menus = make(map[string]tgmodel.MenuPattern)

func addMenu(command string, menu tgmodel.MenuPattern) {
	menus[command] = menu
}

// NewMainMenu :
// arg buttons - is name path (e.g. start.MainMenu-"start:mainMenu") to command
func NewMainMenu(locale string, buttons map[string]string) (*tgmodel.InlineKeyboard, error) {

	mainMenu := tgmodel.InlineKeyboard{Columns: 2}

	for k, v := range buttons {
		mainMenu.AddButton(langs.GetWithLocale(locale, k), v)
	}

	registerButton := langs.GetWithLocale(locale, start.Register)
	mainMenu.AddButton(registerButton, commands.Register)

	settings, ok := menus[commands.Settings]
	if !ok {
		return nil, errors.New(
			fmt.Sprintf("Ошибка при создании главного меню. Для локали %s отсутствует меню настроек ", locale))
	}
	mainMenu.AddStandAloneButton(langs.GetWithLocale(locale, settingsmenu.Name), settings.GetTransitionCommand())

	return &mainMenu, nil
}
