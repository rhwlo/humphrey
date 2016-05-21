package humphrey

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strconv"
)

// Estimate represents a BART arrival estimate
type Estimate struct {
	Minutes   int
	Platform  int
	Direction string
	Length    int
	Color     string
	HexColor  string
	BikeFlag  string
}

type intermediateEstimate struct {
	XMLName   xml.Name `xml:"estimate"`
	Minutes   string   `xml:"minutes"`
	Platform  int      `xml:"platform"`
	Direction string   `xml:"direction"`
	Length    int      `xml:"length"`
	Color     string   `xml:"color"`
	HexColor  string   `xml:"hexcolor"`
	BikeFlag  string   `xml:"bikeflag"`
}

type etdDocument struct {
	XMLName    xml.Name    `xml:"root"`
	Departures []Departure `xml:"station>etd"`
}

type intermediateDeparture struct {
	XMLName                 xml.Name   `xml:"etd"`
	Destination             string     `xml:"destination"`
	DestinationAbbreviation string     `xml:"abbreviation"`
	Estimates               []Estimate `xml:"estimate"`
	Limited                 int        `xml:"limited"`
}

// Departure represents a deprature from a BART station
type Departure struct {
	Destination Station
	Limited     bool
	Estimates   []Estimate
}

// UnmarshalXML unmarshals XML onto the Departure object
func (d *Departure) UnmarshalXML(xmlDecoder *xml.Decoder, startEl xml.StartElement) error {
	var i intermediateDeparture
	err := xmlDecoder.DecodeElement(&i, &startEl)
	if err != nil {
		return err
	}
	d.Destination = Station{Abbreviation: i.DestinationAbbreviation, Name: i.Destination}
	d.Estimates = i.Estimates
	d.Limited = (i.Limited == 1)
	return nil
}

// UnmarshalXML unmarshals XML onto the Estimate object
func (e *Estimate) UnmarshalXML(xmlDecoder *xml.Decoder, startEl xml.StartElement) error {
	var i intermediateEstimate
	err := xmlDecoder.DecodeElement(&i, &startEl)
	if err != nil {
		return err
	}
	minutes, err := strconv.Atoi(i.Minutes)
	if err != nil {
		// BART likes to use _words_ when your train is at the station
		minutes = 0
	}
	e.Minutes = minutes
	e.Platform = i.Platform
	e.Direction = i.Direction
	e.Length = i.Length
	e.Color = i.Color
	e.HexColor = i.HexColor
	e.BikeFlag = i.BikeFlag
	return nil
}

// GetDeparturesByStation gets the estimates for a station that is passed to it
func (c *Client) GetDeparturesByStation(station Station) ([]Departure, error) {
	var response etdDocument
	err := c.MakeRequest("etd", "etd", url.Values{"orig": []string{station.Abbreviation}}, &response)
	if err != nil {
		return nil, fmt.Errorf("making 'etd' request for station %s: %v", station.Name, err)
	}
	return response.Departures, nil
}
