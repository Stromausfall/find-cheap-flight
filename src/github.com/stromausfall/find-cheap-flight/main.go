package main

import (
	"net/http"
	"github.com/stromausfall/find-cheap-flight/data"
	"github.com/stromausfall/find-cheap-flight/utils"
)

func main() {
	googleMapsApiCredentials := "google maps API key"
	
	// install handlers
	go http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// we need this in order to get POST form data
		r.ParseMultipartForm(15485760)
	
		data.DisplayDataEntry(w, r, googleMapsApiCredentials)
	})
	go http.HandleFunc("/selection", func(w http.ResponseWriter, r *http.Request) {
		// we need this in order to get POST form data
		r.ParseMultipartForm(15485760)
		
		data.DisplayDataSelection(w, r, googleMapsApiCredentials)
	})

	utils.OpenURL("http://localhost:80")

	// start web server
	http.ListenAndServe(":80", nil)
}