package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Kamila3820/hoca-backend/config"
	"github.com/sirupsen/logrus"
)

type Client struct {
	APIKey string
}

type DirectionsRequest struct {
	Origin        string
	Destination   string
	Mode          string
	Waypoints     []string
	DepartureTime string
	ArrivalTime   string
}

type DirectionsResponse struct {
	Routes            []Route            `json:"routes"`
	GeocodedWaypoints []GeocodedWaypoint `json:"geocoded_waypoints"`
	Status            string             `json:"status"`
}

type Route struct {
	Summary string `json:"summary"`
	Legs    []Leg  `json:"legs"`
}

type GeocodedWaypoint struct {
	GeocoderStatus string   `json:"geocoder_status"`
	PlaceID        string   `json:"place_id"`
	Types          []string `json:"types"`
}

type Leg struct {
	Distance Distance `json:"distance"`
	Duration Duration `json:"duration"`
}

type Distance struct {
	Text string `json:"text"`
}

type Duration struct {
	Text string `json:"text"`
}

type Polyline struct {
	Points string `json:"points"`
}

func (c *Client) Directions(r *DirectionsRequest) (*DirectionsResponse, error) {
	baseURL := "https://maps.googleapis.com/maps/api/directions/json"

	params := url.Values{}
	params.Add("origin", r.Origin)
	params.Add("destination", r.Destination)
	params.Add("key", config.ConfigGetting().Google.ApiKey)

	if r.Mode != "" {
		params.Add("mode", r.Mode)
	}
	if len(r.Waypoints) > 0 {
		params.Add("waypoints", url.QueryEscape(strings.Join(r.Waypoints, "|")))
	}
	if r.DepartureTime != "" {
		params.Add("departure_time", r.DepartureTime)
	}
	if r.ArrivalTime != "" {
		params.Add("arrival_time", r.ArrivalTime)
	}

	renderURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	fmt.Println(renderURL)

	req, err := http.NewRequest("GET", renderURL, nil)
	if err != nil {
		logrus.Error("Unable to construct request ", err)
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error("Unable to send request ", err)
		panic(err)
	}
	defer res.Body.Close()
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		logrus.Error("Unable to read response body ", err)
		panic(err)
	}

	var directionsResponse DirectionsResponse
	if err := json.Unmarshal(body, &directionsResponse); err != nil {
		logrus.Error("Unable to parse response body ", err)
		panic(err)
	}

	return &directionsResponse, nil
}
