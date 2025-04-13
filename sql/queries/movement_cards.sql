-- name: CreateMovementCard :one
INSERT INTO
	movement_cards (
		id,
		description,
		used,
		player_id,
		game_id,
		type,
		position
	)
VALUES
	(
		$1, 
		$2,
		$3,
		$4,
		$5,
		$6,
		$7
	)
RETURNING *;

-- name: GetMovementCardDeck :many
SELECT *
FROM movement_cards
WHERE game_id = $1 AND player_id IS NULL;

-- name: AssignMovementCard :exec
UPDATE movement_cards
SET player_id = $2
WHERE id=$1;

-- name: MarkCardInPlayerHand :exec
UPDATE movement_cards
SET used = false
WHERE id = $1;