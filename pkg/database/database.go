package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres
	"github.com/vincentcent6/product-catalogue/pkg/config"
)

var (
	db *sqlx.DB
)

func InitConnection() error {
	var err error
	dbCfg := config.Get().Database

	connString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", dbCfg.User, dbCfg.Password, dbCfg.DBName, dbCfg.Host, dbCfg.Port)
	db, err = sqlx.Connect(dbCfg.Driver, connString)
	if err != nil {
		return err
	}
	db.SetMaxIdleConns(dbCfg.MaxIdleConns)
	db.SetMaxOpenConns(dbCfg.MaxOpenConns)
	db.SetConnMaxIdleTime(time.Duration(dbCfg.ConnMaxLifetime) * time.Second)

	return nil
}

func GetConnection() (*sqlx.DB, error) {
	var err error
	if db == nil {
		err = InitConnection()
	}
	return db, err
}
