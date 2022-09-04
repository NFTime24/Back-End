package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/duke/db"
	"github.com/duke/model"
	"github.com/labstack/echo/v4"
)

var KlipRequestMap map[uint64]string

type KlipResponse struct {
	RequestKey     string     `json:"request_key"`
	Status         string     `json:"status"`
	Result         KlipResult `json:"result"`
	ExpirationTime int        `json:"expiration_time"`
	RequestURL     string     `json:"request_url"`
}

type KlipResult struct {
	KlaytnAddress string `json:"klaytn_address"`
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
		NftID:   uint(newItemId),
		WorksID: uint(work_id),
		OwnerID: 0,
		// ExibitionsID: 1,
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

func AddNFTWithWorkId(c echo.Context) error {
	work_id_str := c.QueryParam("work_id")
	work_id, err := strconv.ParseUint(work_id_str, 10, 64)
	if err != nil {
		fmt.Printf(err.Error())
	}

	db := db.DbManager()
	var nfts model.Nft
	var result int

	db.Model(nfts).Select(`MAX(nft_id)`).Scan(&result)

	newItemId := result + 1

	db.Create(model.Nft{
		NftID:   uint(newItemId),
		WorksID: uint(work_id),
		OwnerID: 0,
	})

	resultStr := strconv.Itoa(newItemId)
	return c.String(http.StatusOK, resultStr)
}

func MintArtWithoutPaying(c echo.Context) error {
	work_id_str := c.QueryParam("work_id")
	work_id, err := strconv.ParseUint(work_id_str, 10, 64)
	if err != nil {
		fmt.Printf(err.Error())
	}

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
		NftID:   uint(newItemId),
		WorksID: uint(work_id),
		OwnerID: 0,
	})

	klipKey := rand.Uint64()
	reqBodyStr := fmt.Sprintf(`{
		"type": "auth",
		"bapp": {
			"name" : "NFTime",
			"callback": { "success": "http:\/\/34.212.84.161\/onSuccessKlip?klip_key=%s", "fail": "" }
		}
	}`, strconv.FormatUint(klipKey, 10))
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
	KlipRequestMap[klipKey] = jData.RequestKey
	jData.RequestURL = "https://klipwallet.com/?target=/a2a?request_key="
	jData.RequestURL += jData.RequestKey

	fmt.Println(jData.RequestURL)

	// http.Redirect(w, r, jData.RequestQR, http.StatusFound)
	c.Redirect(http.StatusFound, jData.RequestURL)
	return nil
	//t.Execute(w, jData)
}

func OnSuccessKlip(c echo.Context) error {
	klipKey_str := c.QueryParam("klip_key")
	klipKey, err := strconv.ParseUint(klipKey_str, 10, 64)
	if err != nil {
		fmt.Printf(err.Error())
	}

	requestKey := KlipRequestMap[klipKey]

	fmt.Println(requestKey)

	client := &http.Client{}
	reqStr := fmt.Sprintf("https://a2a-api.klipwallet.com/v2/a2a/result?request_key=%s", requestKey)
	req, err := http.NewRequest("GET", reqStr, nil)
	if err != nil {
		fmt.Printf(err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
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

	fmt.Printf("Klaytn address: %s", jData.Result.KlaytnAddress)

	kasClient := &http.Client{}
	kasReqStr := fmt.Sprintf("https://wallet-api.klaytnapi.com/v2/tx/contract/execute")
	jsonStr := fmt.Sprintf(`{
		"from": "0x7c07C1579aD1980863c83876EC4bec43BC8d6dFa",
		"value": "0x0",
		"to": "0xeb0912eff03e357c4cbb9c9c925ae01b2da1e486",
		"input": "0x697d04130000000000000000000000007c07c1579ad1980863c83876ec4bec43bc8d6dfa0000000000000000000000000000000000000000000000000000000000000bb9000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000047465737400000000000000000000000000000000000000000000000000000000",
		"nonce": 0,
		"gasLimit": 1000000,
		"submit": true
	}`)
	kasReq, err := http.NewRequest("POST", kasReqStr, bytes.NewBufferString(jsonStr))
	if err != nil {
		fmt.Printf(err.Error())
	}
	kasReq.Header.Add("x-chain-id", "8217")
	kasReq.Header.Add("Content-Type", "application/json")
	kasReq.Header.Add("Authorization", "Basic S0FTS0NDRjIxR1VZUUdCOE83Q0JQR09GOm1waHN0cTllSDFTV1d6cXNFX3JrTEM0LTRCMDVFYWhyWmg5SVNFbWI=")
	kasResp, err := kasClient.Do(req)
	if err != nil {
		fmt.Printf(err.Error())
	}
	defer kasResp.Body.Close()
	kasBody, err := io.ReadAll(kasResp.Body)
	if err != nil {
		fmt.Printf(err.Error())
	}
	fmt.Printf("kas body: %s \n", kasBody)

	return c.String(http.StatusOK, requestKey)
}
