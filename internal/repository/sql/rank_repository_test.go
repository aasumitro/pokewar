package sql_test

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aasumitro/pokewar/domain"
	repoSql "github.com/aasumitro/pokewar/internal/repository/sql"
	"github.com/aasumitro/pokewar/pkg/appconfigs"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
)

type rankSQLRepositoryTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	repo domain.IRankRepository
}

func (suite *rankSQLRepositoryTestSuite) SetupSuite() {
	var err error

	appconfigs.DbPool, suite.mock, err = sqlmock.New(
		sqlmock.QueryMatcherOption(
			sqlmock.QueryMatcherRegexp))

	require.NoError(suite.T(), err)

	suite.repo = repoSql.NewRankSQLRepository()
}

// =========== ALL
func (suite *rankSQLRepositoryTestSuite) TestRepository_All_ExpectedReturnData() {
	data := suite.mock.
		NewRows([]string{"id", "origin_id", "name", "avatar", "types", "total_battles", "points", "win_battles", "lost_battles"}).
		AddRow(1, 1, "test1", "test.png", "['test']", 1, 1, 1, 1).
		AddRow(2, 2, "test2", "test.png", "['test']", 2, 2, 2, 2)
	q := "SELECT monsters.id as id, monsters.origin_id as origin_id, monsters.name as name, "
	q += "monsters.avatar as avatar, monsters.types as types, "
	q += "count(p.monster_id) as total_battles, sum(IFNULL(p.point, 0)) as points, "
	q += "(SELECT count(w.rank) FROM battle_players as w where rank = 1  "
	q += "AND monster_id = monsters.id AND annulled_at = 0) as win_battles, "
	q += "(SELECT count(l.rank) FROM battle_players as l where rank > 1 "
	q += "AND monster_id = monsters.id AND annulled_at = 0) as lose_battles "
	q += "FROM monsters LEFT JOIN battle_players as p  "
	q += "ON monsters.id = p.monster_id GROUP BY monsters.id ORDER BY points DESC "
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).WillReturnRows(data)
	res, err := suite.repo.All(context.TODO(), "LIMIT 1")
	require.Nil(suite.T(), err)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), res)
}
func (suite *rankSQLRepositoryTestSuite) TestRepository_All_ExpectedReturnErrorFromQuery() {
	q := "SELECT monsters.id as id, monsters.origin_id as origin_id, monsters.name as name, "
	q += "monsters.avatar as avatar, monsters.types as types, "
	q += "count(p.monster_id) as total_battles, sum(IFNULL(p.point, 0)) as points, "
	q += "(SELECT count(w.rank) FROM battle_players as w where rank = 1  "
	q += "AND monster_id = monsters.id AND annulled_at = 0) as win_battles, "
	q += "(SELECT count(l.rank) FROM battle_players as l where rank > 1 "
	q += "AND monster_id = monsters.id AND annulled_at = 0) as lose_battles "
	q += "FROM monsters LEFT JOIN battle_players as p  "
	q += "ON monsters.id = p.monster_id GROUP BY monsters.id ORDER BY points DESC "
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).WillReturnError(errors.New(""))
	res, err := suite.repo.All(context.TODO())
	require.NotNil(suite.T(), err)
	require.Nil(suite.T(), res)
}
func (suite *rankSQLRepositoryTestSuite) TestRepository_All_ExpectedReturnErrorFromScan() {
	data := suite.mock.
		NewRows([]string{"id", "origin_id", "name", "avatar", "types", "total_battles", "points", "win_battles", "lost_battles"}).
		AddRow(1, 1, "test1", "test.png", "['test']", 1, 1, 1, 1).
		AddRow(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	q := "SELECT monsters.id as id, monsters.origin_id as origin_id, monsters.name as name, "
	q += "monsters.avatar as avatar, monsters.types as types, "
	q += "count(p.monster_id) as total_battles, sum(IFNULL(p.point, 0)) as points, "
	q += "(SELECT count(w.rank) FROM battle_players as w where rank = 1  "
	q += "AND monster_id = monsters.id AND annulled_at = 0) as win_battles, "
	q += "(SELECT count(l.rank) FROM battle_players as l where rank > 1 "
	q += "AND monster_id = monsters.id AND annulled_at = 0) as lose_battles "
	q += "FROM monsters LEFT JOIN battle_players as p  "
	q += "ON monsters.id = p.monster_id GROUP BY monsters.id ORDER BY points DESC "
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).WillReturnRows(data)
	res, err := suite.repo.All(context.TODO())
	require.Nil(suite.T(), res)
	require.NotNil(suite.T(), err)
}

func TestRankRepository(t *testing.T) {
	suite.Run(t, new(rankSQLRepositoryTestSuite))
}
