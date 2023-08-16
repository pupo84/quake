package usecase

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/pupo84/quake/domain/entity"
	"github.com/pupo84/quake/infra/repository"
	"go.uber.org/zap"
)

// GameUsecase is a struct that represents the usecase for parsing Quake games.
type GameUsecase struct {
	logger          *zap.SugaredLogger
	fileRepository  repository.FileRepositorier
	cacheRepository repository.CacheRepositorier
}

// NewGameUsecase returns a new instance of GameUsecase.
func NewGameUsecase(logger *zap.SugaredLogger, fileRepository repository.FileRepositorier, cacheRepository repository.CacheRepositorier) *GameUsecase {
	return &GameUsecase{
		logger:          logger,
		fileRepository:  fileRepository,
		cacheRepository: cacheRepository,
	}
}

// Go is a function that parses all Quake games and returns a slice of Game entities.
func (uc *GameUsecase) Go(ctx context.Context) (entity.Games, error) {
	contents, hash, err := uc.fileRepository.Read(ctx)
	if err != nil {
		return nil, err
	}

	if uc.cacheRepository.IsAlive() {
		if value, err := uc.cacheRepository.Get(ctx, hash); err == nil {
			uc.logger.Infof("Cache hit for hash %s", hash)
			games := entity.Games{}
			err := json.Unmarshal([]byte(value), &games)
			if err != nil {
				return nil, err
			}
			return games, nil
		}
	}

	games := make(entity.Games, 0)
	currentGame := entity.NewGame(1)
	currentGameID := 1

	killRegex := regexp.MustCompile(`\d+:\d+\s+Kill:\s+\d+\s+\d+\s+\d+:\s+(.*?)\skilled\s+(.*?)\s+by\s+(.*)`)
	playerRegex := regexp.MustCompile(`n\\([^\\]+)`)

	for line_number, content := range contents {
		if uc.startedNewGame(content) || uc.reachedEndOfFile(len(contents), line_number) {
			if currentGame.HasStats() {
				games = append(games, currentGame)
				currentGameID++
				currentGame = entity.NewGame(currentGameID)
			}
		} else {

			var killer, victim, cause string

			if matches := killRegex.FindStringSubmatch(content); len(matches) >= 4 {
				killer, victim, cause = matches[1], matches[2], matches[3]

				if killer != entity.WORD {
					currentGame.AddPlayerStats(killer, cause, 1)
				} else {
					currentGame.AddPlayerStats(victim, cause, -1)
				}
			}

			if matches := playerRegex.FindStringSubmatch(content); len(matches) > 0 {
				player := matches[1]
				currentGame.AddPlayer(player)
			}
		}
	}

	if uc.cacheRepository.IsAlive() {
		var data []byte
		if data, err = json.Marshal(games); err == nil {
			uc.logger.Infof("Saving games to cache with hash %s", hash)
			uc.cacheRepository.Set(ctx, hash, string(data))
		}
	}

	return games, nil
}

// startedNewGame is a function that returns true if the content string contains the "InitGame:" word.
func (uc *GameUsecase) startedNewGame(content string) bool {
	return strings.Contains(content, entity.INIT_GAME)
}

// reachedEndOfFile is a function that returns true if the line_number is the last line of the file.
func (uc *GameUsecase) reachedEndOfFile(file_size, line_number int) bool {
	return (file_size - 1) == line_number
}
