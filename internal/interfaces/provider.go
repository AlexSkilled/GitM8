package interfaces

import "gitlab-tg-bot/internal/model"

type UserProvider interface {
	Create(user model.User) error
	Get(id int32) (model.User, error)
	GetByTelegramId(id int64) (model.User, error)
}

type SubscriptionProvider interface {
	GetSubscription() (model.Subscription, error)
	Subscribe([]model.Subscription) error
	Unsubscribe([]model.Subscription) error
}
