package elastic

import (
	"context"
	"crypto/tls"
	"fpg_types"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	log "github.com/sirupsen/logrus"
)

var es *elasticsearch.TypedClient

func Connect() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Setup a default elastic connection, TODO this should be cli stuff
	var err error
	es, err = elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{
			"https://localhost:9200",
		},
		Username: "elastic",
		Password: "changeme",
	})
	if err != nil {
		log.Error(err)
	}
}

func Write(point *fpg_types.FlightLoc) {
	res, err := es.Index("locs").Request(point).Do(context.Background())

	if err != nil {
		log.Error(err)
		log.Error(res)
	}
}
