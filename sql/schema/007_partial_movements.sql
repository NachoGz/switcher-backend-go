-- +goose Up
CREATE TABLE
	partial_movements (
		id UUID PRIMARY KEY,
		pos_from_x INTEGER NOT NULL,
		pos_from_y INTEGER NOT NULL,
		pos_to_x INTEGER NOT NULL,
		pos_to_y INTEGER NOT NULL,
		game_id UUID references games (id) ON DELETE CASCADE NOT NULL,
		player_id UUID references players (id) ON DELETE CASCADE NOT NULL,
		movement_card_id UUID references movement_cards (id) ON DELETE CASCADE NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

-- +goose Down
DROP TABLE IF EXISTS partial_movements;