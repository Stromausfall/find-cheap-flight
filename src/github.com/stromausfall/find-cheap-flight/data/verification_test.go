package data

import (
	"fmt"
	"github.com/stromausfall/find-cheap-flight/utils"
	"testing"
	"time"
)

func TestVerifyAirportsBothNull(t *testing.T) {
	testData := flightsToSearch{}

	error := verifyAirports(&testData)

	utils.FailIfStringDoesntHaveSubstring(t, error, "no start airport selected !</br>")
	utils.FailIfStringDoesntHaveSubstring(t, error, "no destination airport selected !</br>")
}

func TestVerifyAirportsSameAirportInBoth(t *testing.T) {
	testData := flightsToSearch{}
	testData.destAirports = []string{"INN", "VIE"}
	testData.startAirports = []string{"INN", "PEK"}

	error := verifyAirports(&testData)

	utils.FailIfStringDoesntHaveSubstring(t, error, "the same airport (INN) can't be used for start and destination !</br>")
}

func buildValidTestData() flightsToSearch {
	return flightsToSearch{
		minimumStay:           3,
		maximumStay:           5,
		earliestDepartureDate: utils.RawDateNow(1),
		latestDepartureDate:   utils.RawDateNow(3),
		startAirports:         []string{"INN", "VIE"},
		destAirports:          []string{"PEK", "PVG"},
	}
}

func TestNoErrorMessagesIfValid(t *testing.T) {
	testData := buildValidTestData()
	error := verifyDates(&testData)

	if error != "" {
		fmt.Printf("unexpected error '%v' with valid data !</br>", error)
		t.Fail()
	}
}

func TestVerifyDatesLatestBeforeEarliestDate(t *testing.T) {
	testData := buildValidTestData()
	testData.earliestDepartureDate = utils.RawDateNow(5)

	error := verifyDates(&testData)
	expectedMessage := fmt.Sprintf("the latestDepartureDate (%v) is temporally before the earliestDepartureDate (%v) !</br>", testData.latestDepartureDate, testData.earliestDepartureDate)

	utils.FailIfStringDoesntHaveSubstring(t, error, expectedMessage)
}

func TestVerifyDatesAreNil(t *testing.T) {
	testData := buildValidTestData()
	testData.earliestDepartureDate = time.Time{}
	testData.latestDepartureDate = time.Time{}

	error := verifyDates(&testData)

	utils.FailIfStringDoesntHaveSubstring(t, error, "earliestDepartureDate was not set !</br>")
	utils.FailIfStringDoesntHaveSubstring(t, error, "latestDepartureDate was not set !</br>")
}

func TestVerifyStayIsNegative(t *testing.T) {
	testData := buildValidTestData()
	testData.minimumStay = -1
	testData.maximumStay = -3

	error := verifyDates(&testData)

	utils.FailIfStringDoesntHaveSubstring(t, error, fmt.Sprintf("minimumStay can not be negative (%v) !</br>", testData.minimumStay))
	utils.FailIfStringDoesntHaveSubstring(t, error, fmt.Sprintf("maximumStay can not be negative (%v) !</br>", testData.maximumStay))
}

func TestVerifyStayWasZero(t *testing.T) {
	testData := buildValidTestData()
	testData.minimumStay = 0
	testData.maximumStay = 0

	error := verifyDates(&testData)

	utils.FailIfStringDoesntHaveSubstring(t, error, "minimumStay can not be 0 !</br>")
	utils.FailIfStringDoesntHaveSubstring(t, error, "maximumStay can not be 0 !</br>")
}

func TestVerifyMaxStaySmallerThanMinStay(t *testing.T) {
	testData := buildValidTestData()
	testData.minimumStay = 5
	testData.maximumStay = 4

	error := verifyDates(&testData)

	utils.FailIfStringDoesntHaveSubstring(t, error, fmt.Sprintf("minimumStay was bigger (%v) than maximumStay (%v) !</br>", testData.minimumStay, testData.maximumStay))
}

func TestCalculatePossibleQueriesOnlyOne(t *testing.T) {
	actualDate := time.Now()

	argument := flightsToSearch{
		minimumStay:           6,
		maximumStay:           6,
		earliestDepartureDate: actualDate,
		latestDepartureDate:   actualDate,
		startAirports:         []string{"INN"},
		destAirports:          []string{"PEK"},
	}
	expectedQuery := FlightQuery{
		stayDuration:  6,
		departureData: actualDate,
		startAirport:  "INN",
		destAirport:   "PEK",
	}

	queries := calculatePossibleQueries(&argument)

	if len(queries) != 1 {
		t.Errorf("there should be exactly on query - but there were : %v", len(queries))
	}

	if queries[0] != expectedQuery {
		t.Errorf("returned query %v is not as expected %v", queries[0], expectedQuery)
	}
}

func TestCalculatePossibleQueriesCheckStayTime(t *testing.T) {
	actualDate := time.Now()

	argument := flightsToSearch{
		minimumStay:           5,
		maximumStay:           6,
		earliestDepartureDate: actualDate,
		latestDepartureDate:   actualDate,
		startAirports:         []string{"INN"},
		destAirports:          []string{"PEK"},
	}
	expectedQuery1 := FlightQuery{
		stayDuration:  5,
		departureData: actualDate,
		startAirport:  "INN",
		destAirport:   "PEK",
	}
	expectedQuery2 := FlightQuery{
		stayDuration:  6,
		departureData: actualDate,
		startAirport:  "INN",
		destAirport:   "PEK",
	}

	queries := calculatePossibleQueries(&argument)

	if len(queries) != 2 {
		t.Errorf("there should be exactly two queries - but there were : %v", len(queries))
	}

	if queries[0] != expectedQuery1 {
		t.Errorf("returned query1 %v is not as expected %v", queries[0], expectedQuery1)
	}

	if queries[1] != expectedQuery2 {
		t.Errorf("returned query2 %v is not as expected %v", queries[1], expectedQuery2)
	}
}

func TestCalculatePossibleQueriesCheckDepartureDate(t *testing.T) {
	actualDate1 := utils.DateFromString("2015-11-25")
	actualDate2 := utils.DateFromString("2015-11-26")

	argument := flightsToSearch{
		minimumStay:           6,
		maximumStay:           6,
		earliestDepartureDate: actualDate1,
		latestDepartureDate:   actualDate2,
		startAirports:         []string{"INN"},
		destAirports:          []string{"PEK"},
	}
	expectedQuery1 := FlightQuery{
		stayDuration:  6,
		departureData: actualDate1,
		startAirport:  "INN",
		destAirport:   "PEK",
	}
	expectedQuery2 := FlightQuery{
		stayDuration:  6,
		departureData: actualDate2,
		startAirport:  "INN",
		destAirport:   "PEK",
	}

	queries := calculatePossibleQueries(&argument)

	if len(queries) != 2 {
		t.Errorf("there should be exactly two queries - but there were : %v", len(queries))
	}

	if queries[0] != expectedQuery1 {
		t.Errorf("returned query1 %v is not as expected %v", queries[0], expectedQuery1)
	}

	if queries[1] != expectedQuery2 {
		t.Errorf("returned query2 %v is not as expected %v", queries[1], expectedQuery2)
	}
}

func TestCalculatePossibleQueriesCheckStartAirports(t *testing.T) {
	actualDate := time.Now()

	argument := flightsToSearch{
		minimumStay:           6,
		maximumStay:           6,
		earliestDepartureDate: actualDate,
		latestDepartureDate:   actualDate,
		startAirports:         []string{"INN", "VIE"},
		destAirports:          []string{"PEK"},
	}
	expectedQuery1 := FlightQuery{
		stayDuration:  6,
		departureData: actualDate,
		startAirport:  "INN",
		destAirport:   "PEK",
	}
	expectedQuery2 := FlightQuery{
		stayDuration:  6,
		departureData: actualDate,
		startAirport:  "VIE",
		destAirport:   "PEK",
	}

	queries := calculatePossibleQueries(&argument)

	if len(queries) != 2 {
		t.Errorf("there should be exactly two queries - but there were : %v", len(queries))
	}

	if queries[0] != expectedQuery1 {
		t.Errorf("returned query1 %v is not as expected %v", queries[0], expectedQuery1)
	}

	if queries[1] != expectedQuery2 {
		t.Errorf("returned query2 %v is not as expected %v", queries[1], expectedQuery2)
	}
}

func TestCalculatePossibleQueriesCheckDestAirports(t *testing.T) {
	actualDate := time.Now()

	argument := flightsToSearch{
		minimumStay:           6,
		maximumStay:           6,
		earliestDepartureDate: actualDate,
		latestDepartureDate:   actualDate,
		startAirports:         []string{"INN"},
		destAirports:          []string{"PEK", "PVG"},
	}
	expectedQuery1 := FlightQuery{
		stayDuration:  6,
		departureData: actualDate,
		startAirport:  "INN",
		destAirport:   "PEK",
	}
	expectedQuery2 := FlightQuery{
		stayDuration:  6,
		departureData: actualDate,
		startAirport:  "INN",
		destAirport:   "PVG",
	}

	queries := calculatePossibleQueries(&argument)

	if len(queries) != 2 {
		t.Errorf("there should be exactly two queries - but there were : %v", len(queries))
	}

	if queries[0] != expectedQuery1 {
		t.Errorf("returned query1 %v is not as expected %v", queries[0], expectedQuery1)
	}

	if queries[1] != expectedQuery2 {
		t.Errorf("returned query2 %v is not as expected %v", queries[1], expectedQuery2)
	}
}

func TestCalculatePossibleQueriesCheckCount(t *testing.T) {
	actualDate1 := utils.DateFromString("2015-11-25")
	actualDate2 := utils.DateFromString("2015-11-26")

	argument := flightsToSearch{
		minimumStay:           5,
		maximumStay:           6,
		earliestDepartureDate: actualDate1,
		latestDepartureDate:   actualDate2,
		startAirports:         []string{"INN", "VIE"},
		destAirports:          []string{"PEK", "PVG"},
	}

	queries := calculatePossibleQueries(&argument)

	if len(queries) != 16 {
		t.Errorf("there should be exactly eight queries - but there were : %v", len(queries))
	}
}
