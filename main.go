package main

import (
	"github.com/adarsh-a-tw/passwordly/common"
	"github.com/adarsh-a-tw/passwordly/users"
	"github.com/gin-gonic/gin"
)

func migrate() {
	common.DB().AutoMigrate(&users.User{})
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

	users.SetupRoutes(r, common.DB())

	r.Run()
}
