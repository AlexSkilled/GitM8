package model

type User struct {
	Id         int32
	Name       string
	TgUsername string
	TgId       int64

	Gitlabs []GitlabUser
}
