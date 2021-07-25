package entity

type TicketChatId struct {
	tableName struct{} `pg:"tickets_chat_id"`
	ChatId    int64
	TicketId  int32
	IsActive  bool
}
