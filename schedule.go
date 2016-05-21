package humphrey

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

type stationScheduleResponse struct {
	XMLName xml.Name `xml:"root"`
	Trains  []Train  `xml:"station>item"`
}

// Train represents the train items in the response from `stnsched`
type Train struct {
	Line            string
	Origin          Station
	Time            time.Time
	DestinationTime time.Time
	Index           int
	BikeFlag        int
}

type intermediateTrain struct {
	XMLName         xml.Name `xml:"item"`
	Line            string   `xml:"line,attr"`
	Origin          string   `xml:"trainHeadStation,attr"`
	Time            string   `xml:"origTime,attr"`
	DestinationTime string   `xml:"destTime,attr"`
	Index           int      `xml:"trainIdx,attr"`
	BikeFlag        int      `xml:"bikeflag,attr"`
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
	BikeFlag int      `xml:"bikeFlag,attr"`
}

const (
	effectiveTimeFormat   = "01/02/2006 03:04 PM"
	originTimeFormat      = "3:04 PM"
	destinationTimeFormat = "3:04 PM"
)

// UnmarshalXML unmarshals the Schedule by way of an intermediate struct
func (s *Schedule) UnmarshalXML(xmlDecoder *xml.Decoder, startEl xml.StartElement) error {
	var i intermediateSchedule
	err := xmlDecoder.DecodeElement(&i, &startEl)
	if err != nil {
		return err
	}
	s.Number = i.Number
	s.EffectiveDate, err = time.Parse(effectiveTimeFormat, i.EffectiveDate)
	if err != nil {
		return fmt.Errorf("parsing <effectiveDate> %v: %v", i.EffectiveDate, err)
	}
	return nil
}

// UnmarshalXML unmarshals the Train by way of an intermediate struct
func (t *Train) UnmarshalXML(xmlDecoder *xml.Decoder, startEl xml.StartElement) error {
	var i intermediateTrain
	err := xmlDecoder.DecodeElement(&i, &startEl)
	if err != nil {
		return err
	}
	t.Time, err = time.Parse(originTimeFormat, i.Time)
	if err != nil {
		return fmt.Errorf("parsing origTime %v: %v", i.Time, err)
	}
	t.DestinationTime, err = time.Parse(originTimeFormat, i.DestinationTime)
	if err != nil {
		return fmt.Errorf("parsing destTime %v: %v", i.DestinationTime, err)
	}
	t.Line, t.Index, t.BikeFlag = i.Line, i.Index, i.BikeFlag
	t.Origin = Station{Abbreviation: i.Origin}
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

// GetStationSchedule returns a slice of TrainSchedules for the given station
func (c *Client) GetStationSchedule(station Station) ([]Train, error) {
	var response stationScheduleResponse
	err := c.MakeRequest("sched", "stnsched", url.Values{"orig": []string{station.Abbreviation}}, &response)
	if err != nil {
		return nil, fmt.Errorf("making 'stnsched' request for %s: %v", station.Abbreviation, err)
	}
	return response.Trains, nil
}
