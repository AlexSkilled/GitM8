package internal

import "errors"

var (
	ErrorNoTickets    = errors.New("отсутствуют билеты по этому хуку")
	ErrorAccessDenied = errors.New("время жизни токена закончилось или токен негоден")
)
