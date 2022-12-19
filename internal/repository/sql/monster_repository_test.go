package sql_test

import (
	"context"
	"encoding/json"
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

type monsterSQLRepositoryTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	repo domain.IMonsterRepository
	data *domain.Monster
}

func (suite *monsterSQLRepositoryTestSuite) SetupSuite() {
	var err error

	appconfigs.DbPool, suite.mock, err = sqlmock.New(
		sqlmock.QueryMatcherOption(
			sqlmock.QueryMatcherRegexp))

	require.NoError(suite.T(), err)

	suite.repo = repoSql.NewMonsterSQlRepository()

	suite.data = &domain.Monster{OriginID: 1,
		Name:    "test",
		BaseExp: 1,
		Height:  1,
		Weight:  1,
		Avatar:  "test.png",
		Types:   []string{"grass"},
		Stats:   []domain.Stat{{Name: "asd", BaseStat: 1}},
		Skills:  []*domain.Skill{{PP: 1, Name: "as"}},
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
}
func (suite *monsterSQLRepositoryTestSuite) TestRepository_All_ExpectedReturnErrorFromQuery() {
	q := "SELECT id, origin_id, name, base_exp, height, weight, avatar, types, stats, skills FROM monsters"
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).WillReturnError(errors.New(""))
	res, err := suite.repo.All(context.TODO())
	require.NotNil(suite.T(), err)
	require.Nil(suite.T(), res)
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
}

// =========== CREATE
func (suite *monsterSQLRepositoryTestSuite) TestRepository_Create_ExpectedSuccess() {
	data := suite.mock.
		NewRows([]string{"id", "origin_id", "name", "base_exp", "height", "weight", "avatar", "types", "stats", "skills"}).
		AddRow(1, 1, "test1", 1, 1, 1, "test.png", "[\"grass\"]", "[{\"base_stat\":1,\"name\":\"asd\"}]", "[{\"pp\":1,\"name\":\"as\"}]")
	q := "INSERT INTO monsters (origin_id, name, base_exp, height, weight, avatar, types, stats, skills, created_at) "
	q += "VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING "
	q += "id, origin_id, name, base_exp, height, weight, avatar, types, stats, skills"
	expectedQuery := regexp.QuoteMeta(q)
	types, _ := json.Marshal(suite.data.Types)
	stats, _ := json.Marshal(suite.data.Stats)
	skills, _ := json.Marshal(suite.data.Skills)
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs(suite.data.OriginID, suite.data.Name, suite.data.BaseExp,
			suite.data.Height, suite.data.Weight, suite.data.Avatar,
			types, stats, skills, time.Now().Unix()).
		WillReturnError(nil).WillReturnRows(data)
	err := suite.repo.Create(context.TODO(), suite.data)
	require.Nil(suite.T(), err)
}
func (suite *monsterSQLRepositoryTestSuite) TestRepository_Create_ExpectedError() {
	q := "INSERT INTO monsters (origin_id, name, base_exp, height, weight, avatar, types, stats, skills, created_at) "
	q += "VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING "
	q += "id, origin_id, name, base_exp, height, weight, avatar, types, stats, skills"
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs().
		WillReturnError(errors.New(""))
	err := suite.repo.Create(context.TODO(), suite.data)
	require.NotNil(suite.T(), err)
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
	types, _ := json.Marshal(suite.data.Types)
	stats, _ := json.Marshal(suite.data.Stats)
	skills, _ := json.Marshal(suite.data.Skills)
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs(suite.data.Name, suite.data.BaseExp,
			suite.data.Height, suite.data.Weight, suite.data.Avatar,
			types, stats, skills, time.Now().Unix(), suite.data.OriginID).
		WillReturnError(nil).WillReturnRows(data)
	err := suite.repo.Update(context.TODO(), suite.data)
	require.Nil(suite.T(), err)
}
func (suite *monsterSQLRepositoryTestSuite) TestRepository_Update_ExpectedError() {
	q := "UPDATE monsters SET name = ?, base_exp = ?, height = ?, weight = ?, avatar = ?, "
	q += "types = ?, stats = ?, skills = ?, updated_at = ? WHERE origin_id = ? RETURNING "
	q += "id, origin_id, name, base_exp, height, weight, avatar, types, stats, skills"
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs().
		WillReturnError(errors.New(""))
	err := suite.repo.Update(context.TODO(), suite.data)
	require.NotNil(suite.T(), err)
}

// =========== COUNT
func (suite *monsterSQLRepositoryTestSuite) TestRepository_Count_ExpectedReturnData() {
	data := suite.mock.
		NewRows([]string{"count"}).
		AddRow(50)
	q := "SELECT COUNT(*) FROM monsters"
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).WillReturnRows(data)
	res := suite.repo.Count(context.TODO())
	require.NotNil(suite.T(), res)
	require.EqualValues(suite.T(), res, 50)
}
func (suite *monsterSQLRepositoryTestSuite) TestRepository_Count_ExpectedReturnErrorFromQuery() {
	q := "SELECT COUNT(*) FROM monsters"
	expectedQuery := regexp.QuoteMeta(q)
	suite.mock.ExpectQuery(expectedQuery).WillReturnError(errors.New(""))
	res := suite.repo.Count(context.TODO())
	require.EqualValues(suite.T(), res, 0)
}

func TestMonsterRepository(t *testing.T) {
	suite.Run(t, new(monsterSQLRepositoryTestSuite))
}
