package service

import (
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"

	"strings"

	gapi "github.com/plouc/go-gitlab-client/gitlab"
)

type UserService struct {
	UserProvider interfaces.UserProvider
}

var _ interfaces.UserService = (*UserService)(nil)

func NewUserService(provider interfaces.ProviderStorage) *UserService {
	return &UserService{
		UserProvider: provider.User(),
	}
}

func (u *UserService) GetWithGitlabUsersById(tgId int64) (model.User, error) {
	return u.UserProvider.GetWithGitlabUsers(tgId)
}

func (u *UserService) Register(user model.User) (model.User, error) {
	return u.UserProvider.Create(user)
}

func (u *UserService) AddGitAccount(tgId int64, gitAccount model.GitUser) error {
	if !strings.HasPrefix(gitAccount.Domain, "http") {
		gitAccount.Domain = "https://" + gitAccount.Domain
	}

	if !strings.HasSuffix(gitAccount.Domain, "/") {
		gitAccount.Domain += "/"
	}

	client := gapi.NewGitlab(gitAccount.Domain, StandardApiLevel, gitAccount.Token)
	user, _, err := client.CurrentUser()
	if err != nil {
		// тк пока работаем только с гитлабом, оставляем обработку так,
		// в будущем при добавлении пользователя будет посылаться несколько запросов
		// успешный будет определять с каким гитом (github/gitlab/bitbucket/etc...) будет работать
		// этот токен
		return err
	}
	gitAccount.Email = user.Email
	gitAccount.Username = user.Username
	return u.UserProvider.AddGit(tgId, gitAccount)
}
