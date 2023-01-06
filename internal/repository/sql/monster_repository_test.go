package sql_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aasumitro/pokewar/configs"
	"github.com/aasumitro/pokewar/domain"
	repoSql "github.com/aasumitro/pokewar/internal/repository/sql"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
	"time"
)

type monsterSQLRepositoryTestSuite struct {
	suite.Suite
	mock     sqlmock.Sqlmock
	repo     domain.IMonsterRepository
	monster  *domain.Monster
	monsters []*domain.Monster
}

func (suite *monsterSQLRepositoryTestSuite) SetupSuite() {
	var err error

	configs.DbPool, suite.mock, err = sqlmock.New(
		sqlmock.QueryMatcherOption(
			sqlmock.QueryMatcherRegexp))

	require.NoError(suite.T(), err)

	suite.repo = repoSql.NewMonsterSQLRepository()

	suite.monster = &domain.Monster{OriginID: 1,
		Name:    "test",
		BaseExp: 1,
		Height:  1,
		Weight:  1,
		Avatar:  "test.png",
		Types:   []string{"grass"},
		Stats:   []domain.Stat{{Name: "asd", BaseStat: 1}},
		Skills:  []*domain.Skill{{PP: 1, Name: "as"}},
	}

	suite.monsters = []*domain.Monster{
		suite.monster,
		{
			Name:    "test2",
			BaseExp: 2,
			Height:  2,
			Weight:  2,
			Avatar:  "test2.png",
			Types:   []string{"grass2"},
			Stats:   []domain.Stat{{Name: "asd2", BaseStat: 1}},
			Skills:  []*domain.Skill{{PP: 1, Name: "as2"}},
		},
	}
}

// =========== ALL
func (suite *monsterSQLRepositoryTestSuite) TestRepository_All_ExpectedReturnData() {
	data := suite.mock.
		NewRows([]string{"id", "origin_id", "name", "base_exp", "height", "weight", "avatar", "types", "stats", "skills"}).
		AddRow(1, 1, "test1", 1, 1, 1, "lorem.png", "[\"grass\",\"poison\"]", "[{\"base_stat\":45,\"name\":\"hp\"}]", "[{\"pp\":15,\"name\":\"echoed-voice\"}]").
		AddRow(2, 2, "test2", 2, 2, 2, "lorem.png", "[\"grass\",\"poison\"]", "[{\"base_stat\":45,\"name\":\"hp\"}]", "[{\"pp\":15,\"name\":\"echoed-voice\"}]")
	q := "SELECT id, origin_id, name, base_exp, height, weight, avatar, types, stats, skills FROM monsters LIMIT 1"
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).WillReturnRows(data)
	res, err := suite.repo.All(context.TODO(), "LIMIT 1")
	require.Nil(suite.T(), err)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), res)
	require.Nil(suite.T(), suite.mock.ExpectationsWereMet())
}
func (suite *monsterSQLRepositoryTestSuite) TestRepository_All_ExpectedReturnErrorFromQuery() {
	q := "SELECT id, origin_id, name, base_exp, height, weight, avatar, types, stats, skills FROM monsters"
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).WillReturnError(errors.New(""))
	res, err := suite.repo.All(context.TODO())
	require.NotNil(suite.T(), err)
	require.Nil(suite.T(), res)
	require.Nil(suite.T(), suite.mock.ExpectationsWereMet())
}
func (suite *monsterSQLRepositoryTestSuite) TestRepository_All_ExpectedReturnErrorFromScan() {
	data := suite.mock.
		NewRows([]string{"id", "origin_id", "name", "base_exp", "height", "weight", "avatar", "types", "stats", "skills"}).
		AddRow(1, 1, "test1", 1, 1, 1, "lorem.png", "[\"grass\",\"poison\"]", "[{\"base_stat\":45,\"name\":\"hp\"}]", "[{\"pp\":15,\"name\":\"echoed-voice\"}]").
		AddRow(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	q := "SELECT id, origin_id, name, base_exp, height, weight, avatar, types, stats, skills FROM monsters"
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).WillReturnRows(data)
	res, err := suite.repo.All(context.TODO())
	require.Nil(suite.T(), res)
	require.NotNil(suite.T(), err)
	require.Nil(suite.T(), suite.mock.ExpectationsWereMet())
}

// =========== UPDATE
func (suite *monsterSQLRepositoryTestSuite) TestRepository_Update_ExpectedSuccess() {
	data := suite.mock.
		NewRows([]string{"id", "origin_id", "name", "base_exp", "height", "weight", "avatar", "types", "stats", "skills"}).
		AddRow(1, 1, "test", 1, 1, 1, "test.png", "[\"grass\"]", "[{\"base_stat\":1,\"name\":\"asd\"}]", "[{\"pp\":1,\"name\":\"as\"}]")
	q := "UPDATE monsters SET name = ?, base_exp = ?, height = ?, weight = ?, avatar = ?, "
	q += "types = ?, stats = ?, skills = ?, updated_at = ? WHERE origin_id = ? RETURNING "
	q += "id, origin_id, name, base_exp, height, weight, avatar, types, stats, skills"
	expectedQuery := regexp.QuoteMeta(q)
	types, _ := json.Marshal(suite.monster.Types)
	stats, _ := json.Marshal(suite.monster.Stats)
	skills, _ := json.Marshal(suite.monster.Skills)
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs(suite.monster.Name, suite.monster.BaseExp,
			suite.monster.Height, suite.monster.Weight, suite.monster.Avatar,
			types, stats, skills, time.Now().Unix(), suite.monster.OriginID).
		WillReturnError(nil).WillReturnRows(data)
	err := suite.repo.Update(context.TODO(), suite.monster)
	require.Nil(suite.T(), err)
	require.Nil(suite.T(), suite.mock.ExpectationsWereMet())
}
func (suite *monsterSQLRepositoryTestSuite) TestRepository_Update_ExpectedError() {
	q := "UPDATE monsters SET name = ?, base_exp = ?, height = ?, weight = ?, avatar = ?, "
	q += "types = ?, stats = ?, skills = ?, updated_at = ? WHERE origin_id = ? RETURNING "
	q += "id, origin_id, name, base_exp, height, weight, avatar, types, stats, skills"
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs().
		WillReturnError(errors.New(""))
	err := suite.repo.Update(context.TODO(), suite.monster)
	require.NotNil(suite.T(), err)
	require.Nil(suite.T(), suite.mock.ExpectationsWereMet())
}

// =========== COUNT
func (suite *monsterSQLRepositoryTestSuite) TestRepository_Count_ExpectedReturnData() {
	data := suite.mock.
		NewRows([]string{"total"}).
		AddRow(50)
	q := "SELECT COUNT(*) AS total FROM monsters"
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).WillReturnRows(data)
	res := suite.repo.Count(context.TODO())
	require.NotNil(suite.T(), res)
	require.EqualValues(suite.T(), res, 50)
	require.Nil(suite.T(), suite.mock.ExpectationsWereMet())
}
func (suite *monsterSQLRepositoryTestSuite) TestRepository_Count_ExpectedReturnErrorFromQuery() {
	q := "SELECT COUNT(*) AS total FROM monsters"
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).WillReturnError(errors.New(""))
	res := suite.repo.Count(context.TODO())
	require.EqualValues(suite.T(), res, 0)
	require.Nil(suite.T(), suite.mock.ExpectationsWereMet())
}

// =========== CREATE
func (suite *monsterSQLRepositoryTestSuite) TestRepository_Create_ExpectedSuccess() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`INSERT INTO monsters`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()
	err := suite.repo.Create(context.TODO(), suite.monsters)
	require.Nil(suite.T(), err)
	require.Nil(suite.T(), suite.mock.ExpectationsWereMet())
}
func (suite *monsterSQLRepositoryTestSuite) TestRepository_Create_ShouldErrorTxBegin() {
	suite.mock.ExpectBegin().WillReturnError(errors.New("UNEXPECTED"))
	err := suite.repo.Create(context.TODO(), suite.monsters)
	require.NotNil(suite.T(), err)
	require.Equal(suite.T(), err.Error(), "UNEXPECTED")
	require.Nil(suite.T(), suite.mock.ExpectationsWereMet())
}
func (suite *monsterSQLRepositoryTestSuite) TestRepository_Create_ShouldErrorTxExec() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`INSERT INTO monsters`).
		WillReturnError(errors.New("UNEXPECTED"))
	suite.mock.ExpectRollback()
	err := suite.repo.Create(context.TODO(), suite.monsters)
	require.NotNil(suite.T(), err)
	require.Equal(suite.T(), err.Error(), "UNEXPECTED")
	require.Nil(suite.T(), suite.mock.ExpectationsWereMet())
}
func (suite *monsterSQLRepositoryTestSuite) TestRepository_Create_ShouldErrorCommit() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`INSERT INTO monsters`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit().WillReturnError(errors.New("UNEXPECTED"))
	err := suite.repo.Create(context.TODO(), suite.monsters)
	require.NotNil(suite.T(), err)
	require.Equal(suite.T(), err.Error(), "UNEXPECTED")
	require.Nil(suite.T(), suite.mock.ExpectationsWereMet())
}

func TestMonsterRepository(t *testing.T) {
	suite.Run(t, new(monsterSQLRepositoryTestSuite))
}
