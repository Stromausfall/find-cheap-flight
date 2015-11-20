package data

import (
	"fmt"
	"net/http"
	"github.com/stromausfall/find-cheap-flight/data/retrieve/geoinfo"
)

func DisplayDataSelection(w http.ResponseWriter, r *http.Request, googleMapsApiCredentials string, geonameAccount string) {
	arguments := CreateArguments(r, googleMapsApiCredentials, "????", false, "Find Cheap Flights - data entry II")
	
	startAirports := geoinfo.GetAirports(arguments.StartLat, arguments.StartLng, arguments.StartRange, geonameAccount)
	destinationAirports := geoinfo.GetAirports(arguments.DestLat, arguments.DestLng, arguments.DestRange, geonameAccount)
	
	fmt.Println("startAirports : ", startAirports)
	fmt.Println("destinationAirports : ", destinationAirports)
	
	DisplayPage(w, arguments)
}
