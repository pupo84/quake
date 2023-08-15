package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	jwt "gopkg.in/dgrijalva/jwt-go.v3"
)

// JWTAuthentication is a middleware function that handles JWT authentication for API routes.
// It checks if the request URL path is not in the list of routes to skip authentication and if the request contains a valid JWT token.
// If the token is invalid or not provided, it returns a 401 Unauthorized response.
func JWTAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		routes := viper.GetString("API_AUTH_SKIP_ROUTES")

		for _, route := range strings.Split(routes, ",") {
			if route != c.Request.URL.Path && !strings.Contains(c.Request.URL.Path, "/docs/") {
				var token string
				bearerToken := c.Request.Header.Get("Authorization")

				if len(strings.Split(bearerToken, " ")) == 2 {
					token = strings.Split(bearerToken, " ")[1]

					_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
						if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
							return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
						}
						return []byte(viper.GetString("API_SECRET")), nil
					})

					if err != nil {
						c.JSON(401, gin.H{
							"message": "Not Authorized to perform this request",
							"error":   err.Error(),
						})
						c.Abort()
						return
					}

				} else if len(bearerToken) == 0 {
					c.JSON(401, gin.H{
						"message": "Not Authorized to perform this request",
						"error":   "No token provided",
					})
					c.Abort()
					return
				}
			} else {
				return
			}
		}

		c.Next()
	}
}
