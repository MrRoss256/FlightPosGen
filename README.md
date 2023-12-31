# FlightPosGen
Generate a flights geopositional positions given a flight's departure time, origin and destination.

# Installing
Build with golang, run the executable.

# Running

The script runs in one of two modes, "Perpetual" mode the script executes every 30 seconds and generates position information for the flights that are currently active.
In "Batch" mode the script generates flight poition info for flights from "(now -24 hours) to now" without any delay.  

The script generates a fixed JSON document for each flight position, writing to the console, RabbitMQ (AMQP) or Elasticsearch depending on the command line options.

```
Usage of ./fpg:
  -console
    	Output results to the console (default true)
  -elastic
    	Output results to Elasticsearch (localhost)
  -perpetual
    	Generate positions every few seconds, until the end of time (default is for the previous 24 hours).
  -rabbit
    	Output results to RabbitMQ (localhost)
```

## flights.json

The `flights.json` lists the flights, the flight departure, the from-to's geo-location and the flight arrival.

```json
 {
"flights": [
	     {
	       "flightdes": "AA101",
	       "dep": { "code": "LHR", "time": "08:00:00", "long": 51.4706, "lat": -0.461941 , "tz": "Europe/London"},
	       "arr": { "code": "SFO", "time": "23:00:00", "long": 37.37, "lat": -122.375, "tz" : "America/California"}
	     }
    ]
 }
```

# Disclaimer

The script is written as a experiment to see the feasibility of using such a tool to create test data, and as such it's not production ready, but with a little effort caould be.

The script has a few issues:
- Lots of hard-coding, assumes Elasticsearch or RabbitMQ is running locally.
- Will create duplicate records if run multiple times.
- The JSON file is simialr to a flight schedule, using local times with no dates.
  However flights operate across dates (either forwards or backwards), so the flight arrival time needs a [+-][123] adding.
- There is no batching of data being pushed to Elasticsearch, it might not be very fast.
- Terrible error handling.
- Time handling is not quite complete.

# Credits

The script is a simple wrapper around the `geographiclib` library credit due to the authors and porters of that library!
