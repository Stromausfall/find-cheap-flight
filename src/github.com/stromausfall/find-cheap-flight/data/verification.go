package data

import (
	"fmt"
	"net/http"
	"strings"
)

func PrintValue(r *http.Request, prefix string) {
	fmt.Println("--- prefix : ", prefix)
	for key := range r.Form {
		if strings.HasPrefix(key, prefix) {
			fmt.Println("\t", key, " - ", r.FormValue(key))
		}
	}
}

func DisplayDataVerification(w http.ResponseWriter, r *http.Request, googleMapsApiCredentials string, geonameAccount string) {
	arguments := CreateArguments(r, googleMapsApiCredentials, "????", false, "Find Cheap Flights - data verification")
	arguments.SubmitButtonText = "Query flights"
	
	PrintValue(r, "selected-start-airport")
	PrintValue(r, "selected-dest-airport")
	PrintValue(r, "earliestDeparture")
	PrintValue(r, "latestDeparture")
	PrintValue(r, "minimumStayInput")
	PrintValue(r, "maximumStayInput")

	DisplayPage(w, arguments)
}
