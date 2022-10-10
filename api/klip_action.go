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

	workId := c.QueryParam("work_id")

	requestKey := KlipRequestMap[klipKey]

	fmt.Println(requestKey)

	client1 := &http.Client{}
	reqStr1 := fmt.Sprintf("https://a2a-api.klipwallet.com/v2/a2a/result?request_key=%s", requestKey)
	req1, err := http.NewRequest("GET", reqStr1, nil)
	if err != nil {
		fmt.Printf(err.Error())
	}
	req1.Header.Add("Content-Type", "application/json")
	resp1, err := client1.Do(req1)
	if err != nil {
		fmt.Printf(err.Error())
	}
	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		fmt.Printf(err.Error())
	}
	var jData1 KlipResponse
	fmt.Printf("body1: %s \n", body1)
	json.Unmarshal(body1, &jData1)

	fmt.Printf("Klaytn address: %s", jData1.Result.KlaytnAddress)
	fmt.Println("")
	resp1.Body.Close()

	client2 := &http.Client{}
	reqStr2 := fmt.Sprintf("http://34.212.84.161/mintToAddr?address=%s&work_id=%s", jData1.Result.KlaytnAddress, workId)
	req2, err := http.NewRequest("GET", reqStr2, nil)
	if err != nil {
		fmt.Printf(err.Error())
	}
	resp2, err := client2.Do(req2)
	if err != nil {
		fmt.Printf(err.Error())
	}
	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		fmt.Printf(err.Error())
	}
	fmt.Printf("body2: %s \n", body2)

	resString := string(body2[:])

	fmt.Printf("result: %s", resString)
	fmt.Println("")
	resp2.Body.Close()

	return c.String(http.StatusOK, resString)
}

func MintToAddr(c echo.Context) error {
	address := c.QueryParam("address")

	workId_str := c.QueryParam("work_id")
	workId, err := strconv.ParseUint(workId_str, 10, 64)
	if err != nil {
		fmt.Printf(err.Error())
	}

	fmt.Printf("Work Id: %d", workId)
	fmt.Println("")

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

	fmt.Printf("Klaytn address: %s", address)
	fmt.Println("")

	typ := abi.MustNewType("uint256")

	nftId_big := big.NewInt(int64(newItemId))
	nftId_encoded, err := typ.Encode(nftId_big)
	if err != nil {
		panic(err)
	}
	nftId_hex := fmt.Sprintf("%x", nftId_encoded)

	workId_big := big.NewInt(int64(workId))
	workId_encoded, err := typ.Encode(workId_big)
	if err != nil {
		panic(err)
	}
	workId_hex := fmt.Sprintf("%x", workId_encoded)

	addressBase := "0000000000000000000000000000000000000000000000000000000000000000"
	ablen := len(addressBase)
	kalen := len(address)
	addr_hex := fmt.Sprintf("%s%s", addressBase[:(ablen-kalen+2)], address[2:])

	reqCallData := "0x20b7668b"
	reqCallData += addr_hex
	reqCallData += nftId_hex
	reqCallData += workId_hex

	fmt.Println(reqCallData)

	kasClient := &http.Client{}
	kasReqStr := fmt.Sprintf("https://wallet-api.klaytnapi.com/v2/tx/contract/execute")
	jsonStr := fmt.Sprintf(`{
		"from": "0x7c07C1579aD1980863c83876EC4bec43BC8d6dFa",
		"value": "0x0",
		"to": "0x63Ff714D28D84Eb336cEd92b9DAB59C8797D5bCB",
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

	return c.String(http.StatusOK, "Success Mint")
}
