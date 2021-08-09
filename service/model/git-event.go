package model

const (
	Gitlab GitSource = "gitlab" // на этот URI будут регистрироваться хуки для гита.

	// т.е гит будет слать хуки на http://our.service.web/gitlab
)

type GitSource string

func (g GitSource) GetUri() string {
	return "/" + string(g)
}

type GitEvent struct {
	GitSource GitSource

	ProjectId   string
	ProjectName string

	HookType GitHookType
	SubType  GitHookSubtype

	Payload map[PayloadKey]string

	Link string
}
