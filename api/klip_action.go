package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/duke/db"
	"github.com/duke/model"
	"github.com/labstack/echo/v4"
	"github.com/umbracle/ethgo/abi"
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

	klipKey := rand.Uint64()
	reqBodyStr := fmt.Sprintf(`{
		"type": "auth",
		"bapp": {
			"name" : "NFTime",
			"callback": { "success": "http:\/\/34.212.84.161\/onSuccessKlip?klip_key=%s&work_id=%s", "fail": "" }
		}
	}`, strconv.FormatUint(klipKey, 10), strconv.FormatUint(work_id, 10))
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

	workId_str := c.QueryParam("work_id")
	workId, err := strconv.ParseUint(workId_str, 10, 64)
	if err != nil {
		fmt.Printf(err.Error())
	}

	requestKey := KlipRequestMap[klipKey]

	fmt.Println(requestKey)

	db := db.DbManager()
	var nfts model.Nft
	var result uint64

	db.Model(nfts).Select(`MAX(nft_id)`).Scan(&result)

	newItemId := result + 1

	db.Create(model.Nft{
		NftID:   uint(newItemId),
		WorksID: uint(workId),
		OwnerID: 0,
	})

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
	fmt.Println("")

	typ := abi.MustNewType("uint256")
	nftId_big := big.NewInt(int64(newItemId))
	encoded, err := typ.Encode(nftId_big)
	if err != nil {
		panic(err)
	}

	nftId_hex := fmt.Sprintf("%x", encoded)
	addressBase := "0000000000000000000000000000000000000000000000000000000000000000"
	ablen := len(addressBase)
	kalen := len(jData.Result.KlaytnAddress)
	addr_hex := fmt.Sprintf("%s%s", addressBase[:(ablen-kalen+2)], jData.Result.KlaytnAddress[2:])

	reqCallData := "0x697d0413"
	reqCallData += addr_hex
	reqCallData += nftId_hex
	reqCallData += "0000000000000000000000000000000000000000000000000000000000000060"
	reqCallData += "0000000000000000000000000000000000000000000000000000000000000004"
	reqCallData += "7465737400000000000000000000000000000000000000000000000000000000"

	fmt.Println(reqCallData)

	kasClient := &http.Client{}
	kasReqStr := fmt.Sprintf("https://wallet-api.klaytnapi.com/v2/tx/contract/execute")
	jsonStr := fmt.Sprintf(`{
		"from": "0x7c07C1579aD1980863c83876EC4bec43BC8d6dFa",
		"value": "0x0",
		"to": "0xeb0912eff03e357c4cbb9c9c925ae01b2da1e486",
		"input": "%s",
		"nonce": 0,
		"gasLimit": 1000000,
		"submit": true
	}`, reqCallData)
	kasReq, err := http.NewRequest("POST", kasReqStr, bytes.NewBufferString(jsonStr))
	if err != nil {
		fmt.Printf(err.Error())
	}
	kasReq.Header.Add("x-chain-id", "8217")
	kasReq.Header.Add("Content-Type", "application/json")
	kasReq.Header.Add("Authorization", "Basic S0FTS0NDRjIxR1VZUUdCOE83Q0JQR09GOm1waHN0cTllSDFTV1d6cXNFX3JrTEM0LTRCMDVFYWhyWmg5SVNFbWI=")
	kasResp, err := kasClient.Do(kasReq)
	if err != nil {
		fmt.Printf(err.Error())
	}
	defer kasResp.Body.Close()
	kasBody, err := io.ReadAll(kasResp.Body)
	if err != nil {
		fmt.Printf(err.Error())
	}
	fmt.Printf("kas body: %s \n", kasBody)

	return c.Redirect(http.StatusFound, "https://klipwallet.com/")
}
