package controller

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pupo84/quake/config"
	"github.com/pupo84/quake/domain/usecase"
	"github.com/pupo84/quake/infra/api/presenter"

	"go.uber.org/zap"
)

// Controller is the interface that wraps the Get and GetByID methods
type Controller interface {
	Get(c *gin.Context)
	GetByID(c *gin.Context, id int)
}

// GameController is the struct that implements the Controller interface
type GameController struct {
	ctx     context.Context
	logger  *zap.SugaredLogger
	useCase *usecase.GameUsecase
}

// NewGameController returns a new GameController
func NewGameController(gameUseCase *usecase.GameUsecase) *GameController {
	return &GameController{
		ctx:     context.Background(),
		logger:  config.NewLogger(),
		useCase: gameUseCase,
	}
}

// Get godoc
// @Summary Get all games
// @Description Get all games
// @Tags games
// @Produces application/json
// @Success 200 {object} presenter.Response{}
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
// @Router /games [get]
func (gc *GameController) Get(c *gin.Context) {
	games, err := gc.useCase.Go(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{
			"message": "got unexpected error while parsing games",
			"error":   err.Error(),
		})
		return
	}

	if len(games) == 0 {
		message := "could not parse games from log file"
		c.JSON(404, gin.H{
			"message": message,
			"error":   errors.New(message),
		})
		return
	}

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(200, presenter.NewGameResponse(games))
}

// GetByID godoc
// @Summary Get a single game
// @Description Get game given an ID
// @Param gameID path int true "Game identifier"
// @Tags games
// @Produces application/json
// @Success 200 {object} presenter.Response{}
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
// @Router /games/{gameID} [get]
func (gc *GameController) GetByID(c *gin.Context) {
	games, err := gc.useCase.Go(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{
			"message": "got unexpected error while parsing games",
			"error":   err.Error(),
		})
		return
	}

	var gameID int
	if gameID, err = strconv.Atoi(c.Param("gameID")); err != nil {
		message := fmt.Sprintf("could not parse game id %s from request path", c.Param("gameID"))
		c.JSON(400, gin.H{
			"message": message,
			"error":   err.Error(),
		})
		return
	}

	if gameID < 1 || gameID > len(games) {
		message := fmt.Sprintf("game with ID %d not found.", gameID)
		c.JSON(404, gin.H{
			"message": message,
			"error":   errors.New(message),
		})
		return
	}

	game := games[gameID-1 : gameID]
	c.JSON(200, presenter.NewGameResponse(game))
}
