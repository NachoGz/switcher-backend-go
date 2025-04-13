package partialMovements

import (
	"context"
	"database/sql"

	"github.com/NachoGz/switcher-backend-go/internal/board"
	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/NachoGz/switcher-backend-go/internal/movementCard"
	"github.com/google/uuid"
)

// Service handles all partial movements operations
type Service struct {
	db               *sql.DB
	partialMovRepo   PartialMovementRepository
	boardService     board.BoardService
	movementCardRepo movementCard.MovementCardRepository
}

// NewService creates a new partial movements servzicez
func NewService(partialMovRepo PartialMovementRepository, boardService board.BoardService,
	movementCardRepo movementCard.MovementCardRepository) PartialMovementService {
	return &Service{
		partialMovRepo:   partialMovRepo,
		boardService:     boardService,
		movementCardRepo: movementCardRepo,
	}
}

var _ PartialMovementService = (*Service)(nil)

func (s *Service) RevertPartialMovements(ctx context.Context, db *sql.DB, gameID uuid.UUID, playerID uuid.UUID) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	boardServiceTx := s.boardService.WithTx(tx)

	partialMovements, err := s.partialMovRepo.GetPartialMovementsByPlayer(ctx, database.GetPartialMovementsByPlayerParams{
		GameID:   gameID,
		PlayerID: playerID,
	})
	if err != nil {
		return err
	}

	for _, movement := range partialMovements {
		// Use swap colors to revert
		posFrom := board.BoardPosition{PosX: int(movement.PosFromX), PosY: int(movement.PosFromY)}
		posTo := board.BoardPosition{PosX: int(movement.PosToX), PosY: int(movement.PosToY)}

		if err := boardServiceTx.SwapColors(ctx, gameID, posFrom, posTo); err != nil {
			return err
		}

		// Mark card as not used
		s.movementCardRepo.MarkCardInPlayerHand(ctx, movement.MovementCardID)

		// Delete partial movement
		s.partialMovRepo.UndoMovementByID(ctx, movement.ID)
	}

	return tx.Commit()
}
