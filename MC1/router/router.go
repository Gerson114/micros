package router

import (
	"api-go/api"
	"net/http"
)

func Router() {

	http.HandleFunc("/vote", api.CreateVote)
}
