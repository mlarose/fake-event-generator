package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/event", postEventHandler)

	log.Fatal(http.ListenAndServe(":3333", nil))
}

func postEventHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println(string(body))
	w.WriteHeader(200)
}
