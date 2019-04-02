package main

import (
	"github.com/gin-gonic/gin"
)

var DB = make(map[string]int)

func main() {
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := DB[user]
		if ok {
			c.JSON(200, gin.H{"user": user, "value": value})
		} else {
			c.JSON(200, gin.H{"user": user, "status": "no value"})
		}
	})

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"rvasily": "100500", // user:rvasily password:`100500`
	}))

	authorized.GET("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var params struct {
			Value int `form:"user_key" json:"user_key" binding:"required"`
		}

		err := c.Bind(&params)
		if err != nil {
			c.JSON(400, gin.H{"status": "validation_error", "error": err})
			return
		}

		DB[user] = params.Value
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.Run(":8080")
}
