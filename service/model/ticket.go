package model

type Ticket struct {
	Id                 int32
	MaintainerGitlabId int64
	RepositoryId       string
	HookTypes          map[string]interface{}
	ChatIds            []int64
}
