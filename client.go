package humphrey

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
)

// A Client lets you use the default BART client or your own API Key
type Client struct {
	APIKey string
}

// MakeRequest is the heart of the client library
func (c *Client) MakeRequest(path string, cmd string, params url.Values, v interface{}) error {
	params["key"], params["cmd"] = []string{c.APIKey}, []string{cmd}
	requestURL := url.URL{Scheme: "http", Host: "api.bart.gov", Path: fmt.Sprintf("/api/%s.aspx", path), RawQuery: params.Encode()}
	response, err := http.Get(requestURL.String())
	if err != nil {
		return fmt.Errorf("requesting \"%s\": %v", requestURL.String(), err)
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("received non-200 status code (%d) requesting \"%s\"", response.StatusCode, requestURL.String())
	}
	defer response.Body.Close()
	xmlDecoder := xml.NewDecoder(response.Body)
	err = xmlDecoder.Decode(v)
	if err != nil {
		return fmt.Errorf("decoding the response: %v", err)
	}
	return nil
}

// DefaultClient is the BART Client with the default key
var DefaultClient = Client{"MW9S-E7SL-26DU-VV8V"}
