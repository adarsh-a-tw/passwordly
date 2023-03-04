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
}

func main() {
	common.ConfigureDB(
		common.Sqlite3,
		&common.SqliteDBConfig{
			Filename: "sqlite3.db",
		},
	)
	migrate()

	r := gin.Default()

	users.RegisterValidations()

	db := common.DB()

	users.SetupRoutes(r, db)
	vaults.SetupRoutes(r, db)

	r.Run()
}
