package service

import (
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"
	"time"
)

type TokenExpiration struct {
	gits []model.GitUser
}

var _ interfaces.TokenExpiration = (*TokenExpiration)(nil)

func (t *TokenExpiration) Start() {
	ticker := time.NewTicker(time.Hour * 24)

	for {
		<-ticker.C

	}
}

func (t *TokenExpiration) AddGitIds(gits []model.GitUser) {

}

func (t *TokenExpiration) Stop() {

}
