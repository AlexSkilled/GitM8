package entity

import (
	"gitlab-tg-bot/internal/model"
)

type Tickets []Ticket

type Ticket struct {
	tableName          struct{} `pg:"ticket"`
	Id                 int32
	MaintainerGitlabId int64
	RepositoryId       string
	HookTypes          map[string]interface{}
}

func (t *Ticket) FromDto(ticket model.Ticket) {
	t.MaintainerGitlabId = ticket.MaintainerGitlabId
	t.HookTypes = ticket.HookTypes
	t.RepositoryId = ticket.RepositoryId
}
