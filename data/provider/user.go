package provider

import (
	"gitlab-tg-bot/data/entity"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"

	"github.com/sirupsen/logrus"

	"github.com/go-pg/pg/v9"
)

type UserProvider struct {
	*pg.DB
}

var _ interfaces.UserProvider = (*UserProvider)(nil)

func NewUser(conn *pg.DB) interfaces.UserProvider {
	return &UserProvider{conn}
}

func (u *UserProvider) Create(in model.User) (model.User, error) {
	var user entity.User
	user.FromModel(in)

	_, err := u.Model(&user).Insert()
	return user.ToModel(), err
}

func (u *UserProvider) Get(id int64) (model.User, error) {
	var user entity.User

	err := u.Model(&user).Where("id = ?", id).Select()
	return user.ToModel(), err
}

func (u *UserProvider) GetWithGitlabUsers(id int64) (model.User, error) {
	var gitlabEntities entity.GitUsers

	err := u.Model(&gitlabEntities).Where("user_id = ?", id).Select()
	if err != nil {
		// TODO
	}

	gitlabs := gitlabEntities.ToModel()

	user, err := u.Get(id)
	if err != nil {
		return model.User{}, err
	}

	user.Gitlabs = gitlabs
	return user, err
}

func (u *UserProvider) AddGit(userId int64, gitlab model.GitUser) error {
	var gitlabEnt entity.GitUser
	gitlabEnt.FromModel(userId, gitlab)

	_, err := u.Model(&gitlabEnt).Insert()
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (u *UserProvider) Update(user model.User) error {
	var userEnt entity.User
	columns := make([]string, 0, 1)

	if len(user.Locale) > 0 {
		userEnt.Locale = user.Locale
		columns = append(columns, "locale")
	}

	_, err := u.Model(&userEnt).Where("id = ?", user.Id).Column(columns...).Update()
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (u *UserProvider) GetURN(userId int64) (string, error) {
	var user entity.User

	err := u.Model(&user).Where("id = ?", userId).Column("urn").Select()
	return user.Urn, err
}

func (u *UserProvider) SetURN(userId int64, urn string) error {
	user := entity.User{
		Urn: urn,
	}
	_, err := u.Model(&user).Where("id = ?", userId).Column("urn").Update()
	return err
}
