package figureCard

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/NachoGz/switcher-backend-go/internal/player"
	"github.com/google/uuid"
)

// Service handles all movement cards related operations
type Service struct {
	figureCardRepo FigureCardRepository
	playerRepo     player.PlayerRepository
}

// NewService creates a new movement cards service
func NewService(
	figureCardRepo FigureCardRepository,
	playerRepo player.PlayerRepository,
) *Service {
	return &Service{
		figureCardRepo: figureCardRepo,
		playerRepo:     playerRepo,
	}
}

// Ensure Service implements FigureCardService
var _ FigureCardService = (*Service)(nil)

func (s *Service) CreateFigureCardDeck(ctx context.Context, gameID uuid.UUID) error {
	// Create a list with card types
	var hardCards []TypeEnum
	var easyCards []TypeEnum

	typesList := GetAllCardTypes()
	// Filter cards based on their name
	for _, card := range typesList {
		cardName := string(card)

		// Hard cards start with "FIG"
		if strings.HasPrefix(cardName, "FIG") && !strings.HasPrefix(cardName, "FIGE") {
			hardCards = append(hardCards, card)
		}

		// Easy cards start with "FIGE"
		if strings.HasPrefix(cardName, "FIGE") {
			easyCards = append(easyCards, card)
		}
	}

	players, err := s.playerRepo.GetPlayers(ctx, gameID)
	if err != nil {
		return fmt.Errorf("failed to get players: %w", err)
	}

	if len(players) == 0 {
		return errors.New("there are no players")
	}

	hardCardsPerPlayer := 36 / len(players)
	easyCardsPerPlayer := 14 / len(players)

	for _, player := range players {
		// Shuffle the lists
		rand.Shuffle(len(hardCards), func(i, j int) {
			hardCards[i], hardCards[j] = hardCards[j], hardCards[i]
		})
		rand.Shuffle(len(easyCards), func(i, j int) {
			easyCards[i], easyCards[j] = easyCards[j], easyCards[i]
		})

		playerCards := hardCards[:hardCardsPerPlayer]
		playerCards = append(playerCards, easyCards[:easyCardsPerPlayer]...)

		rand.Shuffle(len(playerCards), func(i, j int) {
			playerCards[i], playerCards[j] = playerCards[j], playerCards[i]
		})

		show := true
		for i, figure := range playerCards {
			if i == SHOW_LIMIT {
				show = false
			}
			_, err := s.figureCardRepo.CreateFigureCard(ctx, database.CreateFigureCardParams{
				ID:          uuid.New(),
				Show:        show,
				PlayerID:    player.ID,
				GameID:      gameID,
				Type:        string(figure),
				Blocked:     false,
				SoftBlocked: false,
			})
			if err != nil {
				return fmt.Errorf("failed to create figure card: %w", err)
			}
		}
	}
	return nil
}
