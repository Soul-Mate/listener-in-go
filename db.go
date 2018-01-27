package main

import (
	"bytes"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

var db *sql.DB

func GetConnect() *sql.DB {
	if db == nil {
		db = NewConnect()
	}
	return db
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


func NewConnect() *sql.DB {
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
