package interfaces

import (
	"gitlab-tg-bot/service/model"
)

type ProviderStorage interface {
	User() UserProvider
	Ticket() TicketProvider
}

type UserProvider interface {
	Create(user model.User) (model.User, error)
	Get(id int64) (model.User, error)
	GetWithGitlabUsers(id int64) (model.User, error)
	AddGitlab(userId int64, gitlab model.GitlabUser) error
}

type SubscriptionProvider interface {
	GetSubscription() (model.Subscription, error)
	Subscribe([]model.Subscription) error
	Unsubscribe([]model.Subscription) error
}

type TicketProvider interface {
	AddTicket(model.Ticket) (tickerId int32, err error)
	AddTicketToUser(userId int64, ticketId int32) error
}
