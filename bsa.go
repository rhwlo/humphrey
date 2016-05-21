package humphrey

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"time"
)

// ServiceAdvisory represents the BART service advisories
type ServiceAdvisory struct {
	ID          int
	Type        string
	Description string
	SMSText     string
	Posted      time.Time
	Expires     time.Time
}

type intermediateServiceAdvisory struct {
	XMLName xml.Name `xml:"bsa"`
	ID      int      `xml:"id,attr"`
	//	Station     string   `xml:"station"` // this isn't used yet
	Description string `xml:"description"`
	SMSText     string `xml:"sms_text"`
	Type        string `xml:"type"`
	Posted      string `xml:"posted"`
	Expires     string `xml:"expires"`
}

type serviceAdvisoryResponse struct {
	XMLName    xml.Name          `xml:"root"`
	Advisories []ServiceAdvisory `xml:"bsa"`
}

type trainCountResponse struct {
	XMLName    xml.Name `xml:"root"`
	TrainCount int      `xml:"traincount"`
}

const (
	postedFormat             = "Mon Jan 02 2006 03:04 PM MST"
	expiresFormat            = "Mon Jan 02 2006 03:04 PM MST"
	noDelaysDescription      = "No delays reported."
	elevatorsFineDescription = "Attention passengers: All elevators are in service. Thank You." // maybe? this is untested
)

// UnmarshalXML unmarshals the ServiceAdvisory by way of an intermediate type
func (s *ServiceAdvisory) UnmarshalXML(xmlDecoder *xml.Decoder, startElement xml.StartElement) error {
	var i intermediateServiceAdvisory
	err := xmlDecoder.DecodeElement(&i, &startElement)
	if err != nil {
		return err
	}
	s.ID, s.Description, s.SMSText = i.ID, i.Description, i.SMSText
	if i.Posted == "" {
		// An unexpectedly shaped blob of XML gets handed to us if no delays exist.
		// We return early here because otherwise it gets ugly
		return nil
	}
	switch i.Type {
	case "DELAY", "ELEVATOR", "EMERGENCY":
		s.Type = i.Type
	default:
		return fmt.Errorf("%v is not a valid alert type (must be \"DELAY\" or \"EMERGENCY\")", i.Type)
	}
	s.Posted, err = time.Parse(postedFormat, i.Posted)
	if err != nil {
		return fmt.Errorf("parsing <posted> time: %v", err)
	}
	s.Expires, err = time.Parse(expiresFormat, i.Expires)
	if err != nil {
		return fmt.Errorf("parsing <expires> time: %v", err)
	}
	return nil
}

// GetServiceAdvisories uses the `bsa` command to fetch all service advisories currently posted by BART.
func (c *Client) GetServiceAdvisories() ([]ServiceAdvisory, error) {
	var response serviceAdvisoryResponse
	err := c.MakeRequest("bsa", "bsa", url.Values{}, &response)
	if err != nil {
		return nil, fmt.Errorf("making 'bsa' request: %v", err)
	}
	if response.Advisories[0].Description == noDelaysDescription {
		// An unexpectedly shaped blob of XML gets handed to us if no delays exist
		return []ServiceAdvisory{}, nil
	}
	return response.Advisories, nil
}

// GetElevatorAdvisories uses the `elev` command to fetch all service advisories currently posted by BART.
func (c *Client) GetElevatorAdvisories() ([]ServiceAdvisory, error) {
	var response serviceAdvisoryResponse
	err := c.MakeRequest("bsa", "elev", url.Values{}, &response)
	if err != nil {
		return nil, fmt.Errorf("making 'elev' request: %v", err)
	}
	if response.Advisories[0].Description == elevatorsFineDescription {
		return []ServiceAdvisory{}, nil
	}
	return response.Advisories, nil
}

// GetTrainCount uses the `count` command to count all the trains currently in service.
func (c *Client) GetTrainCount() (int, error) {
	var response trainCountResponse
	err := c.MakeRequest("bsa", "count", url.Values{}, &response)
	if err != nil {
		return 0, fmt.Errorf("making 'count' request: %v", err)
	}
	return response.TrainCount, nil
}
