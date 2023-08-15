package controller

import "github.com/gin-gonic/gin"

// HealthCheck godoc
// @Summary Healthcheck endpoint
// @Description Healthcheck endpoint
// @Tags healthcheck
// @Produces application/json
// @Success 200 {object} map[string]interface{}
// @Router /healthcheck [get]
func HealthCheck(c *gin.Context) {
	c.JSON(200, map[string]interface{}{
		"message": "OK",
	})
}
