-- +goose Up
CREATE TABLE
	game_state (
		id UUID PRIMARY KEY,
		state VARCHAR(255) NOT NULL,
		game_id UUID references games (id) ON DELETE CASCADE NOT NULL,
		current_player_id UUID DEFAULT NULL,
		forbidden_color VARCHAR(255) DEFAULT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

-- +goose Down
DROP TABLE IF EXISTS game_state;