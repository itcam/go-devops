package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/itcam/go-devops/api/app/controller"
	"github.com/itcam/go-devops/api/config"

	"github.com/spf13/viper"
)

func main() {
	cfgTmp := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()
	cfg := *cfgTmp
	if *version {
		fmt.Println(config.VERSION)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	//viper.AddConfigPath(".")
	//viper.AddConfigPath("/")
	viper.AddConfigPath("./config")
	//viper.AddConfigPath("./api/config")
	cfg = strings.Replace(cfg, ".json", "", 1)
	viper.SetConfigName(cfg)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = config.InitLog(viper.GetString("log_level"))
	if err != nil {
		log.Fatal(err)
	}
	err = config.InitDB(viper.GetBool("db.db_bug"), viper.GetViper())
	if err != nil {
		log.Fatalf("db conn failed with error1 %s", err.Error())
	}
	defer config.CloseDB()

	if viper.GetString("log_level") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	routes := gin.Default()
	log.Debugf("will start with port:%v", viper.GetString("web_port"))
	go controller.StartGin(viper.GetString("web_port"), routes)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		os.Exit(0)
	}()
	select {}
}
