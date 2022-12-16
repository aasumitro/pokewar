package http_test

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type rankHTTPHandlerTestSuite struct {
	suite.Suite
}

func (suite *rankHTTPHandlerTestSuite) SetupSuite() {}

// =====
// TODO:
// =====

func TestRankHTTPHandler(t *testing.T) {
	suite.Run(t, new(rankHTTPHandlerTestSuite))
}
