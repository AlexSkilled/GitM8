package service

import (
	"gitlab-tg-bot/internal/interfaces"
)

type Announcer struct {
	data interfaces.ProviderStorage
}

var _ interfaces.AnnouncerService = (*Announcer)(nil)

func NewAnnouncer(providerStorage interfaces.ProviderStorage) *Announcer {
	return &Announcer{
		data: providerStorage,
	}
}

func (a *Announcer) Announce(message interface{}) {
	panic("implement me")
}
