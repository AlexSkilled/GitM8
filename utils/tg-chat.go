package utils

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewTgMessageButtonsMarkup(buttons []tgbotapi.InlineKeyboardButton, buttonsInRaw int) tgbotapi.InlineKeyboardMarkup {
	finalButtonsSet := make([][]tgbotapi.InlineKeyboardButton, 0, 1)
	i := 0
	raw := -1
	for i < len(buttons) {
		if i%buttonsInRaw == 0 {
			finalButtonsSet = append(finalButtonsSet, make([]tgbotapi.InlineKeyboardButton, 0, buttonsInRaw))
			raw++
		}
		finalButtonsSet[raw] = append(finalButtonsSet[raw], buttons[i])
		i++
	}

	return tgbotapi.InlineKeyboardMarkup{InlineKeyboard: finalButtonsSet}
}

func AppendWithPattern(sb *strings.Builder, pattern string, replaces ...string) *strings.Builder {
	if strings.Count(pattern, "%s") != len(replaces) {
		logrus.Error("Для шаблона " + pattern + " не хватает данных, для подстановки")
	}
	sb.WriteString("\n" + fmt.Sprintf(pattern, replaces))
	return sb
}

func EscapeLinkSymbols(in string) (out string) {
	in = strings.Replace(in, "_", "\\_", -1)
	return in
}

func EscapeNameSymbols(in string) (out string) {
	in = EscapeLinkSymbols(in)

	in = strings.Replace(in, ".", "\\.", -1)
	in = strings.Replace(in, "*", "\\*", -1)
	in = strings.Replace(in, "[", "\\[", -1)
	in = strings.Replace(in, "]", "\\]", -1)
	in = strings.Replace(in, "`", "\\`", -1)
	in = strings.Replace(in, "(", "\\(", -1)
	in = strings.Replace(in, ")", "\\)", -1)
	in = strings.Replace(in, ">", "\\>", -1)
	return in
}
