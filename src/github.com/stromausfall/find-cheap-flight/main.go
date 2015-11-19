package main

import (
	"net/http"
	"github.com/stromausfall/find-cheap-flight/data"
	"github.com/stromausfall/find-cheap-flight/utils"
)

func main() {
	googleMapsApiCredentials := "asdf"
	genomeAccount := "asdf"
	
	// install handlers
	go http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {	
		data.DisplayDataEntry(w, r, googleMapsApiCredentials)
	})
	go http.HandleFunc("/selection", func(w http.ResponseWriter, r *http.Request) {
		data.DisplayDataSelection(w, r, googleMapsApiCredentials, genomeAccount)
	})

	utils.OpenURL("http://localhost:80")

	// start web server
	http.ListenAndServe(":80", nil)
}
