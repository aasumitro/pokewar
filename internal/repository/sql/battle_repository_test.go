package sql_test

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aasumitro/pokewar/configs"
	"github.com/aasumitro/pokewar/domain"
	repoSql "github.com/aasumitro/pokewar/internal/repository/sql"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
)

type battleSQLRepositoryTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	repo domain.IBattleRepository
}

func (suite *battleSQLRepositoryTestSuite) SetupSuite() {
	var err error

	configs.DbPool, suite.mock, err = sqlmock.New(
		sqlmock.QueryMatcherOption(
			sqlmock.QueryMatcherRegexp))

	require.NoError(suite.T(), err)

	suite.repo = repoSql.NewBattleSQLRepository()
}

// =========== COUNT
func (suite *battleSQLRepositoryTestSuite) TestRepository_Count_ExpectedReturnData() {
	data := suite.mock.
		NewRows([]string{"total"}).
		AddRow(50)
	q := "SELECT COUNT(*) AS total FROM battles"
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).WillReturnRows(data)
	res := suite.repo.Count(context.TODO())
	require.NotNil(suite.T(), res)
}
func (suite *battleSQLRepositoryTestSuite) TestRepository_Count_ExpectedReturnErrorFromQuery() {
	q := "SELECT COUNT(*) AS total FROM battles"
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
	q += "m.id where bp.battle_id = b.id) AS CHAR) as battle_players FROM battles as b ORDER BY b.id DESC LIMIT 1 "
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).WillReturnRows(data)
	res, err := suite.repo.All(context.TODO(), "LIMIT 1")
	require.Nil(suite.T(), err)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), res)
}
func (suite *battleSQLRepositoryTestSuite) TestRepository_All_ExpectedReturnDataWithWhere() {
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
	q += "m.id where bp.battle_id = b.id) AS CHAR) as battle_players FROM battles as b "
	q += "WHERE b.created_at BETWEEN 1 AND 2 ORDER BY b.id DESC "
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).WillReturnRows(data)
	res, err := suite.repo.All(context.TODO(), "WHERE b.created_at BETWEEN 1 AND 2")
	require.Nil(suite.T(), err)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), res)
}
func (suite *battleSQLRepositoryTestSuite) TestRepository_All_ExpectedReturnErrorFromScan() {
	data := suite.mock.
		NewRows([]string{"id", "started_at", "ended_at", "players", "logs"}).
		AddRow(nil, nil, nil, nil, nil).
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
func (suite *battleSQLRepositoryTestSuite) TestRepository_All_ExpectedReturnErrorFromQuery() {
	q := "SELECT b.id as id, b.started_at as started_at, b.ended_at as ended_at, "
	q += "CAST((SELECT json_group_array(json_object('id', bl.id, 'battle_id', bl.battle_id, 'description', bl.description, "
	q += "'created_at', bl.created_at)) FROM battle_logs as bl where bl.battle_id = b.id) AS CHAR) as battle_logs, "
	q += "CAST((SELECT json_group_array(json_object('id', bp.id, 'battle_id', bp.battle_id, 'monster_id', bp.monster_id, "
	q += "'eliminated_at', bp.eliminated_at, 'annulled_at', bp.annulled_at, 'rank', bp.rank, 'point', bp.point, "
	q += "'name', m.name, 'avatar', m.avatar)) FROM battle_players as bp join monsters as m on bp.monster_id = "
	q += "m.id where bp.battle_id = b.id) AS CHAR) as battle_players FROM battles as b ORDER BY b.id DESC "
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).WillReturnError(errors.New(""))
	res, err := suite.repo.All(context.TODO())
	require.NotNil(suite.T(), err)
	require.Nil(suite.T(), res)
}

// =========== Create
type battleSQLRepositoryTestSuiteForCreate struct {
	suite.Suite
	mock sqlmock.Sqlmock
	repo domain.IBattleRepository
}

func (suite *battleSQLRepositoryTestSuiteForCreate) SetupSuite() {
	var err error

	configs.DbPool, suite.mock, err = sqlmock.New(
		sqlmock.QueryMatcherOption(
			sqlmock.QueryMatcherRegexp))

	require.NoError(suite.T(), err)

	suite.repo = repoSql.NewBattleSQLRepository()
}

func (suite *battleSQLRepositoryTestSuiteForCreate) TestRepository_Create_ShouldSuccess() {
	battle := &domain.Battle{
		StartedAt: 1234567890,
		EndedAt:   1234567891,
		Logs:      []domain.Log{{Description: "Test log"}, {Description: "Test log 2"}},
		Players:   []domain.Player{{MonsterID: 1, EliminatedAt: 1234567891, Rank: 1, Point: 5}, {MonsterID: 2, EliminatedAt: 1234567891, Rank: 2, Point: 4}},
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`INSERT INTO battles`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectExec(`INSERT INTO battle_logs`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectExec(`INSERT INTO battle_players`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()
	err := suite.repo.Create(context.TODO(), battle)
	require.Nil(suite.T(), err)
}
func (suite *battleSQLRepositoryTestSuiteForCreate) TestRepository_Create_ShouldErrorTxBegin() {
	battle := &domain.Battle{
		StartedAt: 1234567890,
		EndedAt:   1234567891,
		Logs:      []domain.Log{{Description: "Test log"}},
		Players:   []domain.Player{{MonsterID: 1, EliminatedAt: 1234567891, Rank: 1, Point: 5}},
	}

	suite.mock.ExpectBegin().WillReturnError(errors.New("UNEXPECTED"))
	err := suite.repo.Create(context.TODO(), battle)
	require.NotNil(suite.T(), err)
	require.Equal(suite.T(), err.Error(), "UNEXPECTED")
}
func (suite *battleSQLRepositoryTestSuiteForCreate) TestRepository_Create_ShouldErrorInsertBattle() {
	battle := &domain.Battle{
		StartedAt: 1234567890,
		EndedAt:   1234567891,
		Logs:      []domain.Log{{Description: "Test log"}},
		Players:   []domain.Player{{MonsterID: 1, EliminatedAt: 1234567891, Rank: 1, Point: 5}},
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`INSERT INTO battles`).
		WillReturnError(errors.New("UNEXPECTED"))
	suite.mock.ExpectRollback()
	err := suite.repo.Create(context.TODO(), battle)
	require.NotNil(suite.T(), err)
	require.Equal(suite.T(), err.Error(), "UNEXPECTED")
}
func (suite *battleSQLRepositoryTestSuiteForCreate) TestRepository_Create_ShouldErrorInsertLogs() {
	battle := &domain.Battle{
		StartedAt: 1234567890,
		EndedAt:   1234567891,
		Logs:      []domain.Log{{Description: "Test log"}},
		Players:   []domain.Player{{MonsterID: 1, EliminatedAt: 1234567891, Rank: 1, Point: 5}},
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`INSERT INTO battles`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectExec(`INSERT INTO battle_logs`).
		WillReturnError(errors.New("UNEXPECTED"))
	suite.mock.ExpectRollback()
	err := suite.repo.Create(context.TODO(), battle)
	require.NotNil(suite.T(), err)
	require.Equal(suite.T(), err.Error(), "UNEXPECTED")
}
func (suite *battleSQLRepositoryTestSuiteForCreate) TestRepository_Create_ShouldErrorInsertPlayers() {
	battle := &domain.Battle{
		StartedAt: 1234567890,
		EndedAt:   1234567891,
		Logs:      []domain.Log{{Description: "Test log"}},
		Players:   []domain.Player{{MonsterID: 1, EliminatedAt: 1234567891, Rank: 1, Point: 5}},
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`INSERT INTO battles`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectExec(`INSERT INTO battle_logs`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectExec(`INSERT INTO battle_players`).
		WillReturnError(errors.New("UNEXPECTED"))
	suite.mock.ExpectRollback()
	err := suite.repo.Create(context.TODO(), battle)
	require.NotNil(suite.T(), err)
	require.Equal(suite.T(), err.Error(), "UNEXPECTED")
}
func (suite *battleSQLRepositoryTestSuiteForCreate) TestRepository_Create_ShouldErrorCommit() {
	battle := &domain.Battle{
		StartedAt: 1234567890,
		EndedAt:   1234567891,
		Logs:      []domain.Log{{Description: "Test log"}},
		Players:   []domain.Player{{MonsterID: 1, EliminatedAt: 1234567891, Rank: 1, Point: 5}},
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`INSERT INTO battles`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectExec(`INSERT INTO battle_logs`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectExec(`INSERT INTO battle_players`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit().WillReturnError(errors.New("UNEXPECTED"))
	err := suite.repo.Create(context.TODO(), battle)
	require.NotNil(suite.T(), err)
	require.Equal(suite.T(), err.Error(), "UNEXPECTED")
}

func TestBattleSQLRepository(t *testing.T) {
	suite.Run(t, new(battleSQLRepositoryTestSuite))
	suite.Run(t, new(battleSQLRepositoryTestSuiteForCreate))
}
