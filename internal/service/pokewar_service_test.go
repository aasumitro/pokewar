package service_test

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type pokewarServiceTestSuite struct {
	suite.Suite
}

func (suite *pokewarServiceTestSuite) SetupSuite() {}

// ============
// TODO: HERE
// ============

func TestPokewarService(t *testing.T) {
	suite.Run(t, new(pokewarServiceTestSuite))
}
