package data

import (
	"net/http"
	"fmt"
	"bytes"
	"net/url"
	"gopkg.in/yaml.v2"
)

type GeonameContainer struct {
	GeonameId int
	Population int
}

type GeonamesContainer struct {
	Geonames []GeonameContainer
}

func DisplayDataSelection(w http.ResponseWriter, r *http.Request, googleMapsApiCredentials string, genomeAccount string) {
	arguments := DisplayPage(w, r, googleMapsApiCredentials, "????", false, "Find Cheap Flights - data entry II")
	
	homepage := "http://api.geonames.org/findNearbyJSON"
	postValues := url.Values{
			"lat": {fmt.Sprintf("%v", arguments.StartLat)},
			"lng": {fmt.Sprintf("%v", arguments.StartLng)},
			"fcode": {"AIRP"},
			"radius": {fmt.Sprintf("%v", arguments.StartRange/1000)},
			"maxRows": {"100"},
			"username": {genomeAccount}}

	resp, err := http.PostForm(homepage, postValues)
	if err != nil {
		fmt.Println("error", err.Error())
	}
	
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	data := buf.Bytes()
	s := buf.String() // Does a complete copy of the bytes in the buffer.
	
	
	value := GeonamesContainer{}
	unmarshallError := yaml.Unmarshal(data, &value)

	if unmarshallError != nil {
		fmt.Printf(fmt.Sprintf("%v", unmarshallError))
		panic("error while unmarshaling configFile")
	}
	
	fmt.Printf("\nconverted: %v", value)
	
	fmt.Println("\n\nbody : " + s)
	
	// choosing IATA as language we only need to return the name :D
	// BUT we also need the name of the airport, therefore use the following !
	// http://api.geonames.org/getJSON?geonameId=6299669&username=stromausfall
}
