package database

import (
	"gorm.io/gorm/schema"
	"sync"

	"github.com/RaymondCode/simple-demo/conf"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type InstanceConnection struct {
}

var (
	primaryDb          *gorm.DB
	secondaryDb        *gorm.DB
	once               sync.Once
	instanceConnection *InstanceConnection
)

func GetInstanceConnection() *InstanceConnection {
	if instanceConnection == nil {
		once.Do(func() {
			instanceConnection = &InstanceConnection{}
		})
	}
	return instanceConnection
}

func (c *InstanceConnection) Init() error {
	var err error

	primaryDatabase := conf.PrimaryDatabase
	dsn := primaryDatabase.User + ":" + primaryDatabase.Passwd + "@tcp(" + primaryDatabase.Host + ")/" + primaryDatabase.Name + "?charset=utf8mb4&parseTime=True&loc=Local"
	primaryDb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		TablePrefix:   "t_",
		SingularTable: true,
	}})
	if err != nil {
		log.Errorf("main(): in ConfigLoad() error:%s", err.Error())
		return err
	}

	return nil
}

func (c *InstanceConnection) GetPrimaryDB() *gorm.DB {
	return primaryDb
}

func (c *InstanceConnection) GetSecondaryDB() *gorm.DB {
	return secondaryDb
}

func DB() *gorm.DB {
	return primaryDb
}

type Transaction func(db *gorm.DB) error

func WithTransaction(tdb *gorm.DB, fn Transaction) (err error) {
	tx := tdb.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error // err is nil; if Commit returns error update err
		}
	}()
	err = fn(tx)
	return err
}

func SetDB(initializedDB *gorm.DB) {
	primaryDb = initializedDB
}
