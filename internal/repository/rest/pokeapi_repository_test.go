package rest_test

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type pokeapiRESTRepositoryTestSuite struct {
	suite.Suite
	//mock func(url string, v interface{}) error
}

func (suite *pokeapiRESTRepositoryTestSuite) SetupSuite() {

}

// ============
// TODO: HERE
// ============

func TestPokeapiRESTRepository(t *testing.T) {
	suite.Run(t, new(pokeapiRESTRepositoryTestSuite))
}
