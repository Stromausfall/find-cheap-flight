package data

import (
	"net/http"
	"fmt"
	"bytes"
	"net/url"
)

func DisplayDataSelection(w http.ResponseWriter, r *http.Request, googleMapsApiCredentials string, genomeAccount string) {
	arguments := DisplayPage(w, r, googleMapsApiCredentials, "????", false, "Find Cheap Flights - data entry II")
	
	homepage := "http://api.geonames.org/findNearby"
	postValues := url.Values{
			"lat": {fmt.Sprintf("%v", arguments.StartLat)},
			"lng": {fmt.Sprintf("%v", arguments.StartLng)},
			"fcode": {"AIRP"},
			"radius": {fmt.Sprintf("%v", arguments.StartRange/1000)},
			"maxRows": {"100"},
			"type": {"json"},
			"username": {genomeAccount}}

	resp, err := http.PostForm(homepage, postValues)
	if err != nil {
		fmt.Println("error", err.Error())
	}
	
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	s := buf.String() // Does a complete copy of the bytes in the buffer.
	
	fmt.Println("\n\nbody : " + s)
	
	// use the following, to get more information about the place
	// including the IATA designation !!
	// http://api.geonames.org/get?geonameId=6299669&username=demo&type=json
}
