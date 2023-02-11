package common

import (
	"fmt"

	"github.com/adarsh-a-tw/passwordly/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

func ConfigureDB(dbType DBType, config DBConfigurable) {
	var err error
	switch dbType {
	case Sqlite3:
		Db, err = gorm.Open(sqlite.Open(config.getDSN()), &gorm.Config{})
	case PostgresDB:
		Db, err = gorm.Open(postgres.Open(config.getDSN()), &gorm.Config{})
	}
	if err != nil {
		panic(err)
	}
	MigrateDB()
}

func MigrateDB() {
	Db.AutoMigrate(&models.User{})
}

type DBConfigurable interface {
	getDSN() string
}

type DBType int

const (
	Sqlite3 DBType = iota
	PostgresDB
)

type DBConfig struct {
	Config DBConfigurable
	Type DBType
}

func (s *SqliteDBConfig) getDSN() string {
	return s.Filename
}

func (p *PostgresDBConfig) getDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", p.Host, p.User, p.Password, p.Name, p.Port) 
}

type SqliteDBConfig struct {
	Filename string
}

type PostgresDBConfig struct {
	Host string
	Port int64
	User string
	Password string
	Name string
}

