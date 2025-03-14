-- +goose Up
CREATE TABLE
	boards (
		id UUID PRIMARY KEY,
		game_id UUID references games (id) ON DELETE CASCADE NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

CREATE TABLE
	boxes (
		id UUID PRIMARY KEY,
		color VARCHAR(256) NOT NULL,
		pos_x INTEGER NOT NULL,
		pos_y INTEGER NOT NULL,
		game_id UUID references games (id) ON DELETE CASCADE NOT NULL,
		board_id UUID references boards (id) ON DELETE CASCADE NOT NULL,
		highlight BOOLEAN NOT NULL,
		figure_id UUID references figure_cards (id) ON DELETE CASCADE,
		figure_type VARCHAR(256)
	);

-- +goose Down
DROP TABLE IF EXISTS boxes;

DROP TABLE IF EXISTS boards;