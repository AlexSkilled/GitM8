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

func (t *TicketProvider) AddTicket(ticket model.Ticket) (int32, error) {
	var ticketEntity entity.Ticket
	ticketEntity.FromDto(ticket)

	_, err := t.Model(&ticketEntity).Insert()
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	return ticketEntity.Id, nil
}

func (t *TicketProvider) AddTicketToUser(tgUserId int64, ticketId int32) error {
	_, err := t.DB.Query(nil, `
		INSERT INTO 
		users_tickets (tg_user_id, ticket_id)
			   VALUES (         ?,         ?)`,
		tgUserId, ticketId)
	if err != nil {
		logrus.Errorln(err)
		return err
	}
	return nil
}
