package http_test

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type battleHTTPHandlerTestSuite struct {
	suite.Suite
}

func (suite *battleHTTPHandlerTestSuite) SetupSuite() {}

// =====
// TODO:
// =====

func TestBattleHTTPHandler(t *testing.T) {
	suite.Run(t, new(battleHTTPHandlerTestSuite))
}
