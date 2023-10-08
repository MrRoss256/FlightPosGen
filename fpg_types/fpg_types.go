package fpg_types

import (
	"strings"
	"time"
)

type TimeNoDate time.Time

// Structs for JSON formatted as:
//
//	{
//	   "flights": [
//	     {
//	       "flightdes": "AA101",
//	       "dep": { "code": "LHR", "time": "08:00:00", "long": 51.4706, "lat": -0.461941 , "tz": "Europe/London"},
//	       "arr": { "code": "SFO", "time": "23:00:00", "long": 37.37, "lat": -122.375, "tz" : "America/California"}
//	     },
type Loc struct {
	Code string     `json:"code"`
	Time TimeNoDate `json:"time"`
	Long float64    `json:"long"`
	Lat  float64    `json:"lat"`
	Tz   string     `json:"tz"`
}

type Flight struct {
	FlightDes string `json:"flightdes"`
	Dep       Loc    `json:"dep"`
	Arr       Loc    `json:"arr"`
	Speed     int64
	Duration  int64
}

type Flights struct {
	Flights []Flight `json:"flights"`
}

// Custom JSON unmarshaller for time with no date (UTC),
// formatted YY:MM:SS.
func (t *TimeNoDate) UnmarshalJSON(data []byte) error {
	value := strings.Trim(string(data), `"`)
	if value == "" || value == "null" {
		return nil
	}

	i, err := time.Parse("15:04:05", value)
	if err != nil {
		return err
	}
	*t = TimeNoDate(i)
	return nil
}

// Struct for output JSON.
//
//	{
//		"flightdes": "AA101",
//		"deptime": "2023-09-20T08:00:00Z",
//		"eventtime": "2023-09-20T08:16:00Z",
//		"location": {
//		  "lon": -2.0211443935216784,
//		  "lat": 52.45964431712232
//		}
type Location struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type FlightLoc struct {
	FlightDes string    `json:"flightdes"`
	DepTime   time.Time `json:"deptime"`
	EventTime time.Time `json:"eventtime"`
	Location  Location  `json:"location"`
}
