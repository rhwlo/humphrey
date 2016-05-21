package bart

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// WEEKDAY, SATURDAY, SUNDAY, TODAY, and NOW are all valid dates to pass to GetRouteSchedule
const (
	WEEKDAY  = "wd"
	SATURDAY = "sd"
	SUNDAY   = "su"
	TODAY    = "today"
	NOW      = "now"
)

// A RouteSchedule represents the BART XML document that starts with a <root>
type RouteSchedule struct {
	XMLName xml.Name        `xml:"root"`
	Number  int             `xml:"sched_num"`
	Trains  []TrainSchedule `xml:"route>train"`
	Message string          `xml:"message>special_schedule"`
}

// A Schedule represents the schedule entry passed from the BART `scheds` command.
type Schedule struct {
	EffectiveDate time.Time
	Number        int
}

type intermediateSchedule struct {
	XMLName       xml.Name `xml:"schedule"`
	EffectiveDate string   `xml:"effectivedate,attr"`
	Number        int      `xml:"id,attr"`
}

// a bartSched represents the BART `scheds` command XML response
type bartSchedDocument struct {
	XMLName   xml.Name   `xml:"root"`
	Schedules []Schedule `xml:"schedules>schedule"`
}

// A TrainSchedule represents the BART XML <train>
type TrainSchedule struct {
	XMLName xml.Name        `xml:"train"`
	Index   int             `xml:"index,attr"`
	Stops   []ScheduledStop `xml:"stop"`
}

// A ScheduledStop represents the BART XML <stop>
type ScheduledStop struct {
	XMLName  xml.Name `xml:"stop"`
	Station  string   `xml:"station,attr"`
	Time     string   `xml:"origTime,attr"`
	BikeFlag string   `xml:"bikeFlag,attr"`
}

// UnmarshalXML unmarshals the Schedule by way of an intermediate struct
func (s *Schedule) UnmarshalXML(xmlDecoder *xml.Decoder, startEl xml.StartElement) error {
	var i intermediateSchedule
	err := xmlDecoder.DecodeElement(&i, &startEl)
	if err != nil {
		return err
	}
	s.Number = i.Number
	s.EffectiveDate, err = time.Parse("01/02/2006 03:04 PM", i.EffectiveDate)
	if err != nil {
		return fmt.Errorf("parsing a time for schedule %d: %v", s.Number, err)
	}
	return nil
}

// GetRouteScheduleByDate requests a route schedule by route number and date, which can be formatted
// either as a "MM/DD/YYYY" string or passed as one of the constants WEEKDAY, SATURDAY, SUNDAY, TODAY, or NOW.
func (c *Client) GetRouteScheduleByDate(routeNumber int, date string) (RouteSchedule, error) {
	var routeSchedule RouteSchedule
	params := url.Values{"route": []string{strconv.Itoa(routeNumber)}, "date": []string{date}, "cmd": []string{"routesched"}}
	err := c.MakeRequest("sched", "routesched", params, &routeSchedule)
	if err != nil {
		return RouteSchedule{}, fmt.Errorf("making 'routesched' request for %d on %v: %v", routeNumber, date, err)
	}
	return routeSchedule, nil
}

// GetRouteScheduleBySchedule requests a route schedule by route number and schedule number.
func (c *Client) GetRouteScheduleBySchedule(routeNumber int, scheduleNumber int) (RouteSchedule, error) {
	var routeSchedule RouteSchedule
	params := url.Values{"route": []string{strconv.Itoa(routeNumber)}, "sched": []string{strconv.Itoa(scheduleNumber)}, "cmd": []string{"routesched"}}
	err := c.MakeRequest("sched", "routesched", params, &routeSchedule)
	if err != nil {
		return RouteSchedule{}, fmt.Errorf("making 'routesched' request for %d on schedule #%d: %v", routeNumber, scheduleNumber, err)
	}
	return routeSchedule, nil
}

// GetCurrentSchedules requests the information about currently relevant schedules.
func (c *Client) GetCurrentSchedules() ([]Schedule, error) {
	var schedResponse bartSchedDocument
	err := c.MakeRequest("sched", "scheds", url.Values{}, &schedResponse)
	if err != nil {
		return nil, fmt.Errorf("making 'scheds' request: %v", err)
	}
	return schedResponse.Schedules, nil
}

// GetAllRouteSchedulesBySchedules returns all routes for the slice of Schedules given.
func (c *Client) GetAllRouteSchedulesBySchedules(schedules []Schedule) ([]RouteSchedule, error) {
	var routeSchedules []RouteSchedule
	for _, schedule := range schedules {
		routes, err := c.GetRoutesBySchedule(schedule)
		if err != nil {
			return nil, fmt.Errorf("getting routes for schedule #%d: %v", schedule.Number, err)
		}
		for _, route := range routes {
			routeSchedule, err := c.GetRouteScheduleBySchedule(route.Number, schedule.Number)
			if err != nil {
				return nil, fmt.Errorf("getting a route schedule for schedule #%d: %v", schedule.Number, err)
			}
			routeSchedules = append(routeSchedules, routeSchedule)
		}
	}
	return routeSchedules, nil
}
