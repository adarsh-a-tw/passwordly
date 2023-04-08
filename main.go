package main

import (
	"github.com/adarsh-a-tw/passwordly/common"
	"github.com/adarsh-a-tw/passwordly/secrets"
	"github.com/adarsh-a-tw/passwordly/users"
	"github.com/adarsh-a-tw/passwordly/vaults"
	"github.com/gin-gonic/gin"
)

func migrate() {
	db := common.DB()
	db.AutoMigrate(&users.User{})
	db.AutoMigrate(&vaults.Vault{})
	db.AutoMigrate(&secrets.Secret{})
	db.AutoMigrate(&secrets.Credential{})
	db.AutoMigrate(&secrets.Key{})
	db.AutoMigrate(&secrets.Document{})
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
	secrets.RegisterValidations()

	db := common.DB()

	users.SetupRoutes(r, db)
	vaults.SetupRoutes(r, db)
	secrets.SetupRoutes(r, db)

	r.Run()
}
