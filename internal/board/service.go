package board

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

// Service handles all board-related operations
type Service struct {
	boardRepo BoardRepository
}

// NewService creates a new board service
func NewService(
	boardRepo BoardRepository,

) *Service {
	return &Service{
		boardRepo: boardRepo,
	}
}

// Ensure Service implements BoardService
var _ BoardService = (*Service)(nil)

func (s *Service) ConfigureBoard(ctx context.Context, gameID uuid.UUID) error {
	// Check if a board hasn't been created for this game
	existingBoard, err := s.boardRepo.GetBoard(ctx, gameID)
	if err != nil && existingBoard.ID != uuid.Nil {
		return fmt.Errorf("board already exists for game %v", gameID)
	} else if !errors.Is(err, sql.ErrNoRows) {
		// another database error ocurred
		return fmt.Errorf("error checking for existing board: %w", err)
	}

	// Create a new board
	newBoard, err := s.boardRepo.CreateBoard(ctx, database.CreateBoardParams{
		ID:     uuid.New(),
		GameID: gameID,
	})
	if err != nil {
		return fmt.Errorf("error creating new board: %v", err)
	}

	// Create a list with the colors of the boxes (9 of each)
	colors := make([]ColorEnum, 0, 36)
	for i := 0; i < 9; i++ {
		colors = append(colors, RED)
		colors = append(colors, GREEN)
		colors = append(colors, BLUE)
		colors = append(colors, YELLOW)
	}

	// Shuffle the colors
	rand.Shuffle(len(colors), func(i, j int) {
		colors[i], colors[j] = colors[j], colors[i]
	})

	// Create each box
	for i, color := range colors {
		posX := i % 6
		posY := i / 6

		_, err := s.boardRepo.AddBoxToBoard(ctx, database.AddBoxToBoardParams{
			ID:        uuid.New(),
			BoardID:   newBoard.ID,
			GameID:    gameID,
			Color:     string(color),
			PosX:      int32(posX),
			PosY:      int32(posY),
			Highlight: false,
		})
		if err != nil {
			return fmt.Errorf("failed to add box at position (%d, %d): %w", posX, posY, err)
		}
	}
	return nil
}
