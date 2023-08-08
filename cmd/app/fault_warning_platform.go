package app

import (
	"flag"
	"github.com/RaymondCode/simple-demo/amqp"
	"github.com/RaymondCode/simple-demo/daemon"
	"github.com/RaymondCode/simple-demo/database"

	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/conf"
	"github.com/RaymondCode/simple-demo/router"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	help = flag.Bool("h", false, "to show help")
	env  = flag.String("env", "dev", "env: dev | test | prod")
)

func InitPlatformService() {
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return
	}

	// load config
	err := conf.ConfigLoad(*env)
	if err != nil {
		log.Errorf("main(): in ConfigLoad() error:%s", err.Error())
		common.AbnormalExit()
	}

	// init datasource
	err = database.GetInstanceConnection().Init()
	if err != nil {
		log.Errorf("main(): in common_database.GetInstanceConnection() error:%s", err.Error())
		common.AbnormalExit()
	}

	// init message
	err = amqp.InitAmqp()
	if err != nil {
		log.Errorf("main(): in amqp.InitAmqp() error:%s", err.Error())
		common.AbnormalExit()
	}

	// start daemon
	daemon.InitDaemon()

	//init router
	r := gin.Default()
	router.InitRouter(r)
	serverPort := conf.ServerConfig.Port
	err = r.Run("0.0.0.0:" + serverPort)
	if err != nil {
		return
	}
}
