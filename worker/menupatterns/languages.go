package menupatterns

import (
	"gitlab-tg-bot/internal/message-handling/info"
	"gitlab-tg-bot/internal/message-handling/langs"
	"gitlab-tg-bot/internal/message-handling/menu/settingsmenu"
	"gitlab-tg-bot/worker/commands"

	tgmodel "github.com/AlexSkilled/go_tg/pkg/model"
)

func NewLanguagesMenu(locale string) (tgmodel.MenuPattern, error) {
	languagesMenu := tgmodel.NewMenuPattern(langs.GetWithLocale(locale, settingsmenu.Language))

	for _, item := range langs.AvailableLangs {
		languagesMenu.AddMenuButton(langs.GetWithLocale(item, info.Name), commands.Settings+" "+commands.ChangeLanguage+" "+item)
	}

	addMenu(locale, languagesMenu)

	return languagesMenu, nil
}
