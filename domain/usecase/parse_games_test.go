package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pupo84/quake/domain/entity"
	"github.com/pupo84/quake/domain/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type mockFileRepository struct {
	mock.Mock
}

func (m *mockFileRepository) Read(ctx context.Context) ([]string, string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.String(1), args.Error(2)
}

type mockCacheRepository struct {
	mock.Mock
}

func (m *mockCacheRepository) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *mockCacheRepository) Set(ctx context.Context, key string, value string) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func TestGameUsecase_Go(t *testing.T) {
	logger := zap.NewNop().Sugar()
	ctx := context.Background()

	gameData := []string{
		`  0:00 InitGame: \sv_floodProtect\1\sv_maxPing\0\sv_minPing\0\sv_maxRate\10000\sv_minRate\0\sv_hostname\Code Miner Server\g_gametype\0\sv_privateClients\2\sv_maxclients\16\sv_allowDownload\0\dmflags\0\fraglimit\20\timelimit\15\g_maxGameClients\0\capturelimit\8\version\ioq3 1.36 linux-x86_64 Apr 12 2009\protocol\68\mapname\q3dm17\gamename\baseq3\g_needpass\0`,
		` 20:34 ClientUserinfoChanged: 2 n\Player1\t\0\model\xian/default\hmodel\xian/default\g_redteam\\g_blueteam\\c1\4\c2\5\hc\100\w\0\l\0\tt\0\tl\0`,
		` 20:34 ClientUserinfoChanged: 2 n\Player2\t\0\model\xian/default\hmodel\xian/default\g_redteam\\g_blueteam\\c1\4\c2\5\hc\100\w\0\l\0\tt\0\tl\0`,
		` 20:54 Kill: 1022 2 22: Player1 killed Player2 by MOD_RAILGUN`,
		` 20:55 Kill: 1022 2 22: Player2 killed Player1 by MOD_TRIGGER_HURT`,
		` 20:56 Kill: 1022 2 22: <world> killed Player1 by MOD_TRIGGER_HURT`,
		`  0:07 ShutdownGame:`,
	}

	expectedGames := entity.Games{{ID: 1,
		Players:        []string{"Player1", "Player2"},
		Kills:          3,
		KillsByPlayers: map[string]int{"Player1": 0, "Player2": 1},
		KillsByCause:   map[string]int{"MOD_RAILGUN": 1, "MOD_TRIGGER_HURT": 2}}}

	expectedHash := "hash"

	t.Run("should return error if file repository returns error", func(t *testing.T) {
		fileRepo := new(mockFileRepository)
		cacheRepo := new(mockCacheRepository)
		usecase := usecase.NewGameUsecase(logger, fileRepo, cacheRepo)

		fileRepo.On("Read", ctx).Return([]string{}, "", errors.New("Could not parse games from logfile"))

		games, err := usecase.Go(ctx)

		assert.Equal(t, len(games), 0)
		assert.NotNil(t, err)
	})

	t.Run("should return games from cache if cache repository has a hit", func(t *testing.T) {
		fileRepo := new(mockFileRepository)
		cacheRepo := new(mockCacheRepository)
		usecase := usecase.NewGameUsecase(logger, fileRepo, cacheRepo)

		fileRepo.On("Read", ctx).Return(gameData, expectedHash, nil)
		cacheRepo.On("Get", ctx, expectedHash).Return("[{\"ID\":1,\"Kills\":3,\"Players\":[\"Player1\",\"Player2\"],\"KillsByPlayers\":{\"Player1\":0,\"Player2\":1},\"KillsByCause\":{\"MOD_RAILGUN\":1,\"MOD_TRIGGER_HURT\":2}}]", nil)

		games, err := usecase.Go(ctx)

		assert.Equal(t, expectedGames, games)
		assert.Nil(t, err)
	})

	t.Run("should parse games and save to cache if cache repository has a miss", func(t *testing.T) {
		fileRepo := new(mockFileRepository)
		cacheRepo := new(mockCacheRepository)
		usecase := usecase.NewGameUsecase(logger, fileRepo, cacheRepo)

		fileRepo.On("Read", ctx).Return(gameData, expectedHash, nil)
		cacheRepo.On("Get", ctx, expectedHash).Return("", errors.New("Cache miss"))
		cacheRepo.On("Set", ctx, expectedHash, mock.Anything).Return(nil)

		games, err := usecase.Go(ctx)

		assert.Equal(t, expectedGames, games)
		assert.Nil(t, err)
	})

	t.Run("should parse games and not return error if cache repository fails to save", func(t *testing.T) {
		fileRepo := new(mockFileRepository)
		cacheRepo := new(mockCacheRepository)
		usecase := usecase.NewGameUsecase(logger, fileRepo, cacheRepo)

		fileRepo.On("Read", ctx).Return(gameData, expectedHash, nil)
		cacheRepo.On("Get", ctx, expectedHash).Return("", errors.New("Cache miss"))
		cacheRepo.On("Set", ctx, expectedHash, mock.Anything).Return(errors.New("Could not save to cache"))

		games, err := usecase.Go(ctx)

		assert.Equal(t, expectedGames, games)
		assert.Nil(t, err)
	})

	t.Run("should return error when redis data unmarshal fail", func(t *testing.T) {
		fileRepo := new(mockFileRepository)
		cacheRepo := new(mockCacheRepository)
		usecase := usecase.NewGameUsecase(logger, fileRepo, cacheRepo)

		fileRepo.On("Read", ctx).Return(gameData, expectedHash, nil)
		cacheRepo.On("Get", ctx, expectedHash).Return("{}", nil)

		_, err := usecase.Go(ctx)

		assert.NotNil(t, err)
	})
}
