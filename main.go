package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/adarsh-a-tw/passwordly/common"
)

func main() {
	r := gin.Default()
	common.ConfigureDB(
		common.Sqlite3,
		&common.SqliteDBConfig{
			Filename: "sqlite3.db",
		},
	)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}