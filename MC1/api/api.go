package api

import (
	"api-go/habbit"
	"api-go/models"
	"encoding/json"
	"log"
	"net/http"
)

func CreateVote(w http.ResponseWriter, r *http.Request) {
	log.Printf("📥 Recebendo voto - Método: %s", r.Method)

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var vote models.Vote

	validar := json.NewDecoder(r.Body).Decode(&vote)
	if validar != nil {
		log.Printf("❌ JSON inválido: %v", validar)
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	log.Printf("📤 Publicando voto do usuário %d", vote.Id)

	//enviar para o produce
	enviar := habbit.PublishVote(vote)
	if enviar != nil {
		log.Printf("❌ Erro ao publicar: %v", enviar)
		http.Error(w, "Erro ao publicar mensagem", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Voto do usuário %d enviado para fila", vote.Id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message":"Voto enviado para fila"}`))
}
