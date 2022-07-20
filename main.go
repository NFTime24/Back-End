package main

import "net/http"

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("hello"))
	})
	http.ListenAndServe(":80", nil)
}
