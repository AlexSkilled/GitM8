package service

import "gitlab-tg-bot/internal/interfaces"

type Announcer struct {
}

var _ interfaces.AnnouncerService = (*Announcer)(nil)

func NewAnnouncer() *Announcer {
	return &Announcer{}
}

func (a *Announcer) Announce(message interface{}) {
	panic("implement me")
}
