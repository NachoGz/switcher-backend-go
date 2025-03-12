-- +goose Up
CREATE TABLE
	figure_cards (
		id UUID PRIMARY KEY,
		show BOOLEAN NOT NULL,
		difficulty VARCHAR(256) NOT NULL,
		player_id UUID references players (id) ON DELETE CASCADE NOT NULL,
		game_id UUID references games (id) ON DELETE CASCADE NOT NULL,
		type VARCHAR(256) NOT NULL,
		blocked BOOLEAN NOT NULL,
		soft_blocked BOOLEAN NOT NULL
	);

-- +goose Down
DROP TABLE IF EXISTS figure_cards;