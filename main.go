package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"midea-clo/db"
	"midea-clo/handles"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Address string `yaml:"address"`
	}
	Sms struct {
		SecretId  string `yaml:"secretid"`
		SecretKey string `yaml:"secretkey"`
	}
	Mysql struct {
		User     string `yaml:"user"`
		Host     string `yaml:"host"`
		Password string `yaml:"password"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	}
}

var c string

func init() {
	flag.StringVar(&c, "c", "config.yml", "set config path")
}

func main() {
	var err error
	var conf Config

	flag.Parse()

	data, _ := ioutil.ReadFile(c)
	yaml.Unmarshal(data, &conf)

	// 初始化数据库
	db.MysqlDB, err = gorm.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Local",
			conf.Mysql.User,
			conf.Mysql.Password,
			conf.Mysql.Host,
			conf.Mysql.Port,
			conf.Mysql.Name,
		),
	)
	if err != nil {
		fmt.Println("failed to connect database:", err)
		return
	}
	defer db.MysqlDB.Close()

	db.MysqlDB.SingularTable(true)

	// 初始化引擎
	router := gin.Default()

	router.Use(corsMiddleware())

	router.GET("/get-list", handles.List)
	router.GET("/get-info", handles.Info)
	router.GET("/get-statistics", handles.Statistics)
	router.POST("/send-code", handles.SendCode)
	router.POST("/sign-up", handles.SignUp)
	router.POST("/sign", handles.Sign)
	router.POST("/approval", handles.Approval)
	router.GET("/export", handles.Export)

	router.Static("/img", "html/img")
	router.Static("/css", "html/css")
	router.Static("/js", "html/js")
	router.Static("/html", "html")

	router.Run(conf.Server.Address)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 核心处理方式
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
		c.Set("content-type", "application/json")

		//放行所有OPTIONS方法
		if c.Request.Method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}

		c.Next()
	}
}
