package routers

import (
	"fmt"
	"strings"

	"mas-kusa-api/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func InitRouter(api *gin.Engine) {
	api.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"*",
		},
		AllowMethods: []string{
			"POST",
			"GET",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
	}))
	v1 := api.Group("/api/v1")
	defineRouter(v1)
}

func middleware(c *gin.Context) {
	authorizationHeader := c.Request.Header.Get("Authorization")
	if authorizationHeader != "" {
		ary := strings.Split(authorizationHeader, " ")
		if len(ary) == 2 {
			if ary[0] == "Bearer" {
				t, err := jwt.ParseWithClaims(ary[1], &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
					return utils.SigningKey, nil
				})

				if claims, ok := t.Claims.(*jwt.MapClaims); ok && t.Valid {
					userId := (*claims)["sub"].(string)
					c.Set("user_id", userId)
				} else {
					fmt.Println(err)
				}
			}
		}
	}
	c.Next()
}
