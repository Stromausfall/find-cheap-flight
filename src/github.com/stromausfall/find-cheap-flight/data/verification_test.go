package data

import (
	"fmt"
	"github.com/stromausfall/find-cheap-flight/utils"
	"testing"
	"time"
)

func TestVerifyAirportsBothNull(t *testing.T) {
	testData := FlightsToSearch{}

	error := verifyAirports(&testData)

	utils.FailIfStringDoesntHaveSubstring(t, error, "no start airport selected !</br>")
	utils.FailIfStringDoesntHaveSubstring(t, error, "no destination airport selected !</br>")
}

func TestVerifyAirportsSameAirportInBoth(t *testing.T) {
	testData := FlightsToSearch{}
	testData.destAirports = []string{"INN", "VIE"}
	testData.startAirports = []string{"INN", "PEK"}

	error := verifyAirports(&testData)

	utils.FailIfStringDoesntHaveSubstring(t, error, "the same airport (INN) can't be used for start and destination !</br>")
}

func buildValidTestData() FlightsToSearch {
	return FlightsToSearch{
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
