package httpHandlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

type KlipResponse struct {
	RequestKey     string `json:"request_key"`
	RequestQR      string `json:"request_qr"`
	Status         string `json:"status"`
	ExpirationTime int    `json:"expiration_time"`
}

func PrepareAuth(w http.ResponseWriter, r *http.Request) {
	reqBody := bytes.NewBufferString(`{"bapp":{"name" : "My BApp"}, "callback": { "success": "", "fail": "" }, "type": "auth"}`)
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
	json.Unmarshal(body, &jData)
	t, err := template.ParseGlob("./templates/qr.html")
	if err != nil {
		fmt.Printf(err.Error())
	}
	jData.RequestQR = "intent://klipwallet/open?url=https://klipwallet.com/?target=/a2a?request_key="
	jData.RequestQR += "0b0ee0ad-62b3-4146-980b-531b3201265d"
	jData.RequestQR += "#Intent;scheme=kakaotalk;package=com.kakao.talk;end"
	fmt.Printf(jData.RequestQR)

	t.Execute(w, jData)
	//http.Redirect(w, r, jData.RequestQR, http.StatusFound)
}
