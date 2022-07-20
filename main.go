package main

import (
	"fmt"

	"mas-kusa-api/db"
	"mas-kusa-api/routers"
	"mas-kusa-api/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	api := gin.Default()
	utils.LoadEnv()
	db.InitDB()
	routers.InitRouter(api)
	api.Run(fmt.Sprintf(":%s", utils.ApiPort))
}
