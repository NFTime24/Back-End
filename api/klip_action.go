package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo"
)

type KlipResponse struct {
	RequestKey     string `json:"request_key"`
	Status         string `json:"status"`
	ExpirationTime int    `json:"expiration_time"`
	RequestURL     string `json:"request_url"`
}

func PrepareAuth(c echo.Context) error {

	reqBody := bytes.NewBufferString(`{"bapp": { "name" : "My BApp" }, "type": "execute_contract", "transaction": { "to": "0xFf1C1e55826DD95C44681BfCd88DCB32eE86B793", "value": "0", "abi": "{\"inputs\": [{\"internalType\": \"string\",\"name\": \"artist_address\",\"type\": \"string\"}],\"name\": \"mintArt\",\"outputs\": [{\"internalType\": \"uint256\",\"name\": \"\",\"type\": \"uint256\"}],\"stateMutability\": \"nonpayable\",\"type\": \"function\"}", "params": "[\"0xF39E4961C046BA913f835c08Bf25De348184F3a8\"]"} }`)
	resp, err := http.Post("https://a2a-api.klipwallet.com/v2/a2a/prepare", "Content-Type: application/json", reqBody)
	if err != nil {
		fmt.Printf(err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf(err.Error())
	}
	var jData KlipResponse
	fmt.Printf("body: %s \n", body)
	json.Unmarshal(body, &jData)
	// t, err := template.ParseGlob("./templates/qr.html")
	// if err != nil {
	// 	fmt.Printf(err.Error())
	// }
	fmt.Printf("requestkey: %s \n", jData.RequestKey)
	jData.RequestURL = "intent://klipwallet/open?url=https://klipwallet.com/?target=/a2a?request_key="
	jData.RequestURL += jData.RequestKey
	jData.RequestURL += "#Intent;scheme=kakaotalk;package=com.kakao.talk;end"
	fmt.Printf(jData.RequestURL)

	// http.Redirect(w, r, jData.RequestQR, http.StatusFound)
	c.Redirect(http.StatusFound, jData.RequestURL)
	return nil
	//t.Execute(w, jData)
}