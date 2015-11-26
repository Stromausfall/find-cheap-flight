package main

import (
	"github.com/stromausfall/find-cheap-flight/data"
	"github.com/stromausfall/find-cheap-flight/utils"
	"net/http"
)

func main() {
	googleMapsApiCredentials := "GOOGLE MAPS API CREDENTIALS"
	geonameAccount := "GEONAME ACCOUNT"  

	// install handlers
	go http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data.DisplayDataEntry(w, r, googleMapsApiCredentials)
	})
	go http.HandleFunc("/selectAirports", func(w http.ResponseWriter, r *http.Request) {
		data.DisplayDataSelection(w, r, googleMapsApiCredentials, geonameAccount)
	})
	go http.HandleFunc("/dataVerification", func(w http.ResponseWriter, r *http.Request) {
		data.DisplayDataVerification(w, r, googleMapsApiCredentials, geonameAccount)
	})
	go http.HandleFunc("/calculate", func(w http.ResponseWriter, r *http.Request) {
	})

	utils.OpenURL("http://localhost:80")

	// start web server
	http.ListenAndServe(":80", nil)
}
