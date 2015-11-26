package data

import (
	"fmt"
	"github.com/stromausfall/find-cheap-flight/utils"
	"html/template"
	"net/http"
	"strings"
	"time"
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

type flightsToSearch struct {
	minimumStay           int32
	maximumStay           int32
	earliestDepartureDate time.Time
	latestDepartureDate   time.Time
	startAirports         []string
	destAirports          []string
}

type FlightQuery struct {
	stayDuration  int32
	departureData time.Time
	startAirport  string
	destAirport   string
}

func parseFormValues(r *http.Request) flightsToSearch {
	var result flightsToSearch
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

func verifyDates(data *flightsToSearch) string {
	error := ""
	nilDate := time.Time{}

	if data.earliestDepartureDate == nilDate {
		error = error + "earliestDepartureDate was not set !</br>"
	}

	if data.latestDepartureDate == nilDate {
		error = error + "latestDepartureDate was not set !</br>"
	}

	if data.earliestDepartureDate.After(data.latestDepartureDate) {
		error = error + fmt.Sprintf("the latestDepartureDate (%v) is temporally before the earliestDepartureDate (%v) !</br>", data.latestDepartureDate, data.earliestDepartureDate)
	}

	if data.minimumStay < 0 {
		error = error + fmt.Sprintf("minimumStay can not be negative (%v) !</br>", data.minimumStay)
	}

	if data.maximumStay < 0 {
		error = error + fmt.Sprintf("maximumStay can not be negative (%v) !</br>", data.maximumStay)
	}

	if data.minimumStay == 0 {
		error = error + "minimumStay can not be 0 !</br>"
	}

	if data.maximumStay == 0 {
		error = error + "maximumStay can not be 0 !</br>"
	}

	if data.maximumStay < data.minimumStay {
		error = error + fmt.Sprintf("minimumStay was bigger (%v) than maximumStay (%v) !</br>", data.minimumStay, data.maximumStay)
	}

	return error
}

func calculatePossibleQueries(data *flightsToSearch) []FlightQuery {
	result := make([]FlightQuery, 0)

	// stay
	for stayDuration := data.minimumStay; stayDuration <= data.maximumStay; stayDuration++ {
		// date
		for departureDate := data.earliestDepartureDate; !departureDate.After(data.latestDepartureDate); departureDate = departureDate.Add(time.Hour * 24) {
			// start airport
			for _, startAirport := range data.startAirports {
				// dest airport
				for _, destAirport := range data.destAirports {
					// create query
					flightQuery := FlightQuery{
						stayDuration:  stayDuration,
						departureData: departureDate,
						startAirport:  startAirport,
						destAirport:   destAirport,
					}

					// add it
					result = append(result, flightQuery)
				}
			}
		}
	}

	return result
}

func verifyAirports(data *flightsToSearch) string {
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
	arguments := CreateArguments(r, googleMapsApiCredentials, "", false, "Find Cheap Flights - data verification")
	arguments.SubmitButtonText = "Change Data"

	result := parseFormValues(r)
	error := ""

	// the verification already assumes that the data
	// has been received by the correct homepage
	// --> submitting incorrect data in the post form will just throw an exception
	// without feedback to the user (no need to handle blatantly incorrect usage)
	error = error + verifyAirports(&result)
	error = error + verifyDates(&result)

	if error == "" {
		possibleQueriesCount := len(calculatePossibleQueries(&result))
		
		if possibleQueriesCount <= 50 {
			arguments.DataToAddBeforeSubmitButton = arguments.DataToAddBeforeSubmitButton + template.HTML(fmt.Sprintf("<b>Found <font color=\"green\">%v</font> queries</b></br>", possibleQueriesCount))
			arguments.DataToAddAfterSubmitForm = template.HTML("<form action=\"/calculate\" method=\"post\"><div><input type=\"submit\" value=\"Calculate\"></div></form>")
		} else {
			arguments.DataToAddBeforeSubmitButton = arguments.DataToAddBeforeSubmitButton + template.HTML(fmt.Sprintf("<b>Only 50 queries are supported but there were <font color=\"red\">%v</font></b></br>", possibleQueriesCount))
		}
	} else {
		arguments.DataToAddBeforeSubmitButton = template.HTML("<font color=\"red\"><b>" + error + "</b></font>")
	}

	// TODO : MARSHAL queries to JSON and transmit that to the calculate page !
	// TODO : OR pass an ID and keep the queries internally in a map or so

	DisplayPage(w, arguments)
}
