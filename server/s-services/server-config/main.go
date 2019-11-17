package main

import (
	"github.com/gin-gonic/gin"
)

var cross_auth_key = "rup j03t s.u6q96"

func main() {

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"user_server":     "http://18.179.14.255",
			"verify_server":   "http://54.250.246.160",
			"points_server":   "http://18.176.60.188",
			"market_server":   "http://52.194.252.79",
			"orderbid_server": "http://52.193.225.229",
		})
	})

	r.Run(":8080")
}
