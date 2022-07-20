package cruds

import (
	"errors"
	"mas-kusa-api/db"
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
