package entity

import (
	"gitlab-tg-bot/service/model"
)

type TicketChatIds []TicketChatId

func (t *TicketChatIds) ToDto() []model.TicketChatId {
	out := make([]model.TicketChatId, len(*t))
	for i, item := range *t {
		out[i] = item.ToDto()
	}
	return out
}

type TicketChatId struct {
	tableName   struct{} `pg:"tickets_chat_id"`
	ChatId      int64
	TicketId    int32
	IsActive    bool
	IsNotifying bool
}

func (t *TicketChatId) ToDto() model.TicketChatId {
	return model.TicketChatId{
		TicketId: t.TicketId,
		ChatId:   t.ChatId,
	}
}
