package webserver

import (
	//	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var MessageMap = make(map[string][]byte)

func orders(w http.ResponseWriter, r *http.Request) {

	log.Printf("Orders update request from %v, User-Agent: %v\n", r.Host, r.UserAgent())
	jsonString, err := json.Marshal(MessageMap)
	if err != nil {
		log.Printf("Cannot marshal map: %v\n", err)
	}

	fmt.Fprintf(w, "%s\n", jsonString)
}

func StartServer(port string) {
	http.HandleFunc("/orders", orders)
	http.ListenAndServe(":"+port, nil)
}
