package figureCard

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

type TypeEnum string

const (
	FIG01  TypeEnum = "FIG01"
	FIG02  TypeEnum = "FIG02"
	FIG03  TypeEnum = "FIG03"
	FIG04  TypeEnum = "FIG04"
	FIG05  TypeEnum = "FIG05"
	FIG06  TypeEnum = "FIG06"
	FIG07  TypeEnum = "FIG07"
	FIG08  TypeEnum = "FIG08"
	FIG09  TypeEnum = "FIG09"
	FIG10  TypeEnum = "FIG10"
	FIG11  TypeEnum = "FIG11"
	FIG12  TypeEnum = "FIG12"
	FIG13  TypeEnum = "FIG13"
	FIG14  TypeEnum = "FIG14"
	FIG15  TypeEnum = "FIG15"
	FIG16  TypeEnum = "FIG16"
	FIG17  TypeEnum = "FIG17"
	FIG18  TypeEnum = "FIG18"
	FIGE01 TypeEnum = "FIGE01"
	FIGE02 TypeEnum = "FIGE02"
	FIGE03 TypeEnum = "FIGE03"
	FIGE04 TypeEnum = "FIGE04"
	FIGE05 TypeEnum = "FIGE05"
	FIGE06 TypeEnum = "FIGE06"
	FIGE07 TypeEnum = "FIGE07"
)

func GetAllCardTypes() []TypeEnum {
	return []TypeEnum{
		FIG01, FIG02, FIG03, FIG04, FIG05, FIG06, FIG07, FIG08, FIG09,
		FIG10, FIG11, FIG12, FIG13, FIG14, FIG15, FIG16, FIG17, FIG18,
		FIGE01, FIGE02, FIGE03, FIGE04, FIGE05, FIGE06, FIGE07,
	}
}

type DifficultyEnum string

const (
	EASY DifficultyEnum = "easy"
	HARD DifficultyEnum = "hard"
)

const SHOW_LIMIT = 3

type FigureCard struct {
	ID          uuid.UUID `json:"id"`
	Type        TypeEnum  `json:"type"`
	Difficulty  string    `json:"difficulty"`
	Show        bool      `json:"show"`
	PlayerID    uuid.UUID `json:"player_id"`
	GameID      uuid.UUID `json:"game_id"`
	Blocked     bool      `json:"blocked"`
	SoftBlocked bool      `json:"soft_blocked"`
}

// DBToModel converts a database movement card to a model movement card
func (s *Service) DBToModel(ctx context.Context, dbFigureCard database.FigureCard) FigureCard {
	return FigureCard{
		ID:          dbFigureCard.ID,
		Type:        TypeEnum(dbFigureCard.Type),
		Difficulty:  dbFigureCard.Difficulty,
		Show:        dbFigureCard.Show,
		PlayerID:    dbFigureCard.PlayerID,
		GameID:      dbFigureCard.GameID,
		Blocked:     dbFigureCard.Blocked,
		SoftBlocked: dbFigureCard.SoftBlocked,
	}
}
