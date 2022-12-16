package sql_test

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type battleSQLRepositoryTestSuite struct {
	suite.Suite
	//mock sqlmock.Sqlmock
}

func (suite *battleSQLRepositoryTestSuite) SetupSuite() {}

// ============
// TODO: HERE
// ============

func TestBattleSQLRepository(t *testing.T) {
	suite.Run(t, new(battleSQLRepositoryTestSuite))
}
