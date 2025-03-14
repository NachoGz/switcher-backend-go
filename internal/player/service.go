package player

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

// Service handles all game-related operations
type Service struct {
	playerRepo PlayerRepository
}

// NewService creates a new game service
func NewService(
	playerRepo PlayerRepository,
) *Service {
	return &Service{
		playerRepo: playerRepo,
	}
}

func (s *Service) CreatePlayer(ctx context.Context, playerData Player) (*Player, error) {
	player, err := s.playerRepo.CreatePlayer(ctx, database.CreatePlayerParams{
		ID:          playerData.ID,
		Name:        playerData.Name,
		Turn:        sql.NullString{String: string(playerData.Turn)},
		GameID:      playerData.GameID,
		GameStateID: playerData.GameStateID,
		Host:        playerData.Host,
	})
	if err != nil {
		return nil, err
	}

	resultPlayer := s.DBToModel(ctx, player)

	return &resultPlayer, nil
}

func (s *Service) GetPlayers(ctx context.Context, gameID uuid.UUID) ([]Player, error) {
	players, err := s.playerRepo.GetPlayers(ctx, gameID)
	if err != nil {
		return nil, err
	}

	returnPlayers := []Player{}
	for _, player := range players {
		returnPlayers = append(returnPlayers, s.DBToModel(ctx, player))
	}

	return returnPlayers, nil
}

// Assign randomly the turns for the players in the game and returrns the id of the first player
func (s *Service) AssignRandomTurns(ctx context.Context, players []Player) (uuid.UUID, error) {
	n := len(players)
	if n == 0 {
		return uuid.Nil, errors.New("there are no players")
	}

	randomTurns := make([]int, n)
	for i := 0; i < n; i++ {
		randomTurns[i] = i + 1
	}

	// Shuffle to randomize order
	rand.Shuffle(n, func(i, j int) {
		randomTurns[i], randomTurns[j] = randomTurns[j], randomTurns[i]
	})

	// Define mapping from turn number to turn enum
	turnMapping := map[int]TurnEnum{
		1: FIRST,
		2: SECOND,
		3: THIRD,
		4: FOURTH,
	}

	var firstPlayer *Player

	// Iterate players and assign random turns
	for i, player := range players {
		turn := randomTurns[i]
		turnEnumVal, ok := turnMapping[turn]
		if !ok {
			return uuid.Nil, fmt.Errorf("invalid turn: %d", turn)
		}

		// Assign turn
		if err := s.playerRepo.AssignTurnPlayer(context.Background(), database.AssignTurnPlayerParams{
			ID:   player.ID,
			Turn: sql.NullString{String: string(turnEnumVal)},
		}); err != nil {
			return uuid.Nil, err
		}

		if turn == 1 {
			firstPlayer = &players[i]
		}
	}

	if firstPlayer == nil {
		return uuid.Nil, errors.New("first player not assigned")
	}

	return firstPlayer.ID, nil
}
