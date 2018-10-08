package tools

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

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
	return sql.Open("mysql", dsn)
}

func InitDB(conf *DBConfig) (*sql.DB, error) {
	db, err := initDB(conf)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}
