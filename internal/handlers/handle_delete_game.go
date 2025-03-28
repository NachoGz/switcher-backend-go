package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/NachoGz/switcher-backend-go/internal/utils"
	"github.com/google/uuid"
)

func (h *GameHandlers) HandleDeleteGame(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting game...")

	gameID, err := uuid.Parse(r.PathValue("gameID"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Couldn't parse game ID", err)
		return
	}

	err = h.gameService.DeleteGame(context.Background(), gameID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't delete game", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusNoContent, map[string]string{
		"message": "Game deleted successfully",
	})

	h.wsHub.BroadcastEvent(uuid.Nil, "GAMES_LIST_UPDATE")

}
