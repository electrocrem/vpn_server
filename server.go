package main

import (
	"log"
	"net/http"

	"github.com/electrocrem/vpn_server/cmd/endpoints"
)

func main() {

	mux := http.NewServeMux()
	mux.Handle("/", endpoints.Page("./static"))
	mux.HandleFunc("/profile/", endpoints.GetProfile)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
