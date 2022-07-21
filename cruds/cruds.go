package cruds

import (
	"errors"
	"mas-kusa-api/db"
	"mas-kusa-api/utils"
	"time"

	"github.com/golang-jwt/jwt"
)

func RegisterUser(usr *db.User, instance string, acct string, token string) (err error) {
	if err = db.Psql.Where("instance = ? AND name = ?", instance, acct).First(&db.User{}).Error; err == nil {
		err = errors.New("this account is already registered")
		return
	}

	*usr = db.User{
		Instance: instance,
		Name:     acct,
		Token:    token,
	}

	err = db.Psql.Create(usr).Error
	return
}

func GenerateJWT(instance string, token string) (jwtInfo utils.JwtInfo, err error) {
	var (
		u   db.User
		jwt string
	)

	if err = db.Psql.Where("instance = ? AND token = ?", instance, token).First(&u).Error; err != nil {
		return
	}

	jwt, err = generateToken(u.ID)
	if err != nil {
		return
	}

	jwtInfo = utils.JwtInfo{Jwt: jwt}
	return
}

func UserInfoFromUserId(userId string) (userInfo db.User, err error) {
	err = db.Psql.Where("id = ?", userId).First(&userInfo).Error
	return
}

func generateToken(userID string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(7 * 24 * time.Hour).Unix(),
	})
	return token.SignedString(utils.SigningKey)
}
