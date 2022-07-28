package routers

import (
	"image/png"
	"mas-kusa-api/cruds"
	"mas-kusa-api/db"
	"mas-kusa-api/utils"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func defineRouter(r *gin.RouterGroup) {
	r.GET("/health", hello)

	r.POST("/users/signup", setMastodonInfo)
	r.POST("/users/signin", generateJWT)
	r.DELETE("/users", middleware, deleteMastodonInfo)

	r.GET("/generate", middleware, generateTempMaskusa)
}

func hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hello!",
	})
}

func setMastodonInfo(c *gin.Context) {
	var (
		payload utils.SignUpPayload
		user    db.User
		acct    string
		err     error
	)

	c.Bind(&payload)
	if payload.Instance == "" || payload.Token == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "invalid payload",
		})
		return
	}
	if instance := payload.Instance; len(instance) > 1 && instance[len(instance)-1:] == "/" {
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
	var (
		payload utils.SignUpPayload
		jwt     utils.JwtInfo
		err     error
	)

	c.Bind(&payload)
	if payload.Instance == "" || payload.Token == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "invalid payload",
		})
		return
	}
	if instance := payload.Instance; len(instance) > 1 && instance[len(instance)-1:] == "/" {
		payload.Instance = instance[:len(instance)-1]
	}

	if jwt, err = cruds.GenerateJWT(payload.Instance, payload.Token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, jwt)
}

func deleteMastodonInfo(c *gin.Context) {
	var (
		userId  any
		isExist bool
		err     error
	)

	if userId, isExist = c.Get("user_id"); !isExist {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "token is invalid",
		})
		return
	}

	if err = cruds.DeleteInfo(userId.(string)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func generateTempMaskusa(c *gin.Context) {
	var (
		u       db.User
		userId  any
		isExist bool
		err     error
	)

	if userId, isExist = c.Get("user_id"); !isExist {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "token is invalid",
		})
		return
	}

	if u, err = cruds.UserInfoFromUserId(userId.(string)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	now := time.Now()
	baseInstance := strings.Replace(strings.Replace(strings.Replace(u.Instance, "https://", "", -1), "http://", "", -1), "/", "", -1)
	savingPath := "static/" + baseInstance + "/" + u.Name + "-" + now.Format("20060102") + ".png"
	urlImagePath := "/" + savingPath
	if _, err := os.Stat("static/" + baseInstance); err != nil {
		os.Mkdir("static/"+baseInstance, 0777)
	}

	if _, err := os.Stat(savingPath); err == nil {
		c.JSON(http.StatusOK, utils.ImagePath{Path: urlImagePath, Refresh: false})
		return
	}

	baseDate, tootList, err := utils.CountToot(u.Instance, u.Token, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	weekDayNum := int(baseDate.Weekday())
	wholeTootCounter := [][7]int{}
	weekCount := [7]int{}
	for _, toot := range tootList {
		weekCount[weekDayNum] = toot
		weekDayNum++
		if weekDayNum >= 7 {
			wholeTootCounter = append(wholeTootCounter, weekCount)
			weekCount = [7]int{}
			weekDayNum = 0
		}
	}
	wholeTootCounter = append(wholeTootCounter, weekCount)

	baseImage := utils.GenKusa(wholeTootCounter)
	var imagePath *os.File
	if imagePath, err = os.Create(savingPath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	png.Encode(imagePath, baseImage)

	c.JSON(http.StatusOK, utils.ImagePath{
		Path:    urlImagePath,
		Refresh: true,
	})
}
