package geoinfo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stromausfall/find-cheap-flight/utils"
	"net/http"
	"net/url"
	"strconv"
)

type geonameContainer struct {
	Distance  string
	GeonameId int32
}

type geonamesContainer struct {
	Geonames []geonameContainer
}

type geonameInfos struct {
	AlternateNames []geonameInfo
}

type geonameInfo struct {
	Name string
	Lang string
}

type AirportData struct {
	Distance float32
	Iata     string
	Name     string
}

func retrieveValueFromPostForm(homepage string, postValues url.Values) []byte {
	resp, err := http.PostForm(homepage, postValues)
	utils.CheckErr(err, "problem while retrieving value from post form")

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf.Bytes()
}

func GetAirports(startLat, startLng, startRange float32, geonameAccount string) []AirportData {
	geonamesIds := retrieveAirportIds(startLat, startLng, startRange, geonameAccount)
	airportData := make([]AirportData, 0)

	for _, geonameId := range geonamesIds {
		iata, name := retrieveAirportNames(geonameId, geonameAccount)

		distance, err := strconv.ParseFloat(geonameId.Distance, 32)
		utils.CheckErr(err, "problem while retrieving from post form")

		airportDataElement := AirportData{
			Distance: float32(distance),
			Iata:     iata,
			Name:     name,
		}

		// only add the aiport if it has correct names
		if (iata != "") && (name != "") {
			airportData = append(airportData, airportDataElement)
		}
	}

	return airportData
}

func retrieveAirportNames(geonameId geonameContainer, geonameAccount string) (iata, name string) {
	homepage := "http://api.geonames.org/getJSON"
	postValues := url.Values{
		"geonameId": {fmt.Sprintf("%v", geonameId.GeonameId)},
		"username":  {geonameAccount}}

	data := retrieveValueFromPostForm(homepage, postValues)

	value := geonameInfos{}
	unmarshallError := json.Unmarshal(data, &value)
	utils.CheckErr(unmarshallError, "error while unmarshaling geo names")

	for _, alternateName := range value.AlternateNames {
		if alternateName.Lang == "iata" {
			iata = alternateName.Name
		}
		if alternateName.Lang == "en" {
			name = alternateName.Name
		}
	}

	return iata, name
}

func retrieveAirportIds(startLat, startLng, startRange float32, geonameAccount string) []geonameContainer {
	homepage := "http://api.geonames.org/findNearbyJSON"
	postValues := url.Values{
		"lat":      {fmt.Sprintf("%v", startLat)},
		"lng":      {fmt.Sprintf("%v", startLng)},
		"fcode":    {"AIRP"},
		"radius":   {fmt.Sprintf("%v", startRange/1000)},
		"maxRows":  {"100"},
		"username": {geonameAccount}}

	data := retrieveValueFromPostForm(homepage, postValues)
	value := geonamesContainer{}
	unmarshallError := json.Unmarshal(data, &value)
	utils.CheckErr(unmarshallError, "error while unmarshaling geo airports")

	return value.Geonames
}
