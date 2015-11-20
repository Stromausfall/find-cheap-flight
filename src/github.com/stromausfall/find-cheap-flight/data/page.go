package data

import (
	"strconv"
	"net/http"
	"html/template"
	"github.com/stromausfall/find-cheap-flight/utils"
)

const dataHtmlData = `
<!DOCTYPE html>
<html>
	<head>
		<meta name="viewport" content="initial-scale=1.0, user-scalable=no">
		<meta charset="utf-8">
		<title>{{.PageTitle}}</title>
		<style>
	html, body {
		height: 100%;
		margin: 0;
		padding: 0;
	}
		#map {
		height: 100%;
		float: right;
		width: 50%;
		height: 100%;
	}
	#controls {
		margin: 0;
		position: absolute;
		top: 5%;
		left: 5%;
		width: 50%;
	}
		</style>
	</head>
	<body>
		<div id="controls">
			<form action="/{{.NextPage}}" method="post">
				<b>Time</b>
				<ul>
					<li>Earliest Departure : <input type="date" name="earliestDeparture" min="{{.MinEarliestDeparture}}" value="{{.EarliestDeparture}}" {{if not .InputEnabled}}disabled{{end}}></li>
					<li>Latest Return : <input type="date" name="latestDeparture" min="{{.MinLatestDeparture}}" value="{{.LatestDeparture}}" {{if not .InputEnabled}}disabled{{end}}></li>
					<li>Minimum stay : <input type="number" name="minimumStayInput" min="1" max="1000" value="{{.MinStay}}" {{if not .InputEnabled}}disabled{{end}}></li>
					<li>Maximum stay : <input type="number" name="maximumStayInput" min="1" max="1000" value="{{.MaxStay}}" {{if not .InputEnabled}}disabled{{end}}></li>
				</ul>

				<b>Start</b>
				<ul>
					<li>latitude: <span id="startLocationLatitude"/></li>
					<li>longitude: <span id="startLocationLongitude"/></li>
					<li>range: <span id="startLocationRange"/></li>
				</ul>
				<input type="hidden" id="startLocationLatitudeInput" name="startLocationLatitudeInput" value="">
				<input type="hidden" id="startLocationLongitudeInput" name="startLocationLongitudeInput" value="">
				<input type="hidden" id="startLocationRangeInput" name="startLocationRangeInput" value="">
				<b>Destination</b>
				<ul>
					<li>latitude: <span id="destinationLocationLatitude"/></li>
					<li>longitude: <span id="destinationLocationLongitude"/></li>
					<li>range: <span id="destinationLocationRange"/></li>
				</ul>
				<input type="hidden" id="destinationLocationLatitudeInput" name="destinationLocationLatitudeInput" value="">
				<input type="hidden" id="destinationLocationLongitudeInput" name="destinationLocationLongitudeInput" value="">
				<input type="hidden" id="destinationLocationRangeInput" name="destinationLocationRangeInput" value="">
				<div><input type="submit" value="Check values"></div>
			</form>
		</div>
		<div id="map"></div>
			<script>
		
function initMap() {
	var pos1 = {lat: {{.StartLat}}, lng: {{.StartLng}}};
	var pos2 = {lat: {{.DestLat}}, lng: {{.DestLng}}};
	var posCenter = {lat: {{.CntrLat}}, lng: {{.CntrLng}}};
	var map = new google.maps.Map(document.getElementById('map'), {
		zoom: 4,
		center: posCenter,
		zoomControl: true,
		mapTypeControl: false,
		scaleControl: true,
		streetViewControl: false,
		rotateControl: false,
		showscale: true,
		disableDefaultUI: {{not .InputEnabled}}
	});
	
	var startCircle = new google.maps.Circle({
		strokeColor: '#00FF00',
		strokeOpacity: 0.8,
		strokeWeight: 2,
		fillColor: '#00FF00',
		fillOpacity: 0.35,
		map: map,
		center: pos1,
		radius: {{.StartRange}},
		draggable: {{.InputEnabled}},
		editable: {{.InputEnabled}},
    	geodesic: true
    });
	
	var destinationCircle = new google.maps.Circle({
		strokeColor: '#0000FF',
		strokeOpacity: 0.8,
		strokeWeight: 2,
		fillColor: '#0000FF',
		fillOpacity: 0.35,
		map: map,
		center: pos2,
		radius: {{.DestRange}},
		draggable: {{.InputEnabled}},
		editable: {{.InputEnabled}},
    	geodesic: true
    });
  
	window.setInterval(function(){
        document.getElementById('startLocationLatitude').innerHTML = startCircle.center.lat();
        document.getElementById('startLocationLongitude').innerHTML = startCircle.center.lng();
        document.getElementById('startLocationRange').innerHTML = Math.round(startCircle.radius / 1000) + " km";
        document.getElementById('startLocationLatitudeInput').value = startCircle.center.lat();
        document.getElementById('startLocationLongitudeInput').value = startCircle.center.lng();
        document.getElementById('startLocationRangeInput').value = startCircle.radius;
		
        document.getElementById('destinationLocationLatitude').innerHTML = destinationCircle.center.lat();
        document.getElementById('destinationLocationLongitude').innerHTML = destinationCircle.center.lng();
        document.getElementById('destinationLocationRange').innerHTML = Math.round(destinationCircle.radius / 1000) + " km";
        document.getElementById('destinationLocationLatitudeInput').value = destinationCircle.center.lat();
        document.getElementById('destinationLocationLongitudeInput').value = destinationCircle.center.lng();
        document.getElementById('destinationLocationRangeInput').value = destinationCircle.radius;
	}, 500);
}


    </script>
    <script src="https://maps.googleapis.com/maps/api/js?key={{.GoogleMapsApiCredentials}}&signed_in=true&callback=initMap"
        async defer></script>
  </body>
</html>
`

type DataEntryDisplayArgs struct {
	GoogleMapsApiCredentials string
	StartLng float32
	StartLat float32
	StartRange float32
	DestLng float32
	DestLat float32
	DestRange float32
	CntrLng float32
	CntrLat float32	
	EarliestDeparture string
	MinEarliestDeparture string
	LatestDeparture string
	MinLatestDeparture string
	MinStay int32
	MaxStay int32
	NextPage string
	InputEnabled bool
	PageTitle string
}

func createDefaultDataEntryDisplayArgs(googleMapsApiCredentials string) (*DataEntryDisplayArgs) {
	// file with default values
	result := DataEntryDisplayArgs{
		GoogleMapsApiCredentials: googleMapsApiCredentials,
		StartLat: 50.0,
		StartLng: 14.4,
		StartRange: 50000,
		DestLat: 48.2,
		DestLng: 16.3,
		DestRange: 50000,
		MinStay: 1,
		MaxStay: 5,
		NextPage: "",
		InputEnabled: false,
		PageTitle: "default page title",
	}
	
	// calculate default dates
	result.MinEarliestDeparture = utils.DateStringNow(0)
	result.MinLatestDeparture = utils.DateStringNow(1)
	result.EarliestDeparture = result.MinEarliestDeparture
	result.LatestDeparture = result.MinLatestDeparture
	
	// calculate center 
	result.CntrLat = (result.StartLat + result.DestLat) / 2
	result.CntrLng = (result.StartLng + result.DestLng) / 2
	
	return &result
}

func getStringFormValue(r *http.Request, storeValue *string, formValueKey string) {
	if r.FormValue(formValueKey) != "" {
		*storeValue = r.FormValue(formValueKey)
	}
}

func getFloatFormValue(r *http.Request, storeValue *float32, formValueKey string) {
	if r.FormValue(formValueKey) != "" {
		value, err := strconv.ParseFloat(r.FormValue(formValueKey), 32)
		
		if err == nil {
			*storeValue = float32(value)
		}
	}
}

func getIntFormValue(r *http.Request, storeValue *int32, formValueKey string) {
	if r.FormValue(formValueKey) != "" {
		value, err := strconv.ParseInt(r.FormValue(formValueKey), 10, 32)
		
		if err == nil {
			*storeValue = int32(value)
		}
	}
}

func CreateArguments(
		r *http.Request,
		googleMapsApiCredentials string,
		nextPage string,
		inputEnabled bool,
		pageTitle string) (*DataEntryDisplayArgs) {
	// we need this in order to get POST form data
	r.ParseMultipartForm(15485760)
		
	arguments := createDefaultDataEntryDisplayArgs(googleMapsApiCredentials)
	arguments.NextPage = nextPage
	arguments.InputEnabled = inputEnabled
	arguments.PageTitle = pageTitle
	
	// get values from POST
	getStringFormValue(r, &arguments.EarliestDeparture, "earliestDeparture")
	getStringFormValue(r, &arguments.LatestDeparture, "latestDeparture")
	getFloatFormValue(r, &arguments.StartLat, "startLocationLatitudeInput")
	getFloatFormValue(r, &arguments.StartLng, "startLocationLongitudeInput")
	getFloatFormValue(r, &arguments.StartRange, "startLocationRangeInput")
	getFloatFormValue(r, &arguments.DestLat, "destinationLocationLatitudeInput")
	getFloatFormValue(r, &arguments.DestLng, "destinationLocationLongitudeInput")
	getFloatFormValue(r, &arguments.DestRange, "destinationLocationRangeInput")
	getIntFormValue(r, &arguments.MinStay, "minimumStayInput")
	getIntFormValue(r, &arguments.MaxStay, "maximumStayInput")
	
	return arguments
}

func DisplayPage(
		w http.ResponseWriter,
		arguments *DataEntryDisplayArgs) {
	// create, initialize and use the template
	uninitializedTemplate := template.New("Data entry template")
	initializedTempalte, err := uninitializedTemplate.Parse(dataHtmlData)
	utils.CheckErr(err, "problem while parsing template")

	err = initializedTempalte.Execute(w, arguments)
	utils.CheckErr(err, "problem while executing template")
}
