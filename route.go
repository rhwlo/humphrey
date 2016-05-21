package bart

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strconv"
)

// Route represents a BART route
type Route struct {
	Name         string
	Abbreviation string
	ID           string
	Number       int
	Color        string
	Holidays     int
	StationCount int
	Stations     []Station
}

type intermediateRoute struct {
	XMLName      xml.Name `xml:"route"`
	Name         string   `xml:"name"`
	Abbreviation string   `xml:"abbr"`
	ID           string   `xml:"routeID"`
	Number       int      `xml:"number"`
	Color        string   `xml:"color"`
	Holidays     int      `xml:"holidays"` // I have no idea what this one means
	StationCount int      `xml:"num_stns"` // This should match len(Route.Stations) but who knows!
	Stations     []string `xml:"config>station"`
}

type routesResponse struct {
	XMLName xml.Name `xml:"root"`
	Routes  []Route  `xml:"routes>route"`
}

// UnmarshalXML unmarshals the XML for the Route from the intermediateRoute
func (r *Route) UnmarshalXML(xmlDecoder *xml.Decoder, startEl xml.StartElement) error {
	var i intermediateRoute
	err := xmlDecoder.DecodeElement(&i, &startEl)
	if err != nil {
		return err
	}
	r.Name, r.Abbreviation, r.ID, r.Number, r.Color, r.Holidays, r.StationCount = i.Name, i.Abbreviation, i.ID, i.Number, i.Color, i.Holidays, i.StationCount
	for _, stationName := range i.Stations {
		r.Stations = append(r.Stations, Station{Name: stationName})
	}
	return nil
}

// GetCurrentRoutes returns a slice of routes currently active
func (c *Client) GetCurrentRoutes() ([]Route, error) {
	var response routesResponse
	err := c.MakeRequest("route", "routes", url.Values{}, &response)
	if err != nil {
		return nil, fmt.Errorf("making 'routes' request: %v", err)
	}
	return response.Routes, nil
}

// GetRoutesBySchedule returns a slice of routes by schedule number.
func (c *Client) GetRoutesBySchedule(schedule Schedule) ([]Route, error) {
	var response routesResponse
	params := url.Values{"sched": []string{strconv.Itoa(schedule.Number)}}
	err := c.MakeRequest("route", "routes", params, &response)
	if err != nil {
		return nil, fmt.Errorf("making 'routes' request for schedule #%d: %v", schedule.Number, err)
	}
	return response.Routes, nil
}

// GetRouteInfoByNumber takes a route number and returns a Route with all the info
func (c *Client) GetRouteInfoByNumber(routeNumber int) (Route, error) {
	var response routesResponse
	params := url.Values{"route": []string{strconv.Itoa(routeNumber)}}
	err := c.MakeRequest("route", "routeinfo", params, &response)
	if err != nil {
		return Route{}, fmt.Errorf("making 'routeinfo' request for route #%d: %v", routeNumber, err)
	}
	return response.Routes[0], nil
}

// GetAllRouteInfo returns a slice of all Routes with all the info
func (c *Client) GetAllRouteInfo() ([]Route, error) {
	var response routesResponse
	params := url.Values{"route": []string{"all"}}
	err := c.MakeRequest("route", "routeinfo", params, &response)
	if err != nil {
		return nil, fmt.Errorf("making 'routeinfo' request for all routes: %v", err)
	}
	return response.Routes, nil
}
