package service_test

import (
	"context"
	"errors"
	"github.com/aasumitro/pokewar/domain"
	"github.com/aasumitro/pokewar/internal/service"
	"github.com/aasumitro/pokewar/mocks"
	"github.com/aasumitro/pokewar/pkg/appconfigs"
	"github.com/aasumitro/pokewar/pkg/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type pokewarServiceTestSuite struct {
	suite.Suite
	monsters []*domain.Monster
	ranks    []*domain.Rank
	svcErr   *utils.ServiceError
}

func (suite *pokewarServiceTestSuite) SetupSuite() {
	viper.SetConfigFile("../../.example.env")

	appconfigs.LoadEnv()

	suite.monsters = []*domain.Monster{
		{
			ID:       1,
			OriginID: 2,
			Name:     "test",
			BaseExp:  1,
			Height:   1,
			Weight:   1,
			Avatar:   "test.png",
			Types:    []string{"test", "test"},
			Stats:    []domain.Stat{{BaseStat: 1, Name: "test"}},
			Skills:   []*domain.Skill{{PP: 1, Name: "test"}},
		},
	}

	suite.ranks = []*domain.Rank{
		{
			ID:           1,
			OriginID:     1,
			Name:         "test",
			Avatar:       "test.png",
			T:            "['test']",
			Types:        []string{"test"},
			TotalBattles: 1,
			WinBattles:   1,
			LoseBattle:   0,
			Points:       5,
		},
	}

	suite.svcErr = &utils.ServiceError{
		Code:    500,
		Message: "UNEXPECTED",
	}
}

// ============= MONSTERS
func (suite *pokewarServiceTestSuite) TestService_MonstersCount_ShouldSuccess() {
	repo := new(mocks.IMonsterRepository)
	svc := service.NewPokewarService(
		context.TODO(), new(mocks.IPokeapiRESTRepository), repo,
		new(mocks.IRankRepository), new(mocks.IBattleRepository))
	repo.
		On("Count", mock.Anything).
		Once().
		Return(10, nil)
	data := svc.MonstersCount()
	require.NotNil(suite.T(), data)
	require.Equal(suite.T(), data, 10)
	repo.AssertExpectations(suite.T())
}
func (suite *pokewarServiceTestSuite) TestService_MonstersCount_ShouldError() {
	repo := new(mocks.IMonsterRepository)
	svc := service.NewPokewarService(
		context.TODO(), new(mocks.IPokeapiRESTRepository), repo,
		new(mocks.IRankRepository), new(mocks.IBattleRepository))
	repo.
		On("Count", mock.Anything).
		Once().
		Return(0, errors.New(""))
	data := svc.MonstersCount()
	require.NotNil(suite.T(), data)
	require.Equal(suite.T(), data, 0)
	repo.AssertExpectations(suite.T())
}

func (suite *pokewarServiceTestSuite) TestService_FetchMonsters_ShouldSuccess() {
	repo := new(mocks.IMonsterRepository)
	svc := service.NewPokewarService(
		context.TODO(), new(mocks.IPokeapiRESTRepository), repo,
		new(mocks.IRankRepository), new(mocks.IBattleRepository))
	repo.
		On("All", mock.Anything).
		Once().
		Return(suite.monsters, nil)
	data, err := svc.FetchMonsters()
	require.Nil(suite.T(), err)
	require.NotNil(suite.T(), data)
	require.Equal(suite.T(), data, suite.monsters)
	repo.AssertExpectations(suite.T())
}
func (suite *pokewarServiceTestSuite) TestService_FetchMonsters_ShouldError() {
	repo := new(mocks.IMonsterRepository)
	svc := service.NewPokewarService(
		context.TODO(), new(mocks.IPokeapiRESTRepository), repo,
		new(mocks.IRankRepository), new(mocks.IBattleRepository))
	repo.
		On("All", mock.Anything).
		Once().
		Return(nil, errors.New("UNEXPECTED"))
	data, err := svc.FetchMonsters()
	require.Nil(suite.T(), data)
	require.NotNil(suite.T(), err)
	require.Equal(suite.T(), err, suite.svcErr)
	repo.AssertExpectations(suite.T())
}

//func (suite *pokewarServiceTestSuite) TestService_SyncMonsters_ShouldSuccess() {}
//func (suite *pokewarServiceTestSuite) TestService_SyncMonsters_ShouldError() {}

// ============= RANKS
func (suite *pokewarServiceTestSuite) TestService_FetchRanks_ShouldSuccess() {
	repo := new(mocks.IRankRepository)
	svc := service.NewPokewarService(
		context.TODO(), new(mocks.IPokeapiRESTRepository),
		new(mocks.IMonsterRepository),
		repo, new(mocks.IBattleRepository))
	repo.
		On("All", mock.Anything).
		Once().
		Return(suite.ranks, nil)
	data, err := svc.FetchRanks()
	require.Nil(suite.T(), err)
	require.NotNil(suite.T(), data)
	require.Equal(suite.T(), data, suite.ranks)
	repo.AssertExpectations(suite.T())
}
func (suite *pokewarServiceTestSuite) TestService_FetchRanks_ShouldError() {
	repo := new(mocks.IRankRepository)
	svc := service.NewPokewarService(
		context.TODO(), new(mocks.IPokeapiRESTRepository),
		new(mocks.IMonsterRepository),
		repo, new(mocks.IBattleRepository))
	repo.
		On("All", mock.Anything).
		Once().
		Return(nil, errors.New("UNEXPECTED"))
	data, err := svc.FetchRanks()
	require.Nil(suite.T(), data)
	require.NotNil(suite.T(), err)
	require.Equal(suite.T(), err, suite.svcErr)
	repo.AssertExpectations(suite.T())
}

// ============= BATTLES
func (suite *pokewarServiceTestSuite) TestService_BattlesCount_ShouldSuccess() {
	repo := new(mocks.IBattleRepository)
	svc := service.NewPokewarService(
		context.TODO(), new(mocks.IPokeapiRESTRepository),
		new(mocks.IMonsterRepository),
		new(mocks.IRankRepository), repo)
	repo.
		On("Count", mock.Anything).
		Once().
		Return(10, nil)
	data := svc.BattlesCount()
	require.NotNil(suite.T(), data)
	require.Equal(suite.T(), data, 10)
	repo.AssertExpectations(suite.T())
}
func (suite *pokewarServiceTestSuite) TestService_BattlesCount_ShouldError() {
	repo := new(mocks.IBattleRepository)
	svc := service.NewPokewarService(
		context.TODO(), new(mocks.IPokeapiRESTRepository),
		new(mocks.IMonsterRepository),
		new(mocks.IRankRepository), repo)
	repo.
		On("Count", mock.Anything).
		Once().
		Return(0, errors.New(""))
	data := svc.BattlesCount()
	require.NotNil(suite.T(), data)
	require.Equal(suite.T(), data, 0)
	repo.AssertExpectations(suite.T())
}

//func (suite *pokewarServiceTestSuite) TestService_FetchBattles_ShouldSuccess() {}
//func (suite *pokewarServiceTestSuite) TestService_FetchBattles_ShouldError()   {}

func (suite *pokewarServiceTestSuite) TestService_PrepareMonstersForBattle_ShouldSuccess() {
	repo := new(mocks.IMonsterRepository)
	svc := service.NewPokewarService(
		context.TODO(), new(mocks.IPokeapiRESTRepository), repo,
		new(mocks.IRankRepository), new(mocks.IBattleRepository))
	repo.
		On("All", mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return(suite.monsters, nil)
	data, err := svc.PrepareMonstersForBattle()
	require.Nil(suite.T(), err)
	require.NotNil(suite.T(), data)
	require.Equal(suite.T(), data, suite.monsters)
	repo.AssertExpectations(suite.T())
}
func (suite *pokewarServiceTestSuite) TestService_PrepareMonstersForBattle_ShouldError() {
	repo := new(mocks.IMonsterRepository)
	svc := service.NewPokewarService(
		context.TODO(), new(mocks.IPokeapiRESTRepository), repo,
		new(mocks.IRankRepository), new(mocks.IBattleRepository))
	repo.
		On("All", mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return(nil, errors.New("UNEXPECTED"))
	data, err := svc.PrepareMonstersForBattle()
	require.Nil(suite.T(), data)
	require.NotNil(suite.T(), err)
	require.Equal(suite.T(), err, suite.svcErr)
	repo.AssertExpectations(suite.T())
}

//func (suite *pokewarServiceTestSuite) TestService_AddBattle_ShouldSuccess() {}
//func (suite *pokewarServiceTestSuite) TestService_AddBattle_ShouldError() {}

func (suite *pokewarServiceTestSuite) TestService_AnnulledPlayer_ShouldSuccess() {
	repo := new(mocks.IBattleRepository)
	svc := service.NewPokewarService(
		context.TODO(), new(mocks.IPokeapiRESTRepository),
		new(mocks.IMonsterRepository),
		new(mocks.IRankRepository), repo)
	t := time.Now().Unix()
	repo.
		On("UpdatePlayer", mock.Anything, mock.Anything).
		Once().
		Return(t, nil)
	data, err := svc.AnnulledPlayer(1)
	require.Nil(suite.T(), err)
	require.NotNil(suite.T(), data)
	require.Equal(suite.T(), data, t)
	repo.AssertExpectations(suite.T())
}
func (suite *pokewarServiceTestSuite) TestService_AnnulledPlayer_ShouldError() {
	repo := new(mocks.IBattleRepository)
	svc := service.NewPokewarService(
		context.TODO(), new(mocks.IPokeapiRESTRepository),
		new(mocks.IMonsterRepository),
		new(mocks.IRankRepository), repo)
	repo.
		On("UpdatePlayer", mock.Anything, mock.Anything).
		Once().
		Return(time.Now().Unix(), errors.New("UNEXPECTED"))
	_, err := svc.AnnulledPlayer(1)
	require.NotNil(suite.T(), err)
	require.Equal(suite.T(), err, suite.svcErr)
	repo.AssertExpectations(suite.T())
}

func TestPokewarService(t *testing.T) {
	suite.Run(t, new(pokewarServiceTestSuite))
}
