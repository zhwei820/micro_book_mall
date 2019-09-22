package db

import (
	"database/sql"
	"fmt"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/basic/config"
)

type db struct {
	Mysql Mysql `json："mysql"`
}

// Mysql mySQL 配置
type Mysql struct {
	URL               string `json:"url"`
	Enable            bool   `json:"enabled"`
	MaxIdleConnection int    `json:"maxIdleConnection"`
	MaxOpenConnection int    `json:"maxOpenConnection"`
}

func initMysql() {
	log.Info(fmt.Sprintf("[initMysql] 初始化Mysql"))

	c := config.C()
	cfg := &db{}

	err := c.App("db", cfg)
	if err != nil {
		log.Info(fmt.Sprintf("[initMysql] %s", err))
	}

	if !cfg.Mysql.Enable {
		log.Info(fmt.Sprintf("[initMysql] 未启用Mysql"))
		return
	}

	// 创建连接
	mysqlDB, err = sql.Open("mysql", cfg.Mysql.URL)
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}

	// 最大连接数
	mysqlDB.SetMaxOpenConns(cfg.Mysql.MaxOpenConnection)

	// 最大闲置数
	mysqlDB.SetMaxIdleConns(cfg.Mysql.MaxIdleConnection)

	// 激活链接
	if err = mysqlDB.Ping(); err != nil {
		log.Fatal(err.Error())
	}

	log.Info(fmt.Sprintf("[initMysql] Mysql 连接成功"))
}
