package rest_test

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type pokeapiRESTRepositoryTestSuite struct {
	suite.Suite
}

func (suite *pokeapiRESTRepositoryTestSuite) SetupSuite() {}

// ============
// TODO: HERE
// ============

func TestPokeapiRESTRepository(t *testing.T) {
	suite.Run(t, new(pokeapiRESTRepositoryTestSuite))
}
