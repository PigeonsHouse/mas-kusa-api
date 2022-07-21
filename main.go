package main

import (
	"fmt"

	"mas-kusa-api/db"
	"mas-kusa-api/routers"
	"mas-kusa-api/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.LoadEnv()

	db.InitDB()

	api := gin.Default()
	routers.InitRouter(api)

	api.Static("/static", "static")

	api.Run(fmt.Sprintf(":%s", utils.ApiPort))
}
