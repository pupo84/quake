package main

import (
	"github.com/pupo84/quake/cmd"
	_ "github.com/pupo84/quake/docs"
)

// @title Quake API
// @description This is a Quake API server
// @version 1.0
// @BasePath /v1
// @host 0.0.0.0:8000
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
// @schemes http
func main() {
	server := cmd.NewAPIServer()
	server.Run()
}
