package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"
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

func main() {
	//fmt.Println(string("stirng"[1:]) == nil)

	body := `{"request_id":"55f6f92b61aa48e4a171ff3e88dddd6c","status":"0","event_id":"1B00000036475D1D","price":"0.05000000","currency":"EUR","estimated_price_messages_sent":"0.04420000"}`

	type verifyObejct struct {
		Request_id string `json:"request_id"`
		Status     string `json:"status"`
		Event_id   string `json:"event_id"`
		Price      string `json:"price"`
		Currency   string `json:"currency"`
	}

	//str := `{"page": 1, "fruits": ["apple", "peach"]}`

	m := verifyObejct{}
	str := string(body)
	json.Unmarshal([]byte(str), &m)
	fmt.Println(m.Status == "0")

	fmt.Println(string(serverconfig.UserServer))
	resp := httpPost(string(serverconfig.UserServer)+"/"+"createUser", url.Values{
		"auth":     {"rup j03t s.u6q96"},
		"username": {"goodsdasdsadas"},
		"phone":    {"0000000000"},
	})

	fmt.Println(resp["error"])
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
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
