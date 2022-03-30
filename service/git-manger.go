package service

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"unsafe"

	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"
)

type GitManager struct {
	sourceToService map[model.GitSource]interfaces.GitApiService
	user            interfaces.UserProvider
}

func NewGitManger(cont interfaces.Configuration, storage interfaces.ProviderStorage) *GitManager {
	return &GitManager{
		sourceToService: map[model.GitSource]interfaces.GitApiService{
			model.Gitlab: NewGitlabApiService(cont),
		},
		user: storage.User(),
	}
}

var _ interfaces.GitApiService = (*GitManager)(nil)

func (g *GitManager) GetRepositories(user model.GitUser) ([]model.Repository, error) {
	if api, ok := g.sourceToService[user.GitSource]; ok {
		return api.GetRepositories(user)
	} else {
		return nil, errors.New(fmt.Sprintf("Для git %s не нашлось обработчика", user.GitSource))
	}
}

func (g *GitManager) GetGroups(user model.GitUser) ([]model.Group, error) {
	if api, ok := g.sourceToService[user.GitSource]; ok {
		return api.GetGroups(user)
	} else {
		return nil, errors.New(fmt.Sprintf("Для git %s не нашлось обработчика", user.GitSource))
	}
}

func (g *GitManager) AddWebHook(user model.GitUser, hookInfo model.Hook) error {
	if api, ok := g.sourceToService[user.GitSource]; ok {
		return api.AddWebHook(user, hookInfo)
	} else {
		return errors.New(fmt.Sprintf("Для git %s не нашлось обработчика", user.GitSource))
	}
}

func (g *GitManager) GetUser(git model.GitUser, userId string) (model.GitUserInfo, error) {
	if api, ok := g.sourceToService[git.GitSource]; ok {
		return api.GetUser(git, userId)
	} else {
		return model.GitUserInfo{}, errors.New(fmt.Sprintf("Для git %s не нашлось обработчика", git.GitSource))
	}
}

func (g *GitManager) GetGitType(s string) model.GitSource {
	switch {
	case strings.Contains(s, string(model.Gitlab)):
		return model.Gitlab
	}
	return model.NotImplemented
}

func (g *GitManager) GetWebhookUrl(domain string, userId int64) (string, error) {
	gs := g.GetGitType(domain)
	api, ok := g.sourceToService[gs]
	if !ok {
		return "", errors.New(fmt.Sprintf("Для git с доменом %s не нашлось обработчика", domain))
	}

	url, _ := api.GetWebhookUrl(domain, userId)
	url += "/"
	urn, err := g.user.GetURN(userId)
	if err != nil {
		return "", err
	}

	if len(urn) > 0 {
		return url + urn, nil
	}

	r := rand.New(rand.NewSource(userId))
	urn = generate(5 + r.Intn(5))
	err = g.user.SetURN(userId, urn)
	tries := 5

	for err != nil || tries > 0 {
		tries--
		urn = generate(5 + r.Intn(5))
		err = g.user.SetURN(userId, urn)
	}
	if tries < 0 {
		return "", errors.New("DB error")
	}

	return url + urn, nil
}

var alphabet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generate(size int) string {
	b := make([]byte, size)
	rand.Read(b)
	for i := 0; i < size; i++ {
		b[i] = alphabet[b[i]/5]
	}
	return *(*string)(unsafe.Pointer(&b))
}
