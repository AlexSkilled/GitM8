package utils

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func NewTgMessageWithButtons(chatId int64, messageText string, buttons []tgbotapi.InlineKeyboardButton, buttonsInRaw int) tgbotapi.Chattable {
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

	markup := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: finalButtonsSet}

	message := tgbotapi.NewMessage(chatId, messageText)
	message.BaseChat.ReplyMarkup = markup

	return message
}
