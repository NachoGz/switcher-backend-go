// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: players.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const assignTurnPlayer = `-- name: AssignTurnPlayer :exec
UPDATE players
SET turn=$2, updated_at = NOW()
WHERE id=$1
`

type AssignTurnPlayerParams struct {
	ID   uuid.UUID
	Turn sql.NullString
}

func (q *Queries) AssignTurnPlayer(ctx context.Context, arg AssignTurnPlayerParams) error {
	_, err := q.db.ExecContext(ctx, assignTurnPlayer, arg.ID, arg.Turn)
	return err
}

const countPlayers = `-- name: CountPlayers :one
SELECT COUNT(*) FROM players WHERE game_id = $1
`

func (q *Queries) CountPlayers(ctx context.Context, gameID uuid.UUID) (int64, error) {
	row := q.db.QueryRowContext(ctx, countPlayers, gameID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createPlayer = `-- name: CreatePlayer :one
INSERT INTO players (id, name, turn, game_id, game_state_id, host)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, name, turn, game_id, game_state_id, host, winner, created_at, updated_at
`

type CreatePlayerParams struct {
	ID          uuid.UUID
	Name        string
	Turn        sql.NullString
	GameID      uuid.UUID
	GameStateID uuid.UUID
	Host        bool
}

func (q *Queries) CreatePlayer(ctx context.Context, arg CreatePlayerParams) (Player, error) {
	row := q.db.QueryRowContext(ctx, createPlayer,
		arg.ID,
		arg.Name,
		arg.Turn,
		arg.GameID,
		arg.GameStateID,
		arg.Host,
	)
	var i Player
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Turn,
		&i.GameID,
		&i.GameStateID,
		&i.Host,
		&i.Winner,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPlayerByID = `-- name: GetPlayerByID :one
SELECT id, name, turn, game_id, game_state_id, host, winner, created_at, updated_at
FROM players
WHERE game_id=$1 AND id=$2
`

type GetPlayerByIDParams struct {
	GameID uuid.UUID
	ID     uuid.UUID
}

func (q *Queries) GetPlayerByID(ctx context.Context, arg GetPlayerByIDParams) (Player, error) {
	row := q.db.QueryRowContext(ctx, getPlayerByID, arg.GameID, arg.ID)
	var i Player
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Turn,
		&i.GameID,
		&i.GameStateID,
		&i.Host,
		&i.Winner,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPlayersInGame = `-- name: GetPlayersInGame :many
SELECT id, name, turn, game_id, game_state_id, host, winner, created_at, updated_at
FROM players
WHERE game_id=$1
`

func (q *Queries) GetPlayersInGame(ctx context.Context, gameID uuid.UUID) ([]Player, error) {
	rows, err := q.db.QueryContext(ctx, getPlayersInGame, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Player
	for rows.Next() {
		var i Player
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Turn,
			&i.GameID,
			&i.GameStateID,
			&i.Host,
			&i.Winner,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getWinner = `-- name: GetWinner :one
SELECT id, name, turn, game_id, game_state_id, host, winner, created_at, updated_at
FROM players
WHERE game_id = $1 AND winner = true limit 1
`

func (q *Queries) GetWinner(ctx context.Context, gameID uuid.UUID) (Player, error) {
	row := q.db.QueryRowContext(ctx, getWinner, gameID)
	var i Player
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Turn,
		&i.GameID,
		&i.GameStateID,
		&i.Host,
		&i.Winner,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
