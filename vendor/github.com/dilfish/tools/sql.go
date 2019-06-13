// Copyright 2018 Sean.ZH

package tools

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // Just import mysql, no need to specify it
	"strconv"
)

// DBConfig for mysql
type DBConfig struct {
	User   string `json:"user"`
	Pass   string `json:"pass"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
	DBName string `json:"db"`
}

func initDB(conf *DBConfig) (*sql.DB, error) {
	dsn := conf.User + ":" + conf.Pass + "@tcp"
	dsn = dsn + "(" + conf.Host + ":"
	dsn = dsn + strconv.Itoa(conf.Port) + ")"
	dsn = dsn + "/" + conf.DBName
    dsn = dsn + "?timeout=10s"
	return sql.Open("mysql", dsn)
}

// InitDB create new db object for mysql
func InitDB(conf *DBConfig) (*sql.DB, error) {
	db, err := initDB(conf)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}
