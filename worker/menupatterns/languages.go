package menupatterns

import (
	"gitlab-tg-bot/internal/message-handling/info"
	"gitlab-tg-bot/internal/message-handling/langs"
	"gitlab-tg-bot/internal/message-handling/menu/settingsmenu/languagemenu"
	"gitlab-tg-bot/worker/commands"

	tgmodel "github.com/AlexSkilled/go_tg/pkg/model"
)

func NewLanguagesMenu() (tgmodel.MenuPattern, error) {

	languagesMenu := tgmodel.NewLocalizedMenuPattern(commands.ChangeLanguage)
	for _, locale := range langs.AvailableLangs {
		languagesMenu.AddMenu(locale, langs.GetWithLocale(locale, languagemenu.LanguageMenuName))
		for _, language := range langs.AvailableLangs {
			languagesMenu.AddMenuButton(locale, langs.GetWithLocale(language, info.Name), commands.Settings+" "+commands.ChangeLanguage+" "+language)
		}

	}

	addMenu(commands.ChangeLanguage, languagesMenu)

	return languagesMenu, nil
}
