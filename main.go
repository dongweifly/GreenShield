package main

import (
	"flag"
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"green_shield/comm"
	"green_shield/controller"
	"green_shield/core"
	"os"
)

var (
	//出现敏感词后，替换后的字符
	replaceDelim = '*'
	logPath      = "/data/logs/green_shield"
)

func init() {

	log.SetFormatter(&nested.Formatter{
		HideKeys: true,
		NoColors: true,
	})

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

var (
	confFile         = flag.String("c", "conf/conf_local.toml", "Config file")
	sensitiveRootDir = flag.String("s", "./sensitive-words", "sensitive words dir")
)

func main() {
	flag.Parse()

	defer comm.PanicRecovery(true)

	if err := InitConfig(*confFile); err != nil {
		log.Fatal("Init config error!")
		return
	}

	if Config.Env == "local" {
		log.SetOutput(os.Stdout)
		gin.SetMode(gin.DebugMode)
	} else {
		log.SetOutput(&lumberjack.Logger{
			Filename:   logPath + "/green_shield.log",
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     7, //days)
		})

		gin.SetMode(gin.ReleaseMode)
	}

	//从V1.3开始，DB是必须要启动的；
	if err := controller.Init(Config.DBAddress); err != nil {
		log.Fatal("Init DB fail ", err.Error())
		return
	}

	//默认是File
	var wordLoader core.WordsLoad
	wordLoader = core.NewLocalWordLoad(*sensitiveRootDir)
	//从数据库中获取词库
	if Config.UseRepo == "DB" {
		wordLoader = core.NewDatabaseWordLoad(Config.DBAddress)
		if wordLoader == nil {
			fmt.Println("Init Database fail")
			return
		}
	}

	core.GWordsMatchManager = core.NewWordsMatchManager(wordLoader)
	if err := core.GWordsMatchManager.Init(); err != nil {
		log.Warn("Init wordsMatchManager error!")
		return
	}

	httpServerRun()
}
