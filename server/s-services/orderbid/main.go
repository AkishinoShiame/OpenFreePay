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

	r.GET("/directlyOrder", func(c *gin.Context) {
		order_id := c.Query("order_id")
		user_id := c.Query("user_id")

		resp := httpGetJSON(string(serverconfig.UserServer) + "/" + "getNicknameById?id=" + user_id)
		username := resp["username"].(string)

		_, err = db.Query("UPDATE market SET `status` = true WHERE order_id = " + order_id)
		if err != nil {
			c.JSON(200, gin.H{
				"error": true,
				"msg2":  err,
			})
			return
		}

		type Market struct {
			Order_id        int     `json:"order_id"`
			Date            string  `json:"date"`
			User_id         string  `json:"user_id"`
			Src_point_type  string  `json:"src_point_type"`
			Dest_point_type string  `json:"dest_point_type"`
			Revision_times  string  `json:"revision_times"`
			Src_bid_points  float64 `json:"src_bid_points"`
			Dest_ask_points float64 `json:"dest_ask_points"`
			Username        string  `json:"username"`
			Status          bool    `json:"status"`
		}

		var market Market
		// Execute the query
		err = db.QueryRow("SELECT * FROM market where order_id = ?", order_id).Scan(&market.Order_id, &market.Date, &market.User_id, &market.Src_point_type, &market.Dest_point_type, &market.Revision_times, &market.Src_bid_points, &market.Dest_ask_points, &market.Username, &market.Status)
		fmt.Println(err)
		if err != nil {
			c.JSON(200, gin.H{
				"error": true,
				"msg":   "query market failed!",
			})
			return
		}

		_, err := db.Query("INSERT INTO orderbid(order_id, user_id, username, target_src_bid_points, target_dest_ask_points, transcation_state) VALUES( " + order_id + ", " + user_id + ", '" + username + "' ,  " + fmt.Sprintf("%.2f", market.Src_bid_points) + ", " + fmt.Sprintf("%.2f", market.Dest_ask_points) + ", true);")
		if err != nil {
			c.JSON(200, gin.H{
				"error": true,
				"msg1":  err,
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
		})
	})

	r.GET("/bidTheOrder", func(c *gin.Context) {
		src_point_type := c.Query("src_point_type")
		dest_point_type := c.Query("dest_point_type")
		order_id := c.Query("order_id")
		target_dest_ask_points := c.Query("target_dest_ask_points")
		user_id := c.Query("user_id")

		//交換 src, dest
		fmt.Println(string(serverconfig.PointsServer) + "/" + "calcPointExchange?src=" + dest_point_type + "&dst=" + src_point_type + "&quantity=" + target_dest_ask_points)
		exchange_ask := httpGet(string(serverconfig.PointsServer) + "/" + "calcPointExchange?src=" + dest_point_type + "&dst=" + src_point_type + "&quantity=" + target_dest_ask_points)
		target_src_bid_points := string(exchange_ask)

		resp := httpGetJSON(string(serverconfig.UserServer) + "/" + "getNicknameById?id=" + user_id)
		username := resp["username"].(string)

		_, err := db.Query("INSERT INTO orderbid(order_id, target_src_bid_points, target_dest_ask_points, user_id, username) VALUES( " + order_id + ",  " + target_src_bid_points + ", " + target_dest_ask_points + ", " + user_id + ", '" + username + "');")
		if err != nil {
			c.JSON(200, gin.H{
				"error": true,
				"msg1":  err,
			})
			return
		}

		_, err = db.Query("UPDATE market SET `revision_times` = revision_times + 1 WHERE order_id = " + order_id)
		if err != nil {
			c.JSON(200, gin.H{
				"error": true,
				"msg2":  err,
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
		})
	})

	//我的訂單 main-logical
	r.GET("/listMyOrder", func(c *gin.Context) {
		user_id := c.Query("user_id")

		json, err := getJSONfromQuery("SELECT * FROM market WHERE user_id = '"+user_id+"';", db)
		if err != nil {
			c.JSON(200, gin.H{
				"error": true,
				"msg":   err,
			})
			return
		}



		jsonparser.ArrayEach([]byte(json), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			
			fmt.Println(jsonparser.Get(value, ""))

		}, "person", "avatars")

		c.String(200, json)
	})

	//我的訂單 sub-logical
	r.GET("/listMyOrderAndBidDetails", func(c *gin.Context) {
		order_id := c.Query("order_id")

		json, err := getJSONfromQuery("SELECT * FROM orderbid WHERE order_id = "+order_id+";", db)
		if err != nil {
			c.JSON(200, gin.H{
				"error": true,
				"msg":   err,
			})
			return
		}

		c.String(200, json)
	})

	//向調價者進行交易 (deal with it and exits user_id problem)
	r.GET("/dealWithIt", func(c *gin.Context) {
		c_id := c.Query("c_id")
		order_id := c.Query("order_id")

		_, err = db.Query("UPDATE market SET `status` = true WHERE order_id = " + order_id)
		if err != nil {
			c.JSON(200, gin.H{
				"error": true,
				"msg2":  err,
			})
			return
		}

		_, err = db.Query("UPDATE orderbid SET transcation_state = true WHERE c_id = " + c_id)
		if err != nil {
			c.JSON(200, gin.H{
				"error": true,
				"msg2":  err,
			})
			return
		}

		c.String(200, "Redirecting...")
		c.Redirect(304, "./")
	})

	//我的出價
	r.GET("/listMyBid", func(c *gin.Context) {
		user_id := c.Query("user_id")

		//SELECT orderbid.*, market.* FROM orderbid INNER JOIN market ON market.order_id=orderbid.order_id WHERE orderbid.user_id = 14;

		//json, err := getJSONfromQuery("SELECT * FROM orderbid WHERE user_id = "+user_id, db)
		json, err := getJSONfromQuery("SELECT orderbid.*, market.* FROM orderbid INNER JOIN market ON market.order_id=orderbid.order_id WHERE orderbid.user_id ="+user_id, db)
		if err != nil {
			c.JSON(200, gin.H{
				"error": true,
				"msg":   err,
			})
			return
		}

		c.String(200, json)
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
