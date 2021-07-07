package interfaces

import "gitlab-tg-bot/internal/model"

type ProviderStorage interface {
	User() UserProvider
}

type UserProvider interface {
	Create(user model.User) error
	Get(id int32) (model.User, error)
	GetByChatId(id int32) (model.User, error)
	GetWithGitlabUsers(id int32) (model.User, error)
	AddGitlab(userId int32, gitlab model.GitlabUser) error
}

type SubscriptionProvider interface {
	GetSubscription() (model.Subscription, error)
	Subscribe([]model.Subscription) error
	Unsubscribe([]model.Subscription) error
}
