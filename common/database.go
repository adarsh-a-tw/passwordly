package common

import (
	"fmt"
	"strconv"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConfigureDB(dbType DBType, config DBConfigurable) {
	var err error
	switch dbType {
	case Sqlite3:
		db, err = gorm.Open(sqlite.Open(config.getDSN()), &gorm.Config{})
	case PostgresDB:
		db, err = gorm.Open(postgres.Open(config.getDSN()), &gorm.Config{})
	}
	if err != nil {
		panic(err)
	}
}

func DB() *gorm.DB {
	if db == nil {
		panic("Database is not configured yet.")
	}
	return db
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
	Type   DBType
}

func (s *SqliteDBConfig) getDSN() string {
	return s.Filename
}

func (p *PostgresDBConfig) getDSN() string {
	ss := strings.Split(p.SourceUrl, "://")
	ss = strings.Split(ss[1], "@")
	up := strings.Split(ss[0], ":")
	user, pwd := up[0], up[1]
	ss = strings.Split(ss[1], "/")
	hp := strings.Split(ss[0], ":")
	host := hp[0]
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		panic(err)
	}
	dbName := ss[1]

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, pwd, dbName, port)
}

type SqliteDBConfig struct {
	Filename string
}

type PostgresDBConfig struct {
	SourceUrl string
}
