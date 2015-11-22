package data

import (
	"net/http"
	"github.com/stromausfall/find-cheap-flight/data/retrieve/geoinfo"
	"bytes"
	"html/template"
	"github.com/stromausfall/find-cheap-flight/utils"
)

type selectionData struct {
	Name string
	ShortName string
	Items *[]geoinfo.AirportData
}

const validAirportsHtmlData = `
	<b>{{$.Name}}</b>
	<ul>
{{range $index, $result := $.Items}}
		<li><input type="checkbox" name="selected-{{$.ShortName}}-airport#{{$index}}" value="{{$result.Iata}}" id="{{$.ShortName}}-id-{{$index}}">
            {{$result.Iata}} ({{$result.Name}}) - Distance : {{$result.Distance}}</br>
		</li>
{{end}}
	</ul>
`

func createAiportsHtmls(name, shortName string, items *[]geoinfo.AirportData) template.HTML {
	arguments := selectionData{
		Name:name,
		ShortName:shortName,
		Items:items,
	}
	
	// create, initialize and use the template
	uninitializedTemplate := template.New("valid airports entry template")
	initializedTempalte, err := uninitializedTemplate.Parse(validAirportsHtmlData)
	utils.CheckErr(err, "problem while parsing template")

	var stringBuffer bytes.Buffer  
	err = initializedTempalte.Execute(&stringBuffer, arguments)
	utils.CheckErr(err, "problem while executing template")
	
	resultString := stringBuffer.String()
	unescapable := template.HTML(resultString)
	
	return unescapable
}

func DisplayDataSelection(w http.ResponseWriter, r *http.Request, googleMapsApiCredentials string, geonameAccount string) {
	arguments := CreateArguments(r, googleMapsApiCredentials, "dataVerification", false, "Find Cheap Flights - data entry II")
	arguments.SubmitButtonText = "Submit Airports"

	startAirports := geoinfo.GetAirports(arguments.StartLat, arguments.StartLng, arguments.StartRange, geonameAccount)
	destinationAirports := geoinfo.GetAirports(arguments.DestLat, arguments.DestLng, arguments.DestRange, geonameAccount)

	startAirportsHtml := createAiportsHtmls("Start Airports", "start", &startAirports)
	destAirportsHtml := createAiportsHtmls("Destination Airports", "dest", &destinationAirports)

	arguments.DataToAddBeforeSubmitButton = startAirportsHtml + destAirportsHtml

	DisplayPage(w, arguments)
}
