package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	UserServer string `json:"user_server"`
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
}

var cross_auth_key = "rup j03t s.u6q96"

//需要注意這裡的 status 是字串的 false
//status == 0 表示成功
func main() {

	r := gin.Default()

	r.POST("/requestSMS", func(c *gin.Context) {
		phone := c.PostForm("phone")
		if string(phone[0]) != "0" {
			c.JSON(200, gin.H{
				"status": "false",
			})
			return
		}

		resp, _ := http.Get("https://api.nexmo.com/verify/json?api_key=xxxxx&api_secret=xxxxx&number=886" + string(phone[1:]) + "&brand=Nexmo&code_length=4")

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		// ...
		c.String(200, string(body))
	})

	r.POST("/verify", func(c *gin.Context) {

		request_id := c.PostForm("request_id")
		verify_code := c.PostForm("verify_code")
		phone := c.PostForm("phone") //security

		if string(request_id) == "" || string(verify_code) == "" || string(phone) == "" {
			log.Println("Empty Form!!")
			c.JSON(200, gin.H{
				"status": "false",
			})
			return
		}

		resp, _ := http.Get("https://api.nexmo.com/verify/check/json?api_key=xxxx&api_secret=xxxx&request_id=" + string(request_id) + "&code=" + string(verify_code))

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		type verifyObejct struct {
			Request_id                    string `json:"request_id"`
			Status                        string `json:"status"`
			Event_id                      string `json:"event_id"`
			Price                         string `json:"price"`
			Currency                      string `json:"currency"`
			Estimated_price_messages_sent string `json:"estimated_price_messages_sent"`
		}

		m := verifyObejct{}
		str := string(body)
		json.Unmarshal([]byte(str), &m)

		if m.Status == "0" {

			/*verify:
			obj := httpGetJSON("https://api.nexmo.com/verify/search/json?api_key=xxxx&api_secret=xxxx&request_id=" + request_id)
			if reflect.ValueOf(obj).IsNil() {
				time.Sleep(time.Duration(1) * time.Second)
				goto verify
			}
			phone := "0" + obj["number"].(string)[3:]*/

			nickname := getNickName()
			//使用者註冊成功
			resp := httpPost(string(serverconfig.UserServer)+"/"+"createUser", url.Values{
				"auth":     {cross_auth_key},
				"username": {nickname},
				"phone":    {phone},
			})

			c.JSON(200, gin.H{
				"status":  "0",
				"user_id": resp["id"],
			})

			return
		}

		c.JSON(200, gin.H{
			"status": "false",
		})
	})

	r.Run(":80")
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
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

func getNickName() string {

	var re = regexp.MustCompile(`(?ms)的稱號就是：\s\n\n\n(.*?)<\/div>`)

	rand := RandStringRunes(3)
	fmt.Println(rand)
	resp, _ := http.Get("https://wtf.hiigara.net/t/titlegen/" + rand)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	name := re.FindAllStringSubmatch(string(body), -1)[0][1]

	return name
}
