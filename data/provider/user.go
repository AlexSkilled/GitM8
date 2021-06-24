package provider

import (
	"gitlab-tg-bot/data/entity"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/internal/model"

	"github.com/go-pg/pg/v9"
)

type UserProvider struct {
	*pg.DB
}

var _ interfaces.UserProvider = (*UserProvider)(nil)

func NewUser(conn *pg.DB) interfaces.UserProvider {
	return &UserProvider{conn}
}

func (u *UserProvider) Create(in model.User) error {
	var user entity.User
	user.FromModel(in)

	_, err := u.Model(&user).Insert()
	return err
}

func (u *UserProvider) Get(id int32) (model.User, error) {
	var user entity.User

	err := u.Model(&user).Where("id = ?", id).Select()
	return user.ToModel(), err
}

func (u *UserProvider) GetByTelegramId(id int64) (model.User, error) {
	var user entity.User

	err := u.Model(&user).Where("telegram_id = ?", id).Select()
	return user.ToModel(), err
}
