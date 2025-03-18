package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/NachoGz/switcher-backend-go/internal/utils"
	"github.com/google/uuid"
)

func (h *PlayerHandlers) HandleGetPlayers(w http.ResponseWriter, r *http.Request) {
	gameID, err := uuid.Parse(r.PathValue("gameID"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Couldn't parse game ID", err)
		return
	}

	log.Println("Getting players from game: ", gameID)

	players, err := h.playerService.GetPlayersInGame(context.Background(), gameID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch players", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, players)
}

func (h *PlayerHandlers) HandleGetPlayer(w http.ResponseWriter, r *http.Request) {
	gameID, err := uuid.Parse(r.PathValue("gameID"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Couldn't parse game ID", err)
		return
	}

	playerID, err := uuid.Parse(r.PathValue("playerID"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Couldn't parse player ID", err)
		return
	}

	log.Printf("Getting player %s from game %s: ", playerID, gameID)

	player, err := h.playerService.GetPlayerByID(context.Background(), gameID, playerID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch player", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, player)
}
