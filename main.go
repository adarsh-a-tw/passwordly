package main

import (
	"net/http"

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
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	users.SetupRoutes(r, common.DB())

	r.Run()
}
