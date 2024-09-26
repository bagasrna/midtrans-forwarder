package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func NewDatabase(cfg MySQLConfig) *sql.DB {
	username := cfg.User
	password := cfg.Password
	host := cfg.Host
	port := cfg.Port
	dbName := cfg.DBName
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}