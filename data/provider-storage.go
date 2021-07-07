package data

import (
	"github.com/go-pg/pg/v9"
	"gitlab-tg-bot/data/provider"
	"gitlab-tg-bot/internal/interfaces"
)

type ProviderStorage struct {
	interfaces.UserProvider
}

var _ interfaces.ProviderStorage = (*ProviderStorage)(nil)

func NewProviderStorage(db *pg.DB) ProviderStorage {
	return ProviderStorage{
		UserProvider: provider.NewUser(db),
	}
}

func (p *ProviderStorage) User() interfaces.UserProvider {
	return p.UserProvider
}
