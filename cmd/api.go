package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pupo84/quake/config"
	"github.com/pupo84/quake/domain/usecase"
	"github.com/pupo84/quake/infra/api/controller"
	"github.com/pupo84/quake/infra/api/middleware"
	"github.com/pupo84/quake/infra/repository"
	"github.com/pupo84/quake/infra/tracing"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// APIServer is the struct that represents the HTTP server
type APIServer struct {
	ctx    context.Context
	server *gin.Engine
}

// NewAPIServer returns a new HTTP server
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
	// OpenTelemetry for Gin
	server.Use(otelgin.Middleware(viper.GetString("SERVICE_NAME")))

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

	return &APIServer{context.Background(), server}
}

// Run starts the HTTP server
func (a *APIServer) Run() {
	shutdown, err := tracing.NewTraceProvider(a.ctx)
	if err != nil {
		log.Printf("Failed to initialize opentelemetry provider: %s", err)
	}
	defer shutdown(a.ctx)

	address := viper.GetString("SERVER_ADDRESS")
	port := viper.GetInt("SERVER_PORT")
	fmt.Printf("Server running on %s:%d", address, port)
	a.server.Run(fmt.Sprintf("%s:%d", address, port))
}
