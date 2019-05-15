package main

import (
	"flag"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
	"github.com/mkideal/log"
	//"github.com/shopspring/decimal"
	"github.com/bububa/depthtree/depthtreed/common"
	"github.com/bububa/depthtree/depthtreed/handler"
	"github.com/bububa/depthtree/depthtreed/router"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var (
		config     common.Config
		configFlag common.Config
		configPath string
	)

	os.Setenv("CONFIGOR_ENV_PREFIX", "-")

	flag.StringVar(&configPath, "config", "config.toml", "configuration file")
	flag.IntVar(&configFlag.Port, "port", 0, "set port")
	flag.StringVar(&configFlag.LogPath, "log", "", "set log file path without filename")
	flag.BoolVar(&configFlag.Debug, "debug", false, "set debug mode")
	flag.BoolVar(&configFlag.EnableWeb, "web", false, "enable http web server")
	flag.Parse()

	configor.New(&configor.Config{Verbose: configFlag.Debug, ErrorOnUnmatchedKeys: true, Environment: "production"}).Load(&config, configPath)

	if configFlag.Port > 0 {
		config.Port = configFlag.Port
	}
	if configFlag.LogPath != "" {
		config.LogPath = configFlag.LogPath
	}
	if configFlag.EnableWeb {
		config.EnableWeb = configFlag.EnableWeb
	}

	if configFlag.Debug {
		config.Debug = configFlag.Debug
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Error(err.Error())
		return
	}
	var logPath string
	if path.IsAbs(config.LogPath) {
		logPath = config.LogPath
	} else {
		logPath = path.Join(wd, config.LogPath)
	}
	defer log.Uninit(log.InitMultiFileAndConsole(logPath, "moyud.log", log.LvERROR))
	service := common.NewService(config)
	defer service.Close()
	// service.Db.Reconnect()
	service.TreeBase.Open()
	if config.EnableWeb {
		handler.InitHandler(service, config)
		handler.Start()
		if config.Debug {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
		//gin.DisableBindValidation()
		templatePath := path.Join(config.Template, "./*")
		log.Info("Template path: %s", templatePath)
		r := router.NewRouter(templatePath, config)
		log.Info("%s started at:0.0.0.0:%d", config.AppName, config.Port)
		defer log.Info("%s exit from:0.0.0.0:%d", config.AppName, config.Port)
		srv := endless.NewServer(fmt.Sprintf(":%d", config.Port), r)
		srv.SignalHooks[endless.PRE_SIGNAL][syscall.SIGINT] = append(
			srv.SignalHooks[endless.PRE_SIGNAL][syscall.SIGINT],
			func() {
				handler.ExitCh <- struct{}{}
			})
		err = srv.ListenAndServe()
		if err != nil {
			log.Error(err.Error())
		}
	} else {
		exitChan := make(chan struct{}, 1)
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGSTOP, syscall.SIGTERM)
			<-ch
			exitChan <- struct{}{}
			close(ch)
		}()
		<-exitChan
	}
	log.Warn("shutting down")
	service.TreeBase.Flush()
}
