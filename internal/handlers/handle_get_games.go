package handlers

import (
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/NachoGz/switcher-backend-go/internal/utils"
	"github.com/google/uuid"
)

func (h *GameHandlers) HandleGetGames(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting games")

	// Parse pagination parameters
	page := 1 // Default page
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		pageVal, err := strconv.Atoi(pageStr)
		if err != nil || pageVal < 1 {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid page", err)
			return
		}
		page = pageVal
	}

	limit := 5 // Default limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limitVal, err := strconv.Atoi(limitStr)
		if err != nil || limitVal < 1 {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid limit", err)
			return
		}
		limit = limitVal
	}

	// filters
	numPlayers := 0 // 0 means no filter
	if numPlayersStr := r.URL.Query().Get("num_players"); numPlayersStr != "" {
		log.Println("num_players:", numPlayersStr)
		numPlayersVal, err := strconv.Atoi(numPlayersStr)
		if err != nil || numPlayersVal < 0 {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid number of players", err)
			return
		}
		numPlayers = numPlayersVal
	}

	name := r.URL.Query().Get("name")

	// Use service to get games
	games, total, err := h.gameService.GetAvailableGames(r.Context(), numPlayers, page, limit, name)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error getting games", err)
		return
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"total_pages": totalPages,
		"games":       games,
	})
}

func (h *GameHandlers) HandleGetGameByID(w http.ResponseWriter, r *http.Request) {
	gameID, err := uuid.Parse(r.PathValue("gameID"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Couldn't parse game ID", err)
		return
	}

	log.Println("Getting game with ID: ", gameID)

	// Use service to get game
	game, err := h.gameService.GetGameByID(r.Context(), gameID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error getting game", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, game)
}
