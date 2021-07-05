package service

import (
	"gitlab-tg-bot/data"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/internal/model"
)

type UserService struct {
	UserProvider interfaces.UserProvider
}

func NewUserService(provider data.ProviderStorage) UserService {
	return UserService{UserProvider: provider.UserProvider}
}

func (u *UserService) CreateUser(user model.User) error {
	err := u.UserProvider.Create(user)
	return err
}

func (u *UserService) GetUser(id int32) (model.User, error) {
	user, err := model.User{
		Id:         1,
		Name:       "Bukov Alexandr",
		TgUsername: "some u n",
		TgId:       1,
		Gitlabs: []model.GitlabUser{
			{
				Id:     1,
				Token:  "DFX3ppBJb7qdBsjz3DsH",
				Domain: "https://gitlab.ru/",
				Email:  "mail@mail.ru",
			},
		},
	}, error(nil)
	return user, err
	//return u.UserProvider.Get(id)
}

//func (u *UserService) GetProjects() ([]model.Project, error) {
//	user, err := u.GetUser(1)
//	if err != nil {
//		return nil, err
//	}
//	projects, _, err := user.GetGitlabClient().Projects(&gitlab.ProjectsOptions{Membership: true})
//	if err != nil {
//		return nil, err
//	}
//
//	out := make([]model.Project, len(projects.Items))
//	for i, item := range projects.Items {
//		out[i] = model.Project{
//			Id:   int32(item.Id),
//			Name: item.Name,
//		}
//	}
//	return out, nil
//}
