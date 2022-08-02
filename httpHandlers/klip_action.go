package httpHandlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

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
	w.Write(body)
}
