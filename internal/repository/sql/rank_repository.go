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

// All
// TODO ADD:
// PAGINATION MAYBE?
// LIMIT? ORDER? etc.
//
// ALTERNATE QUERY -> need to use sql.NullType && need to validate data (return)
//
//	q := `
//	SELECT m.id as id,
//		   m.origin_id as origin_id,
//		   m.name as name,
//		   m.avatar as avatar,
//		   m.types as types,
//		   p.total_battles,
//		   p.points,
//		   p.win_battles,
//		   p.lose_battles
//	FROM monsters as m
//	LEFT JOIN
//	(
//		SELECT monster_id,
//			   COUNT(*) as total_battles,
//			   SUM(IFNULL(point, 0)) as points,
//			   (
//			   		SELECT COUNT(IFNULL(*, 0)) FROM battle_players as w
//			        WHERE rank = 1 AND monster_id = bp.monster_id
//			        AND annulled_at IS NULL
//			   ) as win_battles,
//			   (
//			   		SELECT COUNT(IFNULL(*, 0)) FROM battle_players as l
//					WHERE rank > 1 AND monster_id = bp.monster_id
//				  	AND annulled_at IS NULL
//				) as lose_battles
//		FROM battle_players as bp
//		GROUP BY monster_id
//	) as p ON m.id = p.monster_id
//	`
func (repo *rankSQLRepository) All(ctx context.Context, args ...string) (data []*domain.Rank, err error) {
	// TODO: Optimize this query
	q := "SELECT monsters.id as id, monsters.origin_id as origin_id, monsters.name as name, "
	q += "monsters.avatar as avatar, monsters.types as types, "
	q += "count(p.monster_id) as total_battles, sum(IFNULL(p.point, 0)) as points, "
	q += "(SELECT count(w.rank) FROM battle_players as w where rank = 1  "
	q += "AND monster_id = monsters.id AND annulled_at = null) as win_battles, "
	q += "(SELECT count(l.rank) FROM battle_players as l where rank > 1 "
	q += "AND monster_id = monsters.id AND annulled_at = null) as lose_battles "
	q += "FROM monsters LEFT JOIN battle_players as p  "
	q += "ON monsters.id = p.monster_id GROUP BY monsters.id "
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
