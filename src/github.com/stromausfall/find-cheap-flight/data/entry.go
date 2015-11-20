package data

import (
	"net/http"
)

func DisplayDataEntry(w http.ResponseWriter, r *http.Request, googleMapsApiCredentials string) {
	arguments := CreateArguments(r, googleMapsApiCredentials, "selection", true, "Find Cheap Flights - data entry I")
	DisplayPage(w, arguments)
}
