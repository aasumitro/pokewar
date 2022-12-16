package sql_test

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type monsterSQLRepositoryTestSuite struct {
	suite.Suite
	//mock sqlmock.Sqlmock
}

func (suite *monsterSQLRepositoryTestSuite) SetupSuite() {}

// ============
// TODO: HERE
// ============

func TestMonsterRepository(t *testing.T) {
	suite.Run(t, new(monsterSQLRepositoryTestSuite))
}
