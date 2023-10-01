# FlightPosGen
Generate Flight Positions given a flight's origin and destination 


# Installing
Build with golang, run the executable.




# Running

The script runs in one of two modes, "Perpetual" mode the script executes every 30 seconds and generates position information for the flights that are currently active.
In "Batch" mode the script generates flight poition info for flights from "(now -24 hours) to now" without any delay.  

The script generates a fixed JSON document for each flight position, writing to the console or Elasticsearch depending on the command line options.

```
$ ./fpg --help
Usage of ./fpg:
  -console
    	Output results to the console (default true)
  -elastic
    	Output results to Elasticsearch (localhost)
  -perpetual
    	Generate positions every few seconds, until the end of time (default is for the previous 24 hours).
```

## flights.json

The `flights.json` lists the flights, the flight departure, the from-to's geo-location and the flight arrival.

```json
 {
    "flights": [
      {
        "flightdes": "AA101",
        "dep": { "code": "LHR", "time": "08:00:00", "long": 51.4706, "lat": -0.461941 },
        "arr": { "code": "SFO", "time": "23:00:00", "long": 37.37, "lat": -122.375 }
      }
    ]
 }
```

# Disclaimer

The script is written as a experiment to see the feasibility of using such a tool to create test data, and as such it's not production ready, but with a little effort caould be.

The script has a few issues:
- Lots of hard-coding, assumes Elasticsearch is running locally.
- Will potentially create duplicate records if run multiple times.
- The JSON file is simialr to a flight schedule, using local times with no dates.
  However flights operate across dates (either forwards or backeards), so the flight arrival time needs a [+-][123].
- There is no batching of data being pushed to Elasticsearch, it might not be very fast.

# Credita

The script is a simple wrapper around the `geographiclib` library credit due to the authors and porters of that library!
