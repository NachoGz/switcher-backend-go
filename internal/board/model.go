package board

import (
	"github.com/NachoGz/switcher-backend-go/internal/figureCard"
	"github.com/google/uuid"
)

type ColorEnum string

const (
	RED    ColorEnum = "RED"
	GREEN  ColorEnum = "GREEN"
	BLUE   ColorEnum = "BLUE"
	YELLOW ColorEnum = "YELLOW"
)

type Box struct {
	ID          uuid.UUID `json:"id"`
	Color       ColorEnum `json:"color"`
	PosX        int       `json:"pos_x"`
	PosY        int       `json:"pos_y"`
	Highlighted bool      `json:"highlighted"`
}

type BoxOut struct {
	Color       ColorEnum            `json:"color"`
	PosX        int                  `json:"pos_x"`
	PosY        int                  `json:"pos_y"`
	Highlighted bool                 `json:"highlighted"`
	FigureID    *uuid.UUID           `json:"figure_id"`
	FigureType  *figureCard.TypeEnum `json:"figure_type"`
}

type Board struct {
	ID     uuid.UUID `json:"id"`
	GameID uuid.UUID `json:"game_id"`
}

type BoardAndBoxesOut struct {
	GameID        uuid.UUID  `json:"game_id"`
	BoardID       uuid.UUID  `json:"board_id"`
	Boxes         [][]BoxOut `json:"boxes"`
	FormedFigures [][]BoxOut `json:"formed_figures"`
}

type BoardPosition struct {
	PosX int `json:"pos_x"`
	PosY int `json:"pos_y"`
}
