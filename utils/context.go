package utils

import (
	"context"
	"errors"

	"gitlab-tg-bot/service/model"
)

type ContextKey int32

const (
	ContextKey_User   ContextKey = 1
	ContextKey_ChatId ContextKey = 2
	ContextKey_Locale ContextKey = 3
)

func ExtractUser(ctx context.Context) (model.User, error) {
	user, ok := ctx.Value(ContextKey_User).(model.User)
	if !ok {
		return model.User{}, errors.New("Не удалось получить пользователя.")
	}

	return user, nil
}

func ExtractChatId(ctx context.Context) (chatId int64, err error) {
	user, ok := ctx.Value(ContextKey_ChatId).(int64)
	if !ok {
		return 0, errors.New("Не удалось получить id чата.")
	}

	return user, nil
}

func ExtractLocale(ctx context.Context) (string, error) {
	locale, ok := ctx.Value(ContextKey_Locale).(string)
	if !ok {
		return "", errors.New("Не удалось получить пользователя.")
	}

	return locale, nil
}
