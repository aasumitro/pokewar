package sql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aasumitro/pokewar/domain"
	"github.com/aasumitro/pokewar/pkg/configs"
)

type rankSQLRepository struct {
	db *sql.DB
}

func (repo rankSQLRepository) All(ctx context.Context) (data []*domain.Rank, err error) {
	q := "SELECT monsters.id as id, monsters.origin_id as origin_id, monsters.name as name,  "
	q += "monsters.avatar as avatar, monsters.types as types, count(players.monster_id) as total_battles "
	q += "(SELECT count(w.rank) FROM players as w where rank = 1) as win_battles "
	q += "(SELECT count(l.rank) FROM players as l where rank > 1) as lose_battles "
	q += "FROM monsters LEFT OUTER JOIN players "
	q += "ON monsters.id = players.monster_id "
	q += "GROUP BY monsters.id"
	fmt.Println(data)
	//TODO implement me
	panic("implement me")
}

func NewRankSQLRepository() domain.IRankRepository {
	return rankSQLRepository{db: configs.DbPool}
}
