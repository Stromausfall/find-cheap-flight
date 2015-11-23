package data

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"github.com/stromausfall/find-cheap-flight/utils"
)

func collectValues(r *http.Request, prefix string) []string {
	result := make([]string, 0)
	
	for key := range r.Form {
		if strings.HasPrefix(key, prefix) {
			result = append(result, r.FormValue(key))
		}
	}
	
	return result
}

type FlightsToSearch struct {
	minimumStay int32
	maximumStay int32
	earliestDepartureDate time.Time
	latestDepartureDate time.Time
	startAirports []string
	destAirports []string
}

func parseFormValues(r *http.Request) FlightsToSearch {
	var result FlightsToSearch	
	var earliestDepartureDateRawString string
	var latestDepartureDateRawString string
	
	getIntFormValue(r, &result.minimumStay, "minimumStayInput")
	getIntFormValue(r, &result.maximumStay, "maximumStayInput")
	getStringFormValue(r, &earliestDepartureDateRawString, "earliestDeparture")
	getStringFormValue(r, &latestDepartureDateRawString, "latestDeparture")
	result.earliestDepartureDate = utils.DateFromString(earliestDepartureDateRawString)
	result.latestDepartureDate = utils.DateFromString(latestDepartureDateRawString)
	result.startAirports = collectValues(r, "selected-start-airport")
	result.destAirports = collectValues(r, "selected-dest-airport")
	
	return result
}

func verifyAirports(data *FlightsToSearch) string {
	error := ""
	
	if len(data.startAirports) == 0 {
		error = error + "no start airport selected !</br>"
	}
	
	if len(data.destAirports) == 0 {
		error = error + "no destination airport selected !</br>"
	}
	
	for _, startAirport := range data.startAirports {
		for _, destAirport := range data.destAirports {
			if startAirport == destAirport {
				error = error + "the same airport (" + startAirport + ") can't be used for start and destination !</br>"
			}
		}
	}
	
	return error
}

func DisplayDataVerification(w http.ResponseWriter, r *http.Request, googleMapsApiCredentials string, geonameAccount string) {
	arguments := CreateArguments(r, googleMapsApiCredentials, "????", false, "Find Cheap Flights - data verification")
	arguments.SubmitButtonText = "Query flights"
	
	result := parseFormValues(r)
	
	verifyAirports(&result)
	
	fmt.Println("result : ", result)
}
