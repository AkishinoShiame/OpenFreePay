package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

type ServerConfig struct {
	UserServer   string `json:"user_server"`
	PointsServer string `json:"points_server"`
}

var serverconfig = ServerConfig{}

func init() {
	rand.Seed(time.Now().UnixNano())

	resp, _ := http.Get("https://lab.gris.tw/")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	str := string(body)
	json.Unmarshal([]byte(str), &serverconfig)

	fmt.Println("UserServer: " + serverconfig.UserServer)
	fmt.Println("PointsServer: " + serverconfig.PointsServer)
}

var cross_auth_key = "rup j03t s.u6q96"

//需要注意這裡的 status 是字串的 false
//status == 0 表示成功
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

	r.GET("/listMarketOrder", func(c *gin.Context) {
		json, err := getJSONfromQuery("SELECT * FROM market WHERE status IS NOT true;", db)
		if err != nil {
			c.JSON(200, gin.H{
				"error": true,
				"msg":   err,
			})
			return
		}

		c.String(200, json)
	})

	r.POST("/createMarketOrder", func(c *gin.Context) {
		src_point_type := c.PostForm("src_point_type")
		dest_point_type := c.PostForm("dest_point_type")
		src_bid_points := c.PostForm("src_bid_points")
		//dest_ask_points := c.PostForm("dest_ask_points")
		user_id := c.PostForm("user_id")

		exchange_ask := httpGet(string(serverconfig.PointsServer) + "/" + "calcPointExchange?src=" + src_point_type + "&dst=" + dest_point_type + "&quantity=" + src_bid_points)
		dest_ask_points := string(exchange_ask)

		resp := httpGetJSON(string(serverconfig.UserServer) + "/" + "getNicknameById?id=" + user_id)
		username := resp["username"].(string)

		insert, err := db.Query("INSERT INTO market( user_id, username, src_point_type, dest_point_type, src_bid_points, dest_ask_points) VALUES ( " + user_id + ", '" + username + "', '" + src_point_type + "', '" + dest_point_type + "', '" + src_bid_points + "', '" + dest_ask_points + "'  )")

		// if there is an error inserting, handle it
		if err != nil {
			c.JSON(200, gin.H{
				"error": true,
				"msg":   err,
			})
			return
		}
		// be careful deferring Queries if you are using transactions
		defer insert.Close()

		c.String(200, "Redirecting...")
		c.Redirect(http.StatusMovedPermanently, "./")
	})

	r.Run(":80")
}

func getJSONfromQuery(sqlString string, db *sql.DB) (string, error) {
	rows, err := db.Query(sqlString)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return "", err
	}
	fmt.Println(string(jsonData))
	return string(jsonData), nil
}

func httpGet(url string) string {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func httpGetJSON(url string) map[string]interface{} {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)
	return result
}

func httpPost(url string, formData url.Values) map[string]interface{} {
	/*formData := url.Values{
		"name": {"masnun"},
	}*/

	var result map[string]interface{}

	resp, err := http.PostForm(url, formData)
	if err != nil {
		fmt.Println(err)
		result = map[string]interface{}{"error": true, "msg": "fetch failed!"}
		return result
	}

	json.NewDecoder(resp.Body).Decode(&result)

	//log.Println(result["form"])
	return result
}
