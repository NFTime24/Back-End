package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/duke/db"
	"github.com/duke/model"
	"github.com/labstack/echo/v4"
)

type KlipResponse struct {
	RequestKey     string `json:"request_key"`
	Status         string `json:"status"`
	ExpirationTime int    `json:"expiration_time"`
	RequestURL     string `json:"request_url"`
}

func MintArt(c echo.Context) error {
	work_id_str := c.QueryParam("work_id")
	work_id, err := strconv.ParseUint(work_id_str, 10, 64)
	if err != nil {
		fmt.Printf(err.Error())
	}

	osType := c.QueryParam("os_type") // "ios" or "aos"

	db := db.DbManager()
	var nfts model.Nft
	var result1 uint64

	db.Model(nfts).Select(`MAX(nft_id)`).Scan(&result1)

	var works model.Work
	var result2 string

	db.Model(works).Select(`a.address`).
		Joins("left join test.artists as a on a.id = works.artist_id").
		Where("works.work_id=?", work_id).
		Scan(&result2)

	newItemId := result1 + 1
	artist_address := result2

	fmt.Printf("\n newItemId: %d, artist_address: %s \n", newItemId, artist_address)

	db.Create(model.Nft{
		NftID:        uint(newItemId),
		WorksID:      uint(work_id),
		OwnerID:      0,
		ExibitionsID: 1,
	})

	reqBodyStr := fmt.Sprintf(`{
		"bapp": { "name" : "NFTime" }, 
		"type": "execute_contract", 
		"transaction": { 
			"to": "0xf1cB5DDF7E8E9Af429b79473c41Dd85750Faa7af", 
			"value": "0", 
			"abi": "{ \"inputs\": [ { \"internalType\": \"uint256\", \"name\": \"newItemId\", \"type\": \"uint256\" }, { \"internalType\": \"string\", \"name\": \"artist_address\", \"type\": \"string\" } ], \"name\": \"mintArt\", \"outputs\": [ { \"internalType\": \"uint256\", \"name\": \"\", \"type\": \"uint256\" } ], \"stateMutability\": \"nonpayable\", \"type\": \"function\" }", 
			"params": "[%d, \"%s\"]"
		} 
	}`, newItemId, artist_address)

	reqBody := bytes.NewBufferString(reqBodyStr)
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
	if osType == "ios" {
		jData.RequestURL = "kakaotalk://klipwallet/open?url=https://klipwallet.com/?target=/a2a?request_key="
		jData.RequestURL += jData.RequestKey
	} else {
		jData.RequestURL = "intent://klipwallet/open?url=https://klipwallet.com/?target=/a2a?request_key="
		jData.RequestURL += jData.RequestKey
		jData.RequestURL += "#Intent;scheme=kakaotalk;package=com.kakao.talk;end"
	}

	fmt.Printf(jData.RequestURL)

	// http.Redirect(w, r, jData.RequestQR, http.StatusFound)
	c.Redirect(http.StatusFound, jData.RequestURL)
	return nil
	//t.Execute(w, jData)
}
