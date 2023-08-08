package conf

import (
	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Database struct {
	Host   string
	Name   string
	User   string
	Passwd string
}

type Amqp struct {
	Host   string
	Port   int
	User   string
	Passwd string
}

type Server struct {
	Port string
}

var PrimaryDatabase = &Database{}

var SecondaryDatabase = &Database{}

var ServerConfig = &Server{}

var AmqpConfig = &Amqp{}

var cfg *ini.File

func createSettingMap() {
	mapTo("primary-database", PrimaryDatabase)
	mapTo("secondary-database", SecondaryDatabase)
	mapTo("server", ServerConfig)
	mapTo("amqp", AmqpConfig)
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo RedisSetting err: %v", err.Error())
	}
}
func ConfigLoad(env string) error {
	var err error

	// dev mapping test config test.ini
	if strings.EqualFold(env, "dev") {
		cfg, err = ini.Load("conf/dev.ini")
	} else if strings.EqualFold(env, "test") {
		// prod mapping prod config prod.ini
		cfg, err = ini.Load("conf/test.ini")
	} else if strings.EqualFold(env, "prod") {
		// prod mapping prod config prod.ini
		cfg, err = ini.Load("conf/prod.ini")
	}

	if err != nil {
		log.Fatalf("Setup, fail to parse config file: %v", err.Error())
		return err
	}

	createSettingMap()

	return nil
}
