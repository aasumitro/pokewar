package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/aasumitro/pokewar/domain"
	"github.com/aasumitro/pokewar/pkg/configs"
	"time"
)

type monsterSQLRepository struct {
	db *sql.DB
}

func (repo monsterSQLRepository) All(ctx context.Context) (data []*domain.Monster, err error) {
	q := "SELECT id, origin_id, name, base_exp, height, weight, avatar, types, stats, skills FROM monsters"
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

func (repo monsterSQLRepository) AllWhereIn(ctx context.Context, id []int) (data []*domain.Monster, err error) {
	q := "SELECT id, origin_id, name, base_exp, height, weight, avatar, types, "
	q += "stats, skills FROM monsters WHERE origin_id IN (?,?,?,?,?) LIMIT 5"
	rows, err := repo.db.QueryContext(ctx, q, id[0], id[1], id[2], id[3], id[4])
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

func (repo monsterSQLRepository) Create(ctx context.Context, param *domain.Monster) error {
	var monster domain.MonsterEntity

	types, _ := json.Marshal(param.Types)
	stats, _ := json.Marshal(param.Stats)
	skills, _ := json.Marshal(param.Skills)

	q := "INSERT INTO monsters (origin_id, name, base_exp, height, weight, avatar, types, stats, skills, created_at) "
	q += "VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING "
	q += "id, origin_id, name, base_exp, height, weight, avatar, types, stats, skills"
	if err := repo.db.QueryRowContext(ctx, q,
		param.OriginID, param.Name, param.BaseExp,
		param.Height, param.Weight, param.Avatar,
		types, stats, skills, time.Now().Unix(),
	).Scan(
		&monster.ID, &monster.OriginID, &monster.Name,
		&monster.BaseExp, &monster.Height, &monster.Weight,
		&monster.Avatar, &monster.Types, &monster.Stats, &monster.Stats,
	); err != nil {
		return err
	}

	return nil
}

func (repo monsterSQLRepository) Update(ctx context.Context, param *domain.Monster) error {
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

func NewMonsterSQlRepository() domain.IMonsterRepository {
	return &monsterSQLRepository{db: configs.DbPool}
}
