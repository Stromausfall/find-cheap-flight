package calculate

import (
	"encoding/json"
	"fmt"
	"github.com/stromausfall/find-cheap-flight/data"
	"github.com/stromausfall/find-cheap-flight/utils"
	"net/http"
)

func DisplayPage(w http.ResponseWriter, r *http.Request, googleMapsApiCredentials string) {
	// we need this in order to get POST form data
	r.ParseMultipartForm(15485760)

	postFormData := []byte(r.PostFormValue("flightQueries"))

	flightQueries := make([]data.FlightQuery, 0)
	error := json.Unmarshal(postFormData, &flightQueries)
	utils.CheckErr(error, "Unable to cast received post form data to FlightQueries !")

	fmt.Println(flightQueries)
}
