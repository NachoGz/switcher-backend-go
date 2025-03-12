package movementCard

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand/v2"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/NachoGz/switcher-backend-go/internal/player"
	"github.com/google/uuid"
)

// Service handles all movement cards related operations
type Service struct {
	movementCardRepo MovementCardRepository
	playerRepo       player.PlayerRepository
}

// NewService creates a new movement cards service
func NewService(
	movementCardRepo MovementCardRepository,
	playerRepo player.PlayerRepository,
) *Service {
	return &Service{
		movementCardRepo: movementCardRepo,
		playerRepo:       playerRepo,
	}
}

func (s *Service) CreateMovementDeck(ctx context.Context, gameID uuid.UUID) error {
	// Create a list with the types of movement cards
	typesList := make([]TypeEnum, 0, 40)
	for i := 0; i < 6; i++ {
		typesList = append(typesList, DIAGONAL_CONT)
	}

	for i := 0; i < 6; i++ {
		typesList = append(typesList, DIAGONAL_SPA)
	}

	for i := 0; i < 6; i++ {
		typesList = append(typesList, L_RIGHT)
	}

	for i := 0; i < 6; i++ {
		typesList = append(typesList, L_LEFT)
	}

	for i := 0; i < 5; i++ {
		typesList = append(typesList, LINEAR_LAT)
	}

	for i := 0; i < 5; i++ {
		typesList = append(typesList, LINEAR_CONT)
	}

	for i := 0; i < 6; i++ {
		typesList = append(typesList, LINEAR_SPA)
	}

	// Shuffle the list
	rand.Shuffle(len(typesList), func(i, j int) {
		typesList[i], typesList[j] = typesList[j], typesList[i]
	})

	// Create movement cards
	for i, cardType := range typesList {
		_, err := s.movementCardRepo.CreateMovementCard(ctx, database.CreateMovementCardParams{
			ID:          uuid.New(),
			Description: "",
			Used:        false,
			GameID:      gameID,
			Type:        string(cardType),
			Position:    sql.NullInt32{Int32: int32(i)},
		})
		if err != nil {
			return fmt.Errorf("failed to create movement card: %w", err)
		}
	}

	// assign 3 to each player
	players, err := s.playerRepo.GetPlayers(ctx, gameID)
	if err != nil {
		return fmt.Errorf("failed to get players: %w", err)
	}

	for _, player := range players {
		movDeck, err := s.movementCardRepo.GetMovementDeck(context.Background(), gameID)
		if err != nil {
			return fmt.Errorf("failed to get movement deck: %w", err)
		}

		// Check if there are enough cards in deck
		if len(movDeck) < 3 {
			return fmt.Errorf("not enough cards in deck to assign to player %s", player.ID)
		}

		// Shuffle the deck
		rand.Shuffle(len(movDeck), func(i, j int) {
			movDeck[i], movDeck[j] = movDeck[j], movDeck[i]
		})

		// Take the first three cards
		assignedMovCards := movDeck[:3]

		for _, card := range assignedMovCards {
			err := s.movementCardRepo.AssignMovementCard(context.Background(), database.AssignMovementCardParams{
				ID:       card.ID,
				PlayerID: uuid.NullUUID{UUID: player.ID, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("failed to assign card %s to player %s: %w", card.ID, player.ID, err)
			}
		}
	}
	return nil
}
