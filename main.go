package main

import (
	"net/http"

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
	common.LoadConfig()

	if common.Cfg.DBDriver == "postgres" {
		common.ConfigureDB(
			common.PostgresDB,
			&common.PostgresDBConfig{
				SourceUrl: common.Cfg.DBSource,
			},
		)

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

	r.GET("/health", func(ctx *gin.Context) { ctx.Status(http.StatusOK) })

	r.Run()
}
