package sql

import (
	"context"
	"database/sql"
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
	//TODO implement me
	panic("implement me")
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
