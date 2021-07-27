package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogap/logs"
	"io"
	"os"
	"time"
)

const (
	ErrorInParsing = 1001
	ErrorInGetData = 1002
)

func main() {
	loadParameters()
	loadJenkinsProjectConfig(settings.jenkinsProjectConfigFile)
	r := createGinEngine()
	r.POST(settings.notifyUrl, onNotify)
	if e := r.Run(fmt.Sprintf("%s:%d", settings.hookListeningIp, settings.hookListeningPort)); e == nil {
		logs.Info("Listening on ", settings.hookListeningIp, ":", settings.hookListeningPort)
	} else {
		logs.Error(e)
		panic(e)
	}
}

var settings struct {
	hookRequestLogFile       string
	hookMessageLogFile       string
	hookListeningIp          string
	hookListeningPort        int64
	jenkinsHost              string
	jenkinsNotifyUrl         string
	jenkinsUserName          string
	jenkinsUserApiToken      string
	jenkinsProjectConfigFile string
	notifyUrl                string
	verbose                  bool
}

func loadParameters() {
	flag.StringVar(&settings.hookRequestLogFile, "hook-log-file", "hook-request.log", "Hook request log file")
	flag.StringVar(&settings.hookMessageLogFile, "message-log-file", "message.log", "Hook message log")
	flag.StringVar(&settings.hookListeningIp, "h", "", "Server listening host address(IP or hostname).")
	flag.StringVar(&settings.hookListeningIp, "host", "", "Server listening host address(IP or hostname).")
	flag.Int64Var(&settings.hookListeningPort, "p", 8889, "Server listening port.")
	flag.Int64Var(&settings.hookListeningPort, "port", 8889, "Server listening port.")
	flag.BoolVar(&settings.verbose, "verbose", false, "Print debug logs if verbose is set")
	flag.StringVar(&settings.jenkinsHost, "jenkins-host", "http://cd.mimixiche.cn", "Jenkins host address.")
	flag.StringVar(&settings.jenkinsNotifyUrl, "jenkins-url", "/job/<project>/build?token=<token>", "Jenkins notify URL.")
	flag.StringVar(&settings.jenkinsUserName, "jenkins-user-name", "", "Jenkins User Name.")
	flag.StringVar(&settings.jenkinsUserApiToken, "jenkins-api-token", "", "Jenkins User API Token.")
	flag.StringVar(&settings.jenkinsProjectConfigFile, "jenkins-project-config-file", "/etc/prcd/projects.yaml", "Jenkins Project config file.")
	flag.StringVar(&settings.notifyUrl, "notify-url", "/notify", "Listening url address.")
	flag.Parse()
	logs.SetFileLogger(settings.hookMessageLogFile)
	if !settings.verbose {
		logs.SetLoggerLevel(logs.LevelInfo)
	}
	f, _ := os.Create(settings.hookRequestLogFile)
	gin.DefaultWriter = io.MultiWriter(f)
}

func createGinEngine() *gin.Engine {
	r := gin.New()
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())
	return r
}

func onNotify(c *gin.Context) {
	logs.Debug("receive post")
	errorCode, errorMessage := 0, "ok"
	var e error
	if b, err := c.GetRawData(); err == nil {
		logs.Debug(string(b))
		basicHook := BasicHook{}
		if err := json.Unmarshal(b, &basicHook); err == nil {
			go sendNotice(basicHook, b)
		} else {
			e, errorCode = err, ErrorInParsing
		}
	} else {
		e, errorCode = err, ErrorInGetData
	}

	if e != nil {
		errorMessage = e.Error()
		logs.Error(e)
	}
	c.JSON(200, gin.H{"errcode": errorCode, "errmsg": errorMessage,})
}

func sendNotice(basicHook BasicHook, bytes []byte) {
	agent := createHookAgentByName(basicHook.HookName)
	if e := agent.Parse(bytes); e == nil {
		if agent.CanTriggerEvent() {
			notifier := createNotifierByAgent(agent)
			if err := notifier.Notify(); err != nil {
				logs.Error(err)
			}
		} else {
			logs.Debug("Agent cannot trigger event:", agent.Name())
		}
	} else {
		logs.Error(e)
	}
}
