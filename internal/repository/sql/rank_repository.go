package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aasumitro/pokewar/domain"
	"github.com/aasumitro/pokewar/pkg/appconfigs"
)

type rankSQLRepository struct {
	db *sql.DB
}

func (repo *rankSQLRepository) All(ctx context.Context, args ...string) (data []*domain.Rank, err error) {
	// TODO: Optimize this query
	q := "SELECT monsters.id as id, monsters.origin_id as origin_id, monsters.name as name, "
	q += "monsters.avatar as avatar, monsters.types as types, "
	q += "count(p.monster_id) as total_battles, sum(IFNULL(p.point, 0)) as points, "
	q += "(SELECT count(w.rank) FROM battle_players as w where rank = 1  "
	q += "AND monster_id = monsters.id AND annulled_at = 0) as win_battles, "
	q += "(SELECT count(l.rank) FROM battle_players as l where rank > 1 "
	q += "AND monster_id = monsters.id AND annulled_at = 0) as lose_battles "
	q += "FROM monsters LEFT JOIN battle_players as p  "
	q += "ON monsters.id = p.monster_id GROUP BY monsters.id ORDER BY points DESC "
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
		var rank domain.Rank

		if err := rows.Scan(
			&rank.ID, &rank.OriginID, &rank.Name,
			&rank.Avatar, &rank.T, &rank.TotalBattles,
			&rank.Points, &rank.WinBattles, &rank.LoseBattle,
		); err != nil {
			return nil, err
		}

		var types []string
		_ = json.Unmarshal([]byte(rank.T), &types)
		rank.Types = types
		rank.T = ""
		data = append(data, &rank)
	}

	return data, nil
}

func NewRankSQLRepository() domain.IRankRepository {
	return &rankSQLRepository{db: appconfigs.DbPool}
}
