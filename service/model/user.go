package model

type User struct {
	Id         int64
	Name       string
	TgUsername string
	Locale     string

	Gitlabs []GitUser
}
