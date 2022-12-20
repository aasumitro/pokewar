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
	"time"
)

type battleSQLRepositoryTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	repo domain.IBattleRepository
}

func (suite *battleSQLRepositoryTestSuite) SetupSuite() {
	var err error

	appconfigs.DbPool, suite.mock, err = sqlmock.New(
		sqlmock.QueryMatcherOption(
			sqlmock.QueryMatcherRegexp))

	require.NoError(suite.T(), err)

	suite.repo = repoSql.NewBattleSQLRepository()
}

// =========== COUNT
func (suite *battleSQLRepositoryTestSuite) TestRepository_Count_ExpectedReturnData() {
	data := suite.mock.
		NewRows([]string{"count"}).
		AddRow(50)
	q := "SELECT COUNT(*) FROM battles"
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).WillReturnRows(data)
	res := suite.repo.Count(context.TODO())
	require.NotNil(suite.T(), res)
	require.EqualValues(suite.T(), res, 50)
}
func (suite *battleSQLRepositoryTestSuite) TestRepository_Count_ExpectedReturnErrorFromQuery() {
	q := "SELECT COUNT(*) FROM battles"
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).WillReturnError(errors.New(""))
	res := suite.repo.Count(context.TODO())
	require.EqualValues(suite.T(), res, 0)
}

// =========== ALL
func (suite *battleSQLRepositoryTestSuite) TestRepository_All_ExpectedReturnData() {
	data := suite.mock.
		NewRows([]string{"id", "started_at", "ended_at", "players", "logs"}).
		AddRow(1, 1, 1, "[{\"id\":1,\"description\":\"venomoth attack ekans\"}]", "[{\"id\":1,\"monster_id\":1,\"eliminated_at\":null,\"annulled_at\":null,\"rank\":1,\"point\":5,\"name\":\"venomoth\"}]").
		AddRow(2, 2, 2, "[{\"id\":1,\"description\":\"venomoth attack ekans\"}]", "[{\"id\":1,\"monster_id\":1,\"eliminated_at\":null,\"annulled_at\":null,\"rank\":1,\"point\":5,\"name\":\"venomoth\"}]")
	q := "SELECT b.id as id, b.started_at as started_at, b.ended_at as ended_at, "
	q += "CAST((SELECT json_group_array(json_object('id', bl.id, 'battle_id', bl.battle_id, 'description', bl.description, "
	q += "'created_at', bl.created_at)) FROM battle_logs as bl where bl.battle_id = b.id) AS CHAR) as battle_logs, "
	q += "CAST((SELECT json_group_array(json_object('id', bp.id, 'battle_id', bp.battle_id, 'monster_id', bp.monster_id, "
	q += "'eliminated_at', bp.eliminated_at, 'annulled_at', bp.annulled_at, 'rank', bp.rank, 'point', bp.point, "
	q += "'name', m.name, 'avatar', m.avatar)) FROM battle_players as bp join monsters as m on bp.monster_id = "
	q += "m.id where bp.battle_id = b.id) AS CHAR) as battle_players FROM battles as b LIMIT 1"
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).WillReturnRows(data)
	res, err := suite.repo.All(context.TODO(), "LIMIT 1")
	require.Nil(suite.T(), err)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), res)
}
func (suite *battleSQLRepositoryTestSuite) TestRepository_All_ExpectedReturnErrorFromQuery() {
	q := "SELECT b.id as id, b.started_at as started_at, b.ended_at as ended_at, "
	q += "CAST((SELECT json_group_array(json_object('id', bl.id, 'battle_id', bl.battle_id, 'description', bl.description, "
	q += "'created_at', bl.created_at)) FROM battle_logs as bl where bl.battle_id = b.id) AS CHAR) as battle_logs, "
	q += "CAST((SELECT json_group_array(json_object('id', bp.id, 'battle_id', bp.battle_id, 'monster_id', bp.monster_id, "
	q += "'eliminated_at', bp.eliminated_at, 'annulled_at', bp.annulled_at, 'rank', bp.rank, 'point', bp.point, "
	q += "'name', m.name, 'avatar', m.avatar)) FROM battle_players as bp join monsters as m on bp.monster_id = "
	q += "m.id where bp.battle_id = b.id) AS CHAR) as battle_players FROM battles as b "
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).WillReturnError(errors.New(""))
	res, err := suite.repo.All(context.TODO())
	require.NotNil(suite.T(), err)
	require.Nil(suite.T(), res)
}
func (suite *battleSQLRepositoryTestSuite) TestRepository_All_ExpectedReturnErrorFromScan() {
	data := suite.mock.
		NewRows([]string{"id", "started_at", "ended_at", "players", "logs"}).
		AddRow(1, 1, 1, "[{\"id\":1,\"description\":\"venomoth attack ekans\"}]", "[{\"id\":1,\"monster_id\":1,\"eliminated_at\":null,\"annulled_at\":null,\"rank\":1,\"point\":5,\"name\":\"venomoth\"}]").
		AddRow(nil, nil, nil, nil, nil)
	q := "SELECT b.id as id, b.started_at as started_at, b.ended_at as ended_at, "
	q += "CAST((SELECT json_group_array(json_object('id', bl.id, 'battle_id', bl.battle_id, 'description', bl.description, "
	q += "'created_at', bl.created_at)) FROM battle_logs as bl where bl.battle_id = b.id) AS CHAR) as battle_logs, "
	q += "CAST((SELECT json_group_array(json_object('id', bp.id, 'battle_id', bp.battle_id, 'monster_id', bp.monster_id, "
	q += "'eliminated_at', bp.eliminated_at, 'annulled_at', bp.annulled_at, 'rank', bp.rank, 'point', bp.point, "
	q += "'name', m.name, 'avatar', m.avatar)) FROM battle_players as bp join monsters as m on bp.monster_id = "
	q += "m.id where bp.battle_id = b.id) AS CHAR) as battle_players FROM battles as b "
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).WillReturnRows(data)
	res, err := suite.repo.All(context.TODO())
	require.Nil(suite.T(), res)
	require.NotNil(suite.T(), err)
}

// =========== Create
func (suite *battleSQLRepositoryTestSuite) TestRepository_Create_TestTable() {
	// TODO
	//testCase := struct {
	//	param *domain.Battle
	//	want  error
	//}{
	//	param: &domain.Battle{
	//		StartedAt: 1234567890,
	//		EndedAt:   1234567891,
	//		Logs: []domain.Log{
	//			{Description: "Test log"},
	//		},
	//		Players: []domain.Player{
	//			{MonsterID: 1, EliminatedAt: 1234567891, Rank: 1, Point: 5},
	//		},
	//	},
	//	want: nil,
	//}
}

// =========== UpdatePlayer
func (suite *battleSQLRepositoryTestSuite) TestRepository_UpdatePlayer_ExpectedSuccess() {
	now := time.Now().Unix()
	row := suite.mock.
		NewRows([]string{"annulled_at"}).
		AddRow(now)
	q := "UPDATE battle_players SET annulled_at = ? WHERE id = ?"
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs(now, 1).
		WillReturnRows(row).
		WillReturnError(nil)
	res, err := suite.repo.UpdatePlayer(context.TODO(), 1)
	require.Nil(suite.T(), err)
	require.NotNil(suite.T(), res)
	require.EqualValues(suite.T(), now, res)
}
func (suite *battleSQLRepositoryTestSuite) TestRepository_UpdatePlayer_ExpectedError() {
	now := time.Now().Unix()
	row := suite.mock.
		NewRows([]string{"annulled_at"}).
		AddRow(nil)
	q := "UPDATE battle_players SET annulled_at = ? WHERE id = ?"
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs(now, 1).
		WillReturnRows(row).
		WillReturnError(nil)
	res, err := suite.repo.UpdatePlayer(context.TODO(), 1)
	require.NotNil(suite.T(), err)
	require.EqualValues(suite.T(), 0, res)
}

func TestBattleSQLRepository(t *testing.T) {
	suite.Run(t, new(battleSQLRepositoryTestSuite))
}
