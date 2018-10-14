package igcsdb

import (
	"fmt"
	"time"
)

// time variable that have i
var start time.Time

// ServerStart starting server time when a ipaddress is set
func ServerStart() {
	// needs to be exported for use in main
	start = time.Now()
}

func timeFormatter() string {

	// uptime is being formatted to iso8601 style inform of duration from wiki
	// for now its only days hours and miniutes and seconds since my server wont probally be on for years

	elapsed := time.Since(start)
	uptime := fmt.Sprintf("P%dD%dH%dM%dS", int(elapsed.Hours())/24, int(elapsed.Hours())%24, int(elapsed.Minutes())%60, int(elapsed.Seconds())%60)

	return uptime
}
