package utils

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func NewTgMessageWithButtons(chatId int64, messageText string, buttons []tgbotapi.InlineKeyboardButton) tgbotapi.Chattable {
	markup := tgbotapi.NewInlineKeyboardMarkup(buttons)

	message := tgbotapi.NewMessage(chatId, messageText)
	message.BaseChat.ReplyMarkup = markup

	return message
}
