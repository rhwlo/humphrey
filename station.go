package humphrey

import (
	"encoding/xml"
	"fmt"
	"net/url"
)

// Station stubs out the BART Station interface
type Station struct {
	Abbreviation string  `xml:"abbr"`
	Name         string  `xml:"name"`
	Latitude     float32 `xml:"gtfs_latitude"`
	Longitude    float32 `xml:"gtfs_longitude"`
	Address      string  `xml:"address"`
	City         string  `xml:"city"`
	County       string  `xml:"county"`
	State        string  `xml:"state"`
	ZIPCode      string  `xml:"zipcode"`
}

type stationsResponse struct {
	XMLName  xml.Name  `xml:"root"`
	Stations []Station `xml:"stations>station"`
	Message  string    `xml:"message"`
}

// GetAllStations gets all the BART stations
func (c *Client) GetAllStations() ([]Station, error) {
	var response stationsResponse
	err := c.MakeRequest("stn", "stns", url.Values{}, &response)
	if err != nil {
		return nil, fmt.Errorf("making 'stns' request: %v", err)
	}
	return response.Stations, nil
}
