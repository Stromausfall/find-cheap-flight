package data

import (
	"net/http"
)

func DisplayDataEntry(w http.ResponseWriter, r *http.Request, googleMapsApiCredentials string) {
	DisplayPage(w, r, googleMapsApiCredentials, "selection", true, "Find Cheap Flights - data entry I")
}
