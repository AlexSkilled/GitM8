package provider

import (
	"github.com/go-pg/pg/v9"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/internal/model"
)

type Subscribe struct {
	Email string
}

type SubscriptionProvider struct {
	*pg.DB
}

var _ interfaces.SubscriptionProvider = (*SubscriptionProvider)(nil)

func NewSubscriptionProvider(db *pg.DB) interfaces.SubscriptionProvider {
	return &SubscriptionProvider{
		DB: db,
	}
}

func (s *SubscriptionProvider) GetSubscription() (model.Subscription, error) {
	panic("implement me")
}

func (s *SubscriptionProvider) Subscribe(subscriptions []model.Subscription) error {
	panic("implement me")
}

func (s *SubscriptionProvider) Unsubscribe(subscriptions []model.Subscription) error {
	panic("implement me")
}
