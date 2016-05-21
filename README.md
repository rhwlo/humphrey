**humphrey** is a [Go][] library written for [BART][]’s API and named after the [award-winning free BART transportation to Berkeley of yesteryear][history].

```go
// an example
import (
	"fmt"
	"log"
	"strings"

	"github.com/rhwlo/humphrey"
)

func main() {
	stationMap := make(map[string]humphrey.Station)
	routeMap := make(map[string]humphrey.Route)

	routes, err := humphrey.DefaultClient.GetCurrentRoutes()
	if err != nil {
		log.Fatal(err)
	}
	for _, route := range routes {
		routeMap[route.ID] = route
	}

	stations, err := humphrey.DefaultClient.GetAllStations()
	if err != nil {
		log.Fatal(err)
	}
	for _, station := range stations {
		stationMap[station.Abbreviation] = station
	}

	trains, err := humphrey.DefaultClient.GetStationSchedule(stationMap["MCAR"])
	if err != nil {
		log.Fatal(err)
	}
	for _, train := range trains {
		stationNames := strings.Split(routeMap[train.Line].Abbreviation, "-")
		sourceStation := stationMap[stationNames[0]]
		destStation := stationMap[stationNames[1]]
		trainID := fmt.Sprintf("#%d on %s", train.Index, strings.ToLower(train.Line))
		fmt.Printf("Train %s (headed from %s to %s) is scheduled to pass %s at %s\n", trainID, sourceStation.Name, destStation.Name, stationMap["MCAR"].Name, train.Time.Format("3:04 PM"))
	}
}
```

## still to be done:

- I’d like to wrap the [`sched/load`][] call, but it looks a little tricky
- I’d like to wrap the [`sched/special`][] and [`sched/holiday`][] calls into one function called something like `Client.GetScheduleIrregularities()`
- I might want to wrap the calls for [`sched/arrive`][] and [`sched/depart`][], though [it sounds (by BART’s own admission) like BART deals with time in a creative way][0].

I deliberately left out the [`stn/stninfo`][] and [`stn/stnaccess`][] calls because most of the information they provide is either provided elsewhere or is in a format intended to be presented through a web browser (i.e., formatted HTML).

[0]: http://api.bart.gov/docs/overview/barttime.aspx
[`sched/arrive`]: http://api.bart.gov/docs/sched/arrive.aspx
[`sched/depart`]: http://api.bart.gov/docs/sched/depart.aspx
[`sched/load`]: http://api.bart.gov/docs/sched/load.aspx
[`sched/holiday`]: http://api.bart.gov/docs/sched/holiday.aspx
[`sched/special`]: http://api.bart.gov/docs/sched/special.aspx
[`stn/stnaccess`]: http://api.bart.gov/docs/stn/stnaccess.aspx
[`stn/stninfo`]: http://api.bart.gov/docs/stn/stninfo.aspx
[history]: http://web.archive.org/web/20030325171434/http://www.mtc.ca.gov/whats_happening/awards/1977.htm
[Go]: https://golang.org/
[BART]: http://bart.gov/
