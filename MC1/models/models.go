package models

type Vote struct {
	Id     int    `json:"id"`
	Nome   string `json:"nome"`
	Numero int    `json:"numero"`
	Email  string `json:"email"`
	Votos  string `json:"votos"`
}
