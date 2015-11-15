package data

import (
	"net/http"
	"html/template"
	"github.com/stromausfall/find-cheap-flight/utils"
)

const selectionHtmlData = `
<!DOCTYPE html>
<html>
	<head>
		<meta name="viewport" content="initial-scale=1.0, user-scalable=no">
		<meta charset="utf-8">
		<title>Draggable directions</title>
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
			<script>
				function getToday() {
					var today = new Date();
					return today.getFullYear()+'-'+(today.getMonth()+1)+'-'+today.getDate();
				}
			</script>
			
			<form action="/check" method="post">
				<b>Time</b>
				<ul>
					<li>Earliest Departure : <input id="earliestReturnInput" type="date" name="bday" disabled></li>
					<li>Latest Return : <input id="latestReturnInput" type="date" name="bday" disabled></li>
					<li>Minimum stay : <input type="number" name="minimumStayInput" min="1" max="1000" disabled></li>
					<li>Maximum stay : <input type="number" name="maximumStayInput" min="1" max="1000" disabled></li>
				</ul>

				<b>Start</b>
				<ul>
					<li>latitude: <span id="startLocationLatitude"/></li>
					<li>longitude: <span id="startLocationLongitude"/></li>
					<li>range: <span id="startLocationRange"/></li>
				</ul>
				<input type="hidden" id="startLocationLatitudeInput" name="startLocationLatitudeInput" value="">
				<input type="hidden" id="startLocationLongitudeInput"name="startLocationLongitudeInput" value="">
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
	var pos1 = {lat: 50.0, lng: 14.4};
	var pos2 = {lat: 48.2, lng: 16.3};
	var posCenter = {lat: 49.1, lng: 15.35};
	var map = new google.maps.Map(document.getElementById('map'), {
		zoom: 7,
		center: posCenter,
		zoomControl: true,
		mapTypeControl: false,
		scaleControl: true,
		streetViewControl: false,
		rotateControl: false,
		showscale: true,
		disableDefaultUI: true
	});
	
	var startCircle = new google.maps.Circle({
		strokeColor: '#00FF00',
		strokeOpacity: 0.8,
		strokeWeight: 2,
		fillColor: '#00FF00',
		fillOpacity: 0.35,
		map: map,
		center: pos1,
		radius: 50000,
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
		radius: 50000,
    	geodesic: true
    });
}

    </script>
    <script src="https://maps.googleapis.com/maps/api/js?key={{.GoogleMapsApiCredentials}}&signed_in=true&callback=initMap"
        async defer></script>
  </body>
</html>
`

type dataSelectionDisplayArgs struct {
	GoogleMapsApiCredentials string
}

func DisplayDataSelection(w http.ResponseWriter, r *http.Request, googleMapsApiCredentials string) {
	arguments := dataSelectionDisplayArgs{
		GoogleMapsApiCredentials: googleMapsApiCredentials,
	}
	
	t1 := template.New("Data selection template")
	t3, err := t1.Parse(selectionHtmlData)
	utils.CheckErr(err, "problem while parsing template")

	err = t3.Execute(w, arguments)
	utils.CheckErr(err, "problem while executing template")
}
