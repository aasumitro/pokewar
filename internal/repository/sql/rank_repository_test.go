package sql_test

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type rankSQLRepositoryTestSuite struct {
	suite.Suite
	//mock sqlmock.Sqlmock
}

func (suite *rankSQLRepositoryTestSuite) SetupSuite() {}

// ============
// TODO: HERE
// ============

func TestRankRepository(t *testing.T) {
	suite.Run(t, new(rankSQLRepositoryTestSuite))
}
