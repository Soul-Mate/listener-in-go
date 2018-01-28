package main

import (
	"bytes"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-redis/redis"
	"strings"
	"log"
)

var db *sql.DB

var redisClient *redis.Client

// 获取Mysql连接
func GetMysqlConnect() *sql.DB {
	if db == nil {
		db = NewMysqlConnect()
	}
	return db
}

// 获取redis连接
func GetRedisConnect() *redis.Client {
	if redisClient == nil {
		redisClient = NewRedisConnect()
	}
	return redisClient
}

func GetInsertSql(table string, fields ...string) string {
	var buf bytes.Buffer
	buf.WriteString("INSERT INTO `")
	buf.WriteString(table)
	buf.WriteString("`(id,")
	buf.WriteString(strings.Join(fields, ","))
	buf.WriteString(")")
	buf.WriteString("VALUES(null,")
	fieldsCnt := len(fields)
	placeholders := make([]string, fieldsCnt)
	for i := 0; i < fieldsCnt; i++ {
		placeholders[i] = "?"
	}
	buf.WriteString(strings.Join(placeholders, ","))
	buf.WriteString(")")
	return buf.String()
}

// 建立Mysql连接
func NewMysqlConnect() *sql.DB {
	db, err := sql.Open("mysql", dsn())
	Check(err)
	return db
}

func dsn() string {
	conf := GetConfig()
	var buf bytes.Buffer
	buf.WriteString(conf.Db.User)
	buf.WriteString(":")
	buf.WriteString(conf.Db.Password)
	buf.WriteString("@tcp(")
	buf.WriteString(conf.Db.Host)
	buf.WriteString(":")
	buf.WriteString(conf.Db.Port)
	buf.WriteString(")/")
	buf.WriteString(conf.Db.Database)
	buf.WriteString("?charset=")
	buf.WriteString(conf.Db.Charset)
	return buf.String()
}

// 建立redis连接
func NewRedisConnect() *redis.Client {
	conf := GetConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Host + ":" + conf.Redis.Port,
		Password: "",
		DB:       4,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
	return client
}
