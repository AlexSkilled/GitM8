package interfaces

import (
	"gitlab-tg-bot/service/model"
)

type ProviderStorage interface {
	User() UserProvider
	Ticket() TicketProvider
	MessagePattern() MessagePatternProvider
}

type UserProvider interface {
	Create(user model.User) (model.User, error)
	Get(id int64) (model.User, error)
	GetWithGitlabUsers(id int64) (model.User, error)
	AddGit(userId int64, gitlab model.GitUser) error
	Update(user model.User) error
	GetURN(userId int64) (string, error)
	SetURN(userId int64, urn string) error
}

type SubscriptionProvider interface {
	GetSubscription() (model.Subscription, error)
	Subscribe([]model.Subscription) error
	Unsubscribe([]model.Subscription) error
}

type TicketProvider interface {
	AddTicket(model.Ticket) (tickerId int32, err error)
	AddTicketToChat(chatId int64, ticketId int32) error
	GetTicketsToSend(repoId string, hookType model.GitHookType) ([]model.TicketChatId, error)
	GetGitUsersByTicketId(ticketId int32) ([]model.GitUser, error)
	GetOwnerByTicketId(ticketId int32) (model.GitUser, error)
}

type MessagePatternProvider interface {
	GetMessage(lang string, hookType model.GitHookType, subType model.GitHookSubtype) (string, map[string]string, error)
}
