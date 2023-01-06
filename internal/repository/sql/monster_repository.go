package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aasumitro/pokewar/configs"
	"github.com/aasumitro/pokewar/domain"
	"time"
)

type monsterSQLRepository struct {
	db *sql.DB
}

func (repo *monsterSQLRepository) All(ctx context.Context, args ...string) (data []*domain.Monster, err error) {
	q := "SELECT id, origin_id, name, base_exp, height, weight, avatar, types, stats, skills FROM monsters"
	if len(args) > 0 {
		for _, arg := range args {
			q += fmt.Sprintf(" %s", arg)
		}
	}

	rows, err := repo.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) { _ = rows.Close() }(rows)

	for rows.Next() {
		var monster domain.MonsterEntity
		var types []string
		var stats []domain.Stat
		var skills []*domain.Skill

		if err := rows.Scan(
			&monster.ID, &monster.OriginID, &monster.Name,
			&monster.BaseExp, &monster.Height, &monster.Weight,
			&monster.Avatar, &monster.Types, &monster.Stats, &monster.Skills,
		); err != nil {
			return nil, err
		}

		_ = json.Unmarshal([]byte(monster.Types), &types)
		_ = json.Unmarshal([]byte(monster.Stats), &stats)
		_ = json.Unmarshal([]byte(monster.Skills), &skills)

		data = append(data, &domain.Monster{
			ID:       monster.ID,
			OriginID: monster.OriginID,
			Name:     monster.Name,
			BaseExp:  monster.BaseExp,
			Height:   monster.Height,
			Weight:   monster.Weight,
			Avatar:   monster.Avatar,
			Types:    types,
			Stats:    stats,
			Skills:   skills,
		})
	}

	return data, nil
}

func (repo *monsterSQLRepository) Create(ctx context.Context, params []*domain.Monster) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	query := "INSERT INTO monsters (origin_id, name, base_exp, height, weight, avatar, types, stats, skills, created_at) VALUES"
	for i, monster := range params {
		types, _ := json.Marshal(monster.Types)
		stats, _ := json.Marshal(monster.Stats)
		skills, _ := json.Marshal(monster.Skills)
		now := time.Now().UnixMicro()
		query += fmt.Sprintf("(%d, '%s', %d, %d, %d, '%s', '%s', '%s', '%s', %d)",
			monster.OriginID, monster.Name, monster.BaseExp,
			monster.Height, monster.Weight, monster.Avatar,
			types, stats, skills, now)
		if i != (len(params) - 1) {
			query += ","
		}
	}

	_, err = tx.ExecContext(ctx, query)
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

func (repo *monsterSQLRepository) Update(ctx context.Context, param *domain.Monster) error {
	var monster domain.MonsterEntity

	types, _ := json.Marshal(param.Types)
	stats, _ := json.Marshal(param.Stats)
	skills, _ := json.Marshal(param.Skills)

	q := "UPDATE monsters SET name = ?, base_exp = ?, height = ?, weight = ?, avatar = ?, "
	q += "types = ?, stats = ?, skills = ?, updated_at = ? WHERE origin_id = ? RETURNING "
	q += "id, origin_id, name, base_exp, height, weight, avatar, types, stats, skills"
	if err := repo.db.QueryRowContext(ctx, q,
		param.Name, param.BaseExp, param.Height, param.Weight, param.Avatar,
		types, stats, skills, time.Now().Unix(), param.OriginID,
	).Scan(
		&monster.ID, &monster.OriginID, &monster.Name,
		&monster.BaseExp, &monster.Height, &monster.Weight,
		&monster.Avatar, &monster.Types, &monster.Stats, &monster.Stats,
	); err != nil {
		return err
	}

	return nil
}

func (repo *monsterSQLRepository) Count(ctx context.Context) int {
	var total int

	q := "SELECT COUNT(*) AS total FROM monsters"
	if err := repo.db.QueryRowContext(ctx, q).Scan(&total); err != nil {
		total = 0
	}

	return total
}

func NewMonsterSQLRepository() domain.IMonsterRepository {
	return &monsterSQLRepository{db: configs.DbPool}
}
