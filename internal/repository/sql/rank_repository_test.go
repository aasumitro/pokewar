package sql_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aasumitro/pokewar/domain"
	repoSql "github.com/aasumitro/pokewar/internal/repository/sql"
	"github.com/aasumitro/pokewar/pkg/configs"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type rankSQLRepositoryTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	repo domain.IRankRepository
	// data *domain.Rank
}

func (suite *rankSQLRepositoryTestSuite) SetupSuite() {
	var err error

	configs.DbPool, suite.mock, err = sqlmock.New(
		sqlmock.QueryMatcherOption(
			sqlmock.QueryMatcherRegexp))

	require.NoError(suite.T(), err)

	suite.repo = repoSql.NewRankSQLRepository()

	// suite.data = &domain.Rank{}
}

// =========== ALL
func (suite *rankSQLRepositoryTestSuite) TestRepository_All_Expect() {
	// TODO
}

func TestRankRepository(t *testing.T) {
	suite.Run(t, new(rankSQLRepositoryTestSuite))
}
