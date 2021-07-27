package provider

import (
	"gitlab-tg-bot/data/entity"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"

	"github.com/sirupsen/logrus"

	"github.com/go-pg/pg/v9"
)

type TicketProvider struct {
	*pg.DB
}

var _ interfaces.TicketProvider = (*TicketProvider)(nil)

func NewTicket(conn *pg.DB) *TicketProvider {
	return &TicketProvider{
		conn,
	}
}

func (t *TicketProvider) AddTicket(ticket model.Ticket) (ticketId int32, err error) {

	tx, err := t.DB.Begin()
	defer func() {
		if err != nil {
			logrus.Errorf("Ошибка при попытке создать тикет. %v", err)
			err = tx.Rollback()
			if err != nil {
				logrus.Errorf("Ошибка при откате транзакции. %v", err)
			}
		}
	}()

	var ticketEntity entity.Ticket
	ticketEntity.FromDto(ticket)
	_, err = tx.Model(&ticketEntity).Insert()
	if err != nil {
		return 0, err
	}

	ticketChatIds := make([]entity.TicketChatId, len(ticket.ChatIds))
	for i := range ticketChatIds {
		ticketChatIds[i].TicketId = ticketEntity.Id
		ticketChatIds[i].ChatId = ticket.ChatIds[i]
		ticketChatIds[i].IsActive = true
	}
	_, err = tx.Model(&ticketChatIds).Insert()
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return ticketEntity.Id, nil
}

func (t *TicketProvider) AddTicketToChat(chatId int64, ticketId int32) error {
	usersTickets := entity.TicketChatId{
		ChatId:   chatId,
		TicketId: ticketId,
		IsActive: true,
	}

	_, err := t.DB.Model(&usersTickets).Insert()

	if err != nil {
		logrus.Errorln(err)
		return err
	}
	return nil
}

func (t *TicketProvider) GetTicketsToSend(repoId string, hookType model.GitHookType) ([]model.TicketChatId, error) {

	return nil, nil
}
