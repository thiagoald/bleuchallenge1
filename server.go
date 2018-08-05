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

var etherscanKey = "YCBTHQ1KTEUMYTMYX39YY3U36GSZ3Q2FRP"

func genHash(apiSecret, uri string) (sign string) {
	hash := hmac.New(sha512.New, []byte(apiSecret))
	hash.Write([]byte(uri))
	sign = hex.EncodeToString(hash.Sum(nil))
	return sign
}

func balance(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	nonce := fmt.Sprint(time.Now().Unix() % 1000000000)
	apikey := queryValues.Get("apikey")
	currency := queryValues.Get("currency")
	apisecret := queryValues.Get("apisecret")
	uri := "https://bleutrade.com/api/v2/account/getbalance?apikey=" + apikey + "&nonce=" + nonce
	uri += "&currency=" + currency
	sign := genHash(apisecret, uri)
	uri += "&apisign=" + sign
	resp, _ := http.Get(uri)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(uri)
	w.Write([]byte(body))
}

func balances(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	nonce := fmt.Sprint(time.Now().Unix() % 1000000000)
	apikey := queryValues.Get("apikey")
	apisecret := queryValues.Get("apisecret")
	uri := "https://bleutrade.com/api/v2/account/getbalances?apikey=" + apikey + "&nonce=" + nonce
	sign := genHash(apisecret, uri)
	uri += "&apisign=" + sign
	resp, _ := http.Get(uri)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(uri)
	w.Write([]byte(body))
}

func withdraw(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	currency := queryValues.Get("currency")
	quantity := queryValues.Get("quantity")
	address := queryValues.Get("address")
	apikey := queryValues.Get("apikey")
	apisecret := queryValues.Get("apisecret")
	nonce := fmt.Sprint(time.Now().Unix() % 1000000000)
	uri := "https://bleutrade.com/api/v2/account/withdraw?apikey=" + apikey + "&nonce=" + nonce

	uri += "&currency=" + currency
	uri += "&quantity=" + quantity
	uri += "&address=" + address
	sign := genHash(apisecret, uri)
	uri += "&apisign=" + sign
	fmt.Println(uri)
	resp, _ := http.Get(uri)
	body, _ := ioutil.ReadAll(resp.Body)
	w.Write([]byte(body))
}

func main() {
	http.HandleFunc("/balance", balance)
	http.HandleFunc("/balances", balances)
	http.HandleFunc("/withdraw", withdraw)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
