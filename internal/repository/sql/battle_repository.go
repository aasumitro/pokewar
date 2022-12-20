package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aasumitro/pokewar/domain"
	"github.com/aasumitro/pokewar/pkg/appconfigs"
	"time"
)

type battleSQLRepository struct {
	db *sql.DB
}

func (repo *battleSQLRepository) Count(ctx context.Context) int {
	var total int

	q := "SELECT COUNT(*) FROM battles"
	if err := repo.db.QueryRowContext(ctx, q).Scan(&total); err != nil {
		total = 0
	}

	return total
}

func (repo *battleSQLRepository) All(ctx context.Context, args ...string) (data []*domain.Battle, err error) {
	q := "SELECT b.id as id, b.started_at as started_at, b.ended_at as ended_at, "
	q += "CAST((SELECT json_group_array(json_object('id', bl.id, 'battle_id', bl.battle_id, 'description', bl.description, "
	q += "'created_at', bl.created_at)) FROM battle_logs as bl where bl.battle_id = b.id) AS CHAR) as battle_logs, "
	q += "CAST((SELECT json_group_array(json_object('id', bp.id, 'battle_id', bp.battle_id, 'monster_id', bp.monster_id, "
	q += "'eliminated_at', bp.eliminated_at, 'annulled_at', bp.annulled_at, 'rank', bp.rank, 'point', bp.point, "
	q += "'name', m.name)) FROM battle_players as bp join monsters as m on bp.monster_id = "
	q += "m.id where bp.battle_id = b.id) AS CHAR) as battle_players FROM battles as b "
	if len(args) > 0 {
		for _, arg := range args {
			q += fmt.Sprintf("%s ", arg)
		}
	}

	rows, err := repo.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) { _ = rows.Close() }(rows)

	for rows.Next() {
		var battle domain.BattleEntity
		var battleLogs []domain.Log
		var battlePlayers []domain.Player

		if err := rows.Scan(
			&battle.ID, &battle.StartedAt, &battle.EndedAt,
			&battle.Logs, &battle.Players,
		); err != nil {
			return nil, err
		}

		_ = json.Unmarshal([]byte(battle.Logs), &battleLogs)
		_ = json.Unmarshal([]byte(battle.Players), &battlePlayers)
		data = append(data, &domain.Battle{
			ID:        battle.ID,
			StartedAt: battle.StartedAt,
			EndedAt:   battle.EndedAt,
			Logs:      battleLogs,
			Players:   battlePlayers,
		})
	}

	return data, nil
}

func (repo *battleSQLRepository) Create(ctx context.Context, param *domain.Battle) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	qb := "INSERT INTO battles (started_at, ended_at) VALUES (?, ?) RETURNING id"
	newBattle, err := tx.Exec(qb, param.StartedAt, param.EndedAt)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	battleId, err := newBattle.LastInsertId()
	if err != nil {
		return err
	}

	ql := "INSERT INTO battle_logs (battle_id, description, created_at) VALUES"
	now := time.Now().Unix()
	for i, log := range param.Logs {
		ql += fmt.Sprintf(" (%d, %s, %d)", battleId, log.Description, now)
		if i != (len(param.Logs) - 1) {
			ql += ","
		}
	}
	_, err = tx.Exec(ql)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	qp := "INSERT INTO battle_players (battle_id, monster_id, eliminated_at, rank, point) VALUES"
	for i, player := range param.Players {
		ql += fmt.Sprintf(" (%d, %d, %d, %d, %d)",
			battleId, player.MonsterID, player.EliminatedAt,
			player.Rank, player.Point)
		if i != (len(param.Logs) - 1) {
			ql += ","
		}
	}
	_, err = tx.Exec(qp)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// UpdatePlayer
// When annulled make point 0
// update from rank ...2 add + 1
func (repo *battleSQLRepository) UpdatePlayer(ctx context.Context, id int) (annulledAt int64, err error) {
	var annulledTime int64

	q := "UPDATE battle_players SET annulled_at = ? WHERE id = ? RETURN annulled_at"
	if err := repo.db.QueryRowContext(ctx, q, time.Now().Unix(), id).Scan(&annulledTime); err != nil {
		return 0, err
	}

	return annulledTime, nil
}

func NewBattleSQLRepository() domain.IBattleRepository {
	return &battleSQLRepository{db: appconfigs.DbPool}
}
