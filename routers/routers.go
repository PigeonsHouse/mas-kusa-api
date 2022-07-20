package routers

import (
	"fmt"
	"mas-kusa-api/cruds"
	"mas-kusa-api/db"
	"mas-kusa-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func defineRouter(r *gin.RouterGroup) {
	r.GET("", hello)
	r.POST("/users/signup", setMastodonInfo)
	r.POST("/users/signin", generateJWT)
	r.POST("/generate", generateTempMaskusa)
}

func hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hello!",
	})
}

func setMastodonInfo(c *gin.Context) {
	var (
		payload SignUpPayload
		user    db.User
		acct    string
		err     error
	)

	c.Bind(&payload)

	fmt.Println(payload)

	if payload.Instance == "" || payload.Token == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "invalid payload",
		})
		return
	}

	if instance := payload.Instance; instance[len(instance)-1:] == "/" {
		payload.Instance = instance[:len(instance)-1]
	}
	if acct, err = utils.GetUserName(payload.Instance, payload.Token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err := cruds.RegisterUser(&user, payload.Instance, acct, payload.Token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

func generateJWT(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "temp",
	})
}

func generateTempMaskusa(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "temp",
	})
}
