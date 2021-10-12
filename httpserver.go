package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"green_shield/controller"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, _ := ioutil.ReadAll(c.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.
		c.Request.Body = rdr2
		c.Next()

		log.Debugf("Header: %v Body: %s", c.Request.Header, readBody(rdr1))
	}
}

type bodyWriter struct {
	gin.ResponseWriter
	bodyBuf *bytes.Buffer
}

func (w bodyWriter) Write(b []byte) (int, error) {
	//memory copy here!
	w.bodyBuf.Write(b)
	return w.ResponseWriter.Write(b)
}

// UploadRetrievalLog ...
func UploadRetrievalLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBodyBytes []byte
		if c.Request.Body != nil {
			requestBodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBodyBytes))
		var responseBodyWriter bodyWriter
		//response写缓存
		responseBodyWriter = bodyWriter{
			bodyBuf:        bytes.NewBufferString(""),
			ResponseWriter: c.Writer}
		c.Writer = responseBodyWriter
		c.Next()
		responseBody := strings.Trim(responseBodyWriter.bodyBuf.String(), "\n")

		log.Infof("REQUEST URI : %s RESP: %s", c.Request.URL.RequestURI(), responseBody)
	}
}

// Cors 解决跨域问题
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		//origin := c.Request.Header.Get("Origin") //请求头部
		//if origin != "" {
		//接收客户端发送的origin （重要！）
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		//服务器支持的所有跨域请求的方法
		c.Header("Access-Control-Allow-Methods", "*")
		//允许跨域设置可以返回其他子段，可以自定义字段
		c.Header("Access-Control-Allow-Headers", "*")
		// 允许浏览器（客户端）可以解析的头部 （重要）
		c.Header("Access-Control-Expose-Headers", "*")
		//设置缓存时间
		c.Header("Access-Control-Max-Age", "172800")
		//允许客户端传递校验信息比如 cookie (重要)
		c.Header("Access-Control-Allow-Credentials", "true")
		//}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}

func httpServerRun() {
	router := gin.Default()

	//GIN的日志跟主程的日志是放在一起的, 同时忽略一部分日志打印
	router.Use(gin.LoggerWithWriter(&lumberjack.Logger{
		Filename:   logPath + "/.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     7, //days
	}, "/actuator/health", "/status/nginx-status"))

	if Config.Env == "local" || Config.Env == "dev" {
		router.Use(UploadRetrievalLog())
	}

	router.Use(Cors())
	router.POST("/filter", controller.TextMatchHandler)

	router.POST("/words/modify", controller.ModifyWordsHandler)
	router.POST("/words/search", controller.SearchWordsHandler)
	router.POST("/words/fuzzySearch", controller.FuzzySearchWordsHandler)
	router.POST("/words/add", controller.AddWordsHandler)
	router.POST("/words/delete", controller.RemoveWordByIdHandler)
	router.POST("/words/exportAll", controller.ExportAllWords)
	router.POST("/words/import", controller.ImportWordsHandler)

	//词库信息维护
	router.GET("/words-info/query", controller.WordsInfoQuery)
	router.GET("/words-info/names", controller.WordsInfoNames)

	//操作记录相关
	router.POST("/record/add", controller.AddRecordHandler)
	router.POST("/record/search", controller.SearchRecordHandler)

	//给发布系统做健康检查使用
	router.GET("/actuator/health", func(c *gin.Context) {
		c.Status(200)
	})

	log.Info("Start http server at " + Config.HTTPServerAddr)

	if err := router.Run(Config.HTTPServerAddr); err != nil {
		fmt.Printf("Start http server %s fail : %s", Config.HTTPServerAddr, err.Error())
		log.Fatal(err)
	}
}
