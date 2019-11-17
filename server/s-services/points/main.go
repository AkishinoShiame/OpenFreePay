package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

//respect to twd / per
var points = map[string]float32{
	"line":    1,
	"rakuten": 1,
	"shopee":  1,
	"p_point": 1,
}

func main() {
	r := gin.Default()

	r.GET("/calcPointExchange", func(c *gin.Context) {
		src := c.Query("src")
		dst := c.Query("dst")
		quantity := c.Query("quantity")

		q, err := strconv.Atoi(quantity)
		if err != nil {
			c.JSON(200, gin.H{
				"error": true,
				"msg":   "connot conv quantity",
			})
			return
		}

		total := (points[src] * points[dst]) * float32(q)
		c.String(200, fmt.Sprintf("%.2f", total))
	})

	r.Run(":80")
}
