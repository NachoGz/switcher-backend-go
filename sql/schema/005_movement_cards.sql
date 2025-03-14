-- +goose Up
CREATE TABLE
	movement_cards (
		id UUID PRIMARY KEY,
		description VARCHAR(256) NOT NULL,
		used BOOLEAN NOT NULL,
		player_id UUID references players (id) ON DELETE CASCADE,
		game_id UUID references games (id) ON DELETE CASCADE NOT NULL,
		type VARCHAR(256) NOT NULL,
		position INTEGER
	);

-- +goose Down
DROP TABLE IF EXISTS movement_cards;