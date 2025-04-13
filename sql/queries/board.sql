-- name: CreateBoard :one
INSERT INTO
	boards (id, game_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetBoard :one
SELECT *
FROM boards
WHERE game_id=$1;

-- name: AddBoxToBoard :one
INSERT INTO
	boxes (id, color, pos_x, pos_y, game_id, board_id, highlight)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetBox :one
SELECT *
FROM boxes
WHERE game_id = $1 AND pos_x = $2 AND pos_y = $3;

-- name: ChangeBoxColor :exec
UPDATE boxes
SET color = $2
WHERE id = $1;
