package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/NachoGz/switcher-backend-go/internal/utils"
	"github.com/google/uuid"
)

func (h *GameHandlers) HandlerGetWinner(w http.ResponseWriter, r *http.Request) {
	gameID, err := uuid.Parse(r.PathValue("gameID"))
	log.Println(gameID.String())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Couldn't parse game ID", err)
		return
	}

	log.Println("Getting winner from game: ", gameID)

	winner, err := h.playerService.GetWinner(context.Background(), gameID)
	if err != nil && winner != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't fetch winner", err)
		return
	}

	if winner == nil {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{
			"message": "There is no winner",
		})
		return
	} else {
		utils.RespondWithJSON(w, http.StatusOK, winner)
		return
	}
}
