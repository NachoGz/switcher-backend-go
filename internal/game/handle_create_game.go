package game

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	gameState "github.com/NachoGz/switcher-backend-go/internal/game_state"
	"github.com/NachoGz/switcher-backend-go/internal/middleware"
	"github.com/NachoGz/switcher-backend-go/internal/player"
	"github.com/NachoGz/switcher-backend-go/internal/utils"
	"github.com/google/uuid"
)

func HandleCreateGame(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating game")
	// Get config from context
	cfg := middleware.GetConfig(r)
	if cfg == nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error", nil)
		return
	}

	type GameRequest struct {
		Game   Game          `json:"game"`
		Player player.Player `json:"player"`
	}

	decoder := json.NewDecoder(r.Body)
	params := GameRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	log.Println(params)
	if params.Game.Password != nil {
		hashedPasswd, err := utils.HashPassword(*params.Game.Password)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
			return
		}
		params.Game.Password = &hashedPasswd
	} else {
		params.Game.IsPrivate = false
	}

	// Create game
	game, err := cfg.DB.CreateGame(r.Context(), database.CreateGameParams{
		ID:         uuid.New(),
		Name:       params.Game.Name,
		MaxPlayers: int32(params.Game.MaxPlayers),
		MinPlayers: int32(params.Game.MinPlayers),
		IsPrivate:  params.Game.IsPrivate,
		Password:   sql.NullString{String: *params.Game.Password, Valid: params.Game.Password != nil},
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating game", err)
		return
	}
	log.Println("Game created")

	// Create game state
	game_state, err := cfg.DB.CreateGameState(r.Context(), database.CreateGameStateParams{
		ID:              uuid.New(),
		State:           gameState.WAITING,
		GameID:          uuid.NullUUID{UUID: game.ID, Valid: true},
		CurrentPlayerID: uuid.NullUUID{UUID: uuid.Nil, Valid: false},
		ForbiddenColor:  sql.NullString{String: "", Valid: false},
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating game state", err)
		return
	}
	log.Println("Game state created")

	// Create creator player
	new_player, err := cfg.DB.CreatePlayer(r.Context(), database.CreatePlayerParams{
		ID:          uuid.New(),
		Name:        params.Player.Name,
		Turn:        sql.NullString{String: params.Player.Turn, Valid: params.Player.Turn != ""},
		GameID:      uuid.NullUUID{UUID: game.ID, Valid: true},
		GameStateID: uuid.NullUUID{UUID: game_state.ID, Valid: true},
		Host:        params.Player.Host,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating player", err)
		return
	}
	log.Println("Player created")

	type Response struct {
		Game      Game                `json:"game"`
		GameState gameState.GameState `json:"game_state"`
		Player    player.Player       `json:"player"`
	}

	res := Response{
		Game: Game{
			ID:         game.ID,
			Name:       game.Name,
			MaxPlayers: int(game.MaxPlayers),
			MinPlayers: int(game.MinPlayers),
			IsPrivate:  game.IsPrivate,
		},
		GameState: gameState.GameState{
			ID:              game_state.ID,
			State:           game_state.State,
			GameID:          game_state.GameID.UUID,
			CurrentPlayerID: game_state.CurrentPlayerID.UUID,
			ForbiddenColor:  game_state.ForbiddenColor,
		},
		Player: player.Player{
			ID:          new_player.ID,
			Name:        new_player.Name,
			Turn:        new_player.Turn.String,
			GameID:      new_player.GameID.UUID,
			GameStateID: new_player.GameStateID.UUID,
			Host:        new_player.Host,
		},
	}

	utils.RespondWithJSON(w, http.StatusCreated, res)
}
