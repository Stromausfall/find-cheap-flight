package data

import (
	"net/http"
)

func DisplayDataSelection(w http.ResponseWriter, r *http.Request, googleMapsApiCredentials string) {
	
	DisplayPage(w, r, googleMapsApiCredentials, "????", false, "Find Cheap Flights - data entry II")
}
