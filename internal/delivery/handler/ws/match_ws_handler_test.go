package ws_test

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type matchWSHandlerTestSuite struct {
	suite.Suite
}

func (suite *matchWSHandlerTestSuite) SetupSuite() {}

// =====
// TODO:
// =====

func TestMatchWSHandler(t *testing.T) {
	suite.Run(t, new(matchWSHandlerTestSuite))
}
