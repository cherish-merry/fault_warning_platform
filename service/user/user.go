package user

import (
	"github.com/RaymondCode/simple-demo/database"
	"github.com/RaymondCode/simple-demo/models/user"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Register(user *user.User) error {
	var (
		db = database.GetInstanceConnection().GetPrimaryDB()
	)
	if err := database.WithTransaction(db, func(tx *gorm.DB) error {
		return user.Create(db)
	}); err != nil {
		log.Errorf("create user %v fail, err: %v", user, err)
		return err
	}
	return nil
}
