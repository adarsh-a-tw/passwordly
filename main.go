package main

import (
	"github.com/adarsh-a-tw/passwordly/common"
	"github.com/adarsh-a-tw/passwordly/users"
	"github.com/adarsh-a-tw/passwordly/vaults"
	"github.com/gin-gonic/gin"
)

func migrate() {
	db := common.DB()
	db.AutoMigrate(&users.User{})
	db.AutoMigrate(&vaults.Vault{})
	db.AutoMigrate(&vaults.Credential{})
	db.AutoMigrate(&vaults.Key{})
	db.AutoMigrate(&vaults.Document{})
}

func main() {
	err := common.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	if common.Cfg.DBDriver == "postgres" {

	} else {
		common.ConfigureDB(
			common.Sqlite3,
			&common.SqliteDBConfig{
				Filename: common.Cfg.DBSource,
			},
		)
	}
	migrate()

	if common.Cfg.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	users.RegisterValidations()
	vaults.RegisterValidations()

	db := common.DB()

	users.SetupRoutes(r, db)
	vaults.SetupRoutes(r, db)

	r.Run()
}
