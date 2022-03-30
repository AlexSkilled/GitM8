package model

type Repository struct {
	Id     string
	Name   string
	WebUrl string
}

type Group struct {
	Id           string
	Name         string
	WebUrl       string
	Repositories []Repository
}
