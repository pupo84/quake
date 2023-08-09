package cmd

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pupo84/quake/config"
	"github.com/pupo84/quake/domain/usecase"
	"github.com/pupo84/quake/infra/api/controller"
	"github.com/pupo84/quake/infra/api/middleware"
	"github.com/pupo84/quake/infra/repository"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type APIServer struct {
	server *gin.Engine
}

func NewAPIServer() *APIServer {
	config.Load()

	gin.SetMode(viper.GetString("GIN_MODE"))

	server := gin.New()

	server.Use(gin.Logger())
	server.Use(gin.Recovery())
	// Be aware with CORS in production
	server.Use(middleware.CORS())
	// JWT Authentication
	server.Use(middleware.JWTAuthentication())
	// Swagger configuration
	server.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logger := config.NewLogger()

	fileRepository := repository.NewFileRepository(logger)
	cacheRepository := repository.NewCacheRepository(logger)
	gameUseCase := usecase.NewGameUsecase(logger, fileRepository, cacheRepository)
	gameController := controller.NewGameController(gameUseCase)

	games := server.Group("/v1/games")
	games.GET("", gameController.Get)
	games.GET("/:gameID", gameController.GetByID)

	hc := server.Group("/v1/healthcheck")
	hc.GET("", controller.HealthCheck)

	return &APIServer{server}
}

func (a *APIServer) Run() {
	address := viper.GetString("SERVER_ADDRESS")
	port := viper.GetInt("SERVER_PORT")
	fmt.Printf("Server running on %s:%d\n", address, port)
	a.server.Run(fmt.Sprintf("%s:%d", address, port))
}
