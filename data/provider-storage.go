package data

import (
	"gitlab-tg-bot/data/provider"
	"gitlab-tg-bot/internal/interfaces"

	"github.com/go-pg/pg/v9"
)

type ProviderStorage struct {
	interfaces.UserProvider
	interfaces.TicketProvider
	interfaces.MessagePatternProvider
}

var _ interfaces.ProviderStorage = (*ProviderStorage)(nil)

func NewProviderStorage(db *pg.DB) *ProviderStorage {
	return &ProviderStorage{
		UserProvider:           provider.NewUser(db),
		TicketProvider:         provider.NewTicket(db),
		MessagePatternProvider: provider.NewMessagePatternProvider(db),
	}
}

func (p *ProviderStorage) User() interfaces.UserProvider {
	return p.UserProvider
}

func (p *ProviderStorage) Ticket() interfaces.TicketProvider {
	return p.TicketProvider
}

func (p *ProviderStorage) MessagePattern() interfaces.MessagePatternProvider {
	return p.MessagePatternProvider
}
