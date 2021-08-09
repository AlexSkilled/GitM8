package entity

import "gitlab-tg-bot/service/model"

type MessagePattern struct {
	HookType           model.GitHookType
	Lang               string
	Message            string
	AdditionalPatterns map[string]string
}
