package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func getPage(uri string) []byte {
	apiKey := "55251e9d20b4ad8e7e625a406d717b25"
	apiSecret := "004ccbdf594fedec1d547ff75219f6335242f9ed"
	nonce := fmt.Sprint(time.Now().Unix() % 1000000000)
	uri = uri + apiKey + "&nonce=" + nonce
	hash := hmac.New(sha512.New, []byte(apiSecret))
	hash.Write([]byte(uri))
	sign := hex.EncodeToString(hash.Sum(nil))
	uri = uri + "&apisign=" + sign
	fmt.Println(uri)
	resp, _ := http.Get(uri)
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

func balance(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	uri := "https://bleutrade.com/api/v2/account/getbalance?"
	currency := queryValues.Get("currency")
	uri += "currency=" + currency
	text := getPage(uri + "&apikey=")
	w.Write([]byte(text))
}

func withdraw(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	currency := queryValues.Get("currency")
	quantity := queryValues.Get("quantity")
	address := "DTBpf4fLrpzSjcYWq6tBqqczznzjHhphgX"
	comments := "test2"
	uri := "https://bleutrade.com/api/v2/account/withdraw?"
	uri += "currency=" + currency
	uri += "&quantity=" + quantity
	uri += "&address=" + address
	uri += "&comments=" + comments
	// text := getPage(uri + "&apikey=")
	// w.Write([]byte(text))
}

func main() {
	http.HandleFunc("/balance", balance)
	http.HandleFunc("/withdraw", withdraw)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	// http.HandleFunc("/hello", form)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
