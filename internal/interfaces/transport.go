package interfaces

import "gitlab-tg-bot/service/model"

type GitMapper interface {
	ToModel() model.GitEvent
}
