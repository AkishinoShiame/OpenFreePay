package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var cross_auth_key = "rup j03t s.u6q96"

func main() {

	db, err := sql.Open("mysql", "admin:awszo4zo4@tcp(openfreepay-db.ckfrdtxaxflb.ap-northeast-1.rds.amazonaws.com:3306)/innodb")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	r := gin.Default()

	r.POST("/login", func(c *gin.Context) {
		phone := c.PostForm("phone")

		type User struct {
			Id int `json:"id"`
		}

		var user User
		// Execute the query
		err = db.QueryRow("SELECT id FROM users where phone = ?", phone).Scan(&user.Id)
		if err != nil {
			c.JSON(200, gin.H{
				"error": true,
				"msg":   "query failed",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"id":    user.Id,
		})
	})

	r.POST("/createUser", func(c *gin.Context) {
		auth := c.PostForm("auth")
		username := c.PostForm("username")
		phone := c.PostForm("phone")

		if auth != cross_auth_key {
			c.JSON(404, gin.H{
				"error": true,
				"msg":   "auth failed!",
			})
			return
		}

		insert, err := db.Query("INSERT INTO users(username, phone) VALUES ( '" + username + "', '" + phone + "' );")

		// if there is an error inserting, handle it
		if err != nil {
			c.JSON(200, gin.H{
				"error": true,
				"msg":   "insert failed!",
			})
			return
		}
		// be careful deferring Queries if you are using transactions
		defer insert.Close()

		type User struct {
			Id int `json:"id"`
		}

		var user User
		// Execute the query
		err = db.QueryRow("SELECT id FROM users where phone = ?", phone).Scan(&user.Id)
		if err != nil {
			c.JSON(200, gin.H{
				"error": true,
				"msg":   "query failed",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"id":    user.Id,
		})
	})

	r.GET("/getNicknameById", func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			c.JSON(200, gin.H{
				"error": true,
				"msg":   "parameter failed",
			})
			return
		}

		type User struct {
			Id       int    `json:"id"`
			Username string `json:"username"`
		}

		var user User
		// Execute the query
		err = db.QueryRow("SELECT id, username FROM users where id = ?", id).Scan(&user.Id, &user.Username)
		fmt.Println(err)
		if err != nil {
			c.JSON(200, gin.H{
				"error": true,
				"msg":   "query user failed!",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":    false,
			"id":       user.Id,
			"username": user.Username,
		})

	})

	r.Run(":80")
}
