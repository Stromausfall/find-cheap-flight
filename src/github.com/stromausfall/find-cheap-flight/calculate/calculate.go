package calculate

import (
	"code.google.com/p/google-api-go-client/googleapi/transport"
	"code.google.com/p/google-api-go-client/qpxexpress/v1"
	"encoding/json"
	"fmt"
	"github.com/stromausfall/find-cheap-flight/data"
	"github.com/stromausfall/find-cheap-flight/utils"
	"html/template"
	"net/http"
	"sync"
	"time"
)

type queryDisplayData struct {
	Description string
	Cost        string
}

const submitDataHtmlData = `
<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="refresh" content="0; url=http://localhost/refreshResults">
		<meta name="viewport" content="initial-scale=1.0, user-scalable=no">
		<meta charset="utf-8">
		<title>Submiting Queries</title>
	</head>
	<body>
	</body>
</html>
`

const dataHtmlData = `
<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="refresh" content="2; url=http://localhost/refreshResults">
		<meta name="viewport" content="initial-scale=1.0, user-scalable=no">
		<meta charset="utf-8">
		<title>Calculating Flights</title>
	</head>
	<body>
	<b>Flight queries</b>
	<ul>
{{range $index, $result := $.Results}}
		<li><font color="gray">#{{$index}} - <b>{{$result.Cost}}</b> <i>{{$result.Description}}</i></font></li>
{{end}}
	</ul>
  </body>
</html>
`

type pageArgumentsData struct {
	Results []queryDisplayData
}

var flightQueries []data.FlightQuery
var costLock = make(map[data.FlightQuery]string)
var costLockMutex sync.Mutex

func countDays(start, end time.Time) int {
	endDateString := utils.DateToString(end)
	days := 0

	for element := start; utils.DateToString(element) != endDateString; element = element.Add(time.Hour * time.Duration(24)) {
		days = days + 1
	}

	return days
}

func DisplayRefreshPagePage(w http.ResponseWriter, r *http.Request) {
	toDisplay := make([]queryDisplayData, 0)

	for _, element := range flightQueries {
		startDateString := utils.DateToString(element.StartDate)
		endDateString := utils.DateToString(element.BackDate)
		duration := countDays(element.StartDate, element.BackDate)

		costLockMutex.Lock()
		price, ok := costLock[element]
		costLockMutex.Unlock()

		priceString := "???"

		if ok {
			priceString = price
		}

		displayData := queryDisplayData{
			Description: fmt.Sprintf("From %v to %v (%v days) - %v to %v", startDateString, endDateString, duration, element.StartAirport, element.DestAirport),
			Cost:        priceString,
		}

		toDisplay = append(toDisplay, displayData)
	}

	// create, initialize and use the template
	uninitializedTemplate := template.New("query calculation template")
	initializedTempalte, err := uninitializedTemplate.Parse(dataHtmlData)
	utils.CheckErr(err, "problem while parsing template")
	arguments := pageArgumentsData{
		Results: toDisplay,
	}
	err = initializedTempalte.Execute(w, arguments)
	utils.CheckErr(err, "problem while executing template")
}

func DisplayCalculatePage(w http.ResponseWriter, r *http.Request, googleQPXExpressCredentials string) {
	// we need this in order to get POST form data
	r.ParseMultipartForm(15485760)

	postFormData := []byte(r.PostFormValue("flightQueries"))

	flightQueries = make([]data.FlightQuery, 0)
	error := json.Unmarshal(postFormData, &flightQueries)
	utils.CheckErr(error, "Unable to cast received post form data to FlightQueries !")

	// create, initialize and use the template
	uninitializedTemplate := template.New("submit calculation template")
	initializedTempalte, err := uninitializedTemplate.Parse(submitDataHtmlData)
	utils.CheckErr(err, "problem while parsing template")
	err = initializedTempalte.Execute(w, nil)
	utils.CheckErr(err, "problem while executing template")

	for _, query := range flightQueries {
		go performQuery(query, googleQPXExpressCredentials)
		
		// there is a limit of 10 queries/second/user for QPX Express FREE
		// if we do about 1 query/second/user then it should work !
		time.Sleep(time.Second)
	}
}

func performQuery(query data.FlightQuery, googleQPXExpressCredentials string) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: googleQPXExpressCredentials}}
	service, err := qpxexpress.New(client)
	if err != nil {
		panic(err)
	}

	tripRequest := qpxexpress.TripsSearchRequest{}

	trip1 := qpxexpress.SliceInput{}
	trip1.Date = utils.DateToString(query.StartDate)
	trip1.Destination = query.DestAirport
	trip1.Origin = query.StartAirport

	trip2 := qpxexpress.SliceInput{}
	trip2.Date = utils.DateToString(query.BackDate)
	trip2.Destination = query.StartAirport
	trip2.Origin = query.DestAirport

	trips := []*qpxexpress.SliceInput{}
	trips = append(trips, &trip1)
	trips = append(trips, &trip2)

	tripRequest.Request = &qpxexpress.TripOptionsRequest{}
	tripRequest.Request.Slice = trips
	tripRequest.Request.Refundable = false
	tripRequest.Request.Passengers = &qpxexpress.PassengerCounts{}
	tripRequest.Request.Passengers.AdultCount = 1
	// we only need one solution (it is already the cheapest one)
	tripRequest.Request.Solutions = 1

	result, err := service.Trips.Search(&tripRequest).Do()
	utils.CheckErr(err, "unable to retrieve trips from Google QPX Express")

	costLockMutex.Lock()
	// the first trip is already the cheapest one !
	if len(result.Trips.TripOption) == 0 {
		costLock[query] = "not found !"
	} else {
		costLock[query] = result.Trips.TripOption[0].SaleTotal
	}
	costLockMutex.Unlock()
}
