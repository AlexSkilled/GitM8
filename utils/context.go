package utils

import (
	"context"
	"errors"
	"gitlab-tg-bot/service/model"
)

type ContextKey int32

const (
	ContextKey_User ContextKey = 1
)

func ExtractUser(ctx context.Context) (model.User, error) {
	user, ok := ctx.Value(ContextKey_User).(model.User)
	if !ok {
		return model.User{}, errors.New("Не удалось получить пользователя.")
	}

	return user, nil
}
