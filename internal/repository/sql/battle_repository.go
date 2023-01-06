package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aasumitro/pokewar/configs"
	"github.com/aasumitro/pokewar/domain"
	"strings"
	"time"
)

type battleSQLRepository struct {
	db *sql.DB
}

func (repo *battleSQLRepository) Count(ctx context.Context) int {
	var total int

	q := "SELECT COUNT(*) AS total FROM battles"
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
	q += "'name', m.name, 'avatar', m.avatar)) FROM battle_players as bp join monsters as m on bp.monster_id = "
	q += "m.id where bp.battle_id = b.id) AS CHAR) as battle_players FROM battles as b "
	if len(args) > 0 {
		for _, arg := range args {
			if strings.Contains(arg, "WHERE") {
				q += fmt.Sprintf("%s ", arg)
				q += "ORDER BY b.id DESC "
			} else {
				if !strings.Contains(q, "ORDER BY") {
					q += "ORDER BY b.id DESC "
				}
				q += fmt.Sprintf("%s ", arg)
			}
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
	newBattle, err := tx.ExecContext(ctx, qb, param.StartedAt, param.EndedAt)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	// doesn't need to validate err, because
	// the database table has an auto-incrementing primary key
	battleID, _ := newBattle.LastInsertId()

	ql := "INSERT INTO battle_logs (battle_id, description, created_at) VALUES"
	now := time.Now().UnixMicro()
	for i, log := range param.Logs {
		ql += fmt.Sprintf(" (%d, '%s', %d)", battleID, log.Description, now)
		if i != (len(param.Logs) - 1) {
			ql += ","
		}
	}
	_, err = tx.ExecContext(ctx, ql)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	qp := "INSERT INTO battle_players (battle_id, monster_id, eliminated_at, annulled_at, rank, point) VALUES"
	for i, player := range param.Players {
		qp += fmt.Sprintf(" (%d, %d, %d, %d, %d, %d)",
			battleID, player.MonsterID, player.EliminatedAt,
			player.AnnulledAt, player.Rank, player.Point)
		if i != (len(param.Players) - 1) {
			qp += ","
		}
	}
	_, err = tx.ExecContext(ctx, qp)
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

func NewBattleSQLRepository() domain.IBattleRepository {
	return &battleSQLRepository{db: configs.DbPool}
}
