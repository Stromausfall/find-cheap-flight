package data

import (
	"testing"
	"strings"
)

func TestVerifyAirportsBothNull(t *testing.T) {
	testData := FlightsToSearch{}
	
	error := verifyAirports(&testData)
	
	if strings.Contains(error, "no start airport selected !") {
		if strings.Contains(error, "no destination airport selected !") {
			// we expect this !
			return
		}
	}
	
	t.Fail()
}

func TestVerifyAirportsSameAirportInBoth(t *testing.T) {
	testData := FlightsToSearch{}
	testData.destAirports = []string { "INN", "VIE" }
	testData.startAirports = []string { "INN", "PEK" }
	
	error := verifyAirports(&testData)
	
	if strings.Contains(error, "the same airport (INN) can't be used for start and destination !") {
		return
	}
	
	t.Fail()
}