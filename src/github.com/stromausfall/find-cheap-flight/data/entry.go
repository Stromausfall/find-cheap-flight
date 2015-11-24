package data

import (
	"net/http"
)

func DisplayDataEntry(w http.ResponseWriter, r *http.Request, googleMapsApiCredentials string) {
	arguments := CreateArguments(r, googleMapsApiCredentials, "selectAirports", true, "Find Cheap Flights - data entry I")
	arguments.SubmitButtonText = "Submit Data"

	DisplayPage(w, arguments)
}
