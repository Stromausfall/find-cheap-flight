package utils

import (
	"fmt"
	"runtime"
	"os/exec"
	"time"
)

// returns the date in format YYYY-MM-DD
func DateToString(date time.Time) string {
	return date.Format("2006-01-02")
}

func DateFromString(date string) time.Time {
	result, err := time.Parse("2006-01-02", date)
	CheckErr(err, "unable to parse date !")
	
	return result
}

func DateStringNow(offsetInDays int) string {
	duration := time.Hour * 24 * time.Duration(int64(offsetInDays))
	
	dateToConvert := time.Now().Add(duration)
	
	return DateToString(dateToConvert)
}

func CheckErr(err error, info string) {
	if err != nil {
		message := fmt.Sprintf(info, " - ", err, " = ", err.Error())
		panic(message)
	}
}

// openURL opens a browser window to the specified location.
// This code originally appeared at:
//   http://stackoverflow.com/questions/10377243/how-can-i-launch-a-process-that-is-not-a-file-in-go
func OpenURL(url string) error {
        var err error
        switch runtime.GOOS {
        case "linux":
                err = exec.Command("xdg-open", url).Start()
        case "windows":
                err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
        case "darwin":
                err = exec.Command("open", url).Start()
        default:
                err = fmt.Errorf("Cannot open URL %s on this platform", url)
        }
        return err
}
