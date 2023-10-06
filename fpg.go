package main

import (
	el "elastic"
	"encoding/json"
	"flag"
	"fmt"
	fpg_types "fpg_types"
	"io"
	"os"
	rmq "rabbit"
	"time"

	"github.com/pymaxion/geographiclib-go/geodesic"
	"github.com/pymaxion/geographiclib-go/geodesic/capabilities"
	log "github.com/sirupsen/logrus"
)

// Given a flight with a departure and arrival time, if the
// flight is active at the time give.  Calculate how far we have
// travelled since the departure time (based on pre-canned speed)
// and find the point for that distance.
func weThereYet(flt *fpg_types.Flight, dt *time.Time) *fpg_types.FlightLoc {

	var dep = time.Time(flt.Dep.Time)
	var arr = time.Time(flt.Arr.Time)
	var now, _ = time.Parse("15:04:05", dt.Format("15:04:05"))

	if now.Before(dep) || now.After(arr) {
		return nil
	}

	log.Infof("At [%s] %s is in the air\n", dt.String(), flt.FlightDes)

	// How long has the flight been in the air?
	var dur = now.Sub(dep)

	// How far have we travelled in that time?
	var dist = int64(dur.Seconds()) * flt.Speed
	l := geodesic.WGS84.InverseLine(flt.Dep.Long, flt.Dep.Lat, flt.Arr.Long, flt.Arr.Lat)
	r := l.PositionWithCapabilities(float64(dist), capabilities.Standard|capabilities.LongUnroll)

	var loc fpg_types.FlightLoc
	loc.FlightDes = flt.FlightDes
	t := dt
	loc.DepTime = time.Date(t.Year(), t.Month(), t.Day(), dep.Hour(), dep.Minute(), dep.Second(), 0, t.Location())
	loc.Location.Lat = r.Lat2
	loc.Location.Lon = r.Lon2
	loc.EventTime = *dt
	return &loc
}

func main() {
	log.SetLevel(log.InfoLevel)

	useElastic := flag.Bool("elastic", false, "Output results to Elasticsearch (localhost)")
	useRabbit := flag.Bool("rabbit", false, "Output results to RabbitMQ (localhost)")
	useConsole := flag.Bool("console", true, "Output results to the console")
	perpetual := flag.Bool("perpetual", false, "Generate positions every few seconds, until the end of time (default is for the previous 24 hours).")
	flag.Parse()

	if *useElastic {
		el.Connect()
	}

	if *useRabbit {
		rmq.Connect()
	}

	// Populate the flts with the json data, should be coming from the TDG
	// view of active flights
	var flts fpg_types.Flights
	jsonFile, err := os.Open("flights.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &flts)
	if err != nil {
		fmt.Println(err)
	}

	// Calculate some useful values for the flights, flight duration and speed
	log.Info("Parsed and enriched flights:")
	for i := range flts.Flights {
		rag := geodesic.WGS84.Inverse(flts.Flights[i].Dep.Long, flts.Flights[i].Dep.Lat, flts.Flights[i].Arr.Long, flts.Flights[i].Arr.Lat)
		flts.Flights[i].Duration = int64(time.Time(flts.Flights[i].Arr.Time).Sub(time.Time(flts.Flights[i].Dep.Time)).Seconds())
		flts.Flights[i].Speed = int64(int64(rag.S12) / flts.Flights[i].Duration)
		log.Infof("\t%#v\n", flts.Flights[i])
	}

	if *perpetual {
		for {
			t := time.Now()
			for _, element := range flts.Flights {
				loc := weThereYet(&element, &t)
				if loc != nil {
					if *useConsole == true {
						var b, _ = json.Marshal(loc)
						fmt.Println(string(b))
					}
					if *useElastic == true {
						el.Write(loc)
					}
					if *useRabbit == true {
						rmq.Write(loc)
					}
				}
			}
			time.Sleep(10 * time.Second)
		}
	} else {
		//
		end := time.Now()
		start := end.Add(-(24 * time.Hour))
		for d := start; d.After(end) == false; d = d.Add(1 * time.Minute) {
			for _, element := range flts.Flights {
				loc := weThereYet(&element, &d)
				if loc != nil {
					if *useConsole == true {
						var b, _ = json.Marshal(loc)
						fmt.Println(string(b))
					}
					if *useElastic == true {
						el.Write(loc)
					}
					if *useRabbit == true {
						rmq.Write(loc)
					}
				}
			}
		}
	}
}
