package api

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/jakekeeys/bindicator/internal/collection"
	"github.com/sirupsen/logrus"
)

const (
	url = "https://secure.ashford.gov.uk/wastecollections/collectiondaylookup/"

	defaultIp   = "0.0.0.0"
	defaultPort = "8080"
)

func NewHTTP(ctx context.Context, debug bool) error {
	ip := os.Getenv("IP")
	if ip == "" {
		ip = defaultIp
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	logrus.Debugf("listening on %s:%s", ip, port)
	http.HandleFunc("/", serve(ctx, debug))
	return http.ListenAndServe(net.JoinHostPort(ip, port), nil)
}

func serve(ctx context.Context, debug bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("token") != os.Getenv("TOKEN") {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		postcode := strings.TrimSpace(r.URL.Query().Get("postcode"))
		if postcode == "" {
			http.Error(w, "postcode is required", http.StatusBadRequest)
			return
		}

		number := strings.TrimSpace(r.URL.Query().Get("number"))
		if number == "" {
			http.Error(w, "number is required", http.StatusBadRequest)
			return
		}

		logrus.Debug("attempting to retrieve next bin collection")
		collectionDates, err := collection.GetNext(r.Context(), debug, url, postcode, number)
		if err != nil {
			logrus.WithError(err).Error("could not get next collection")
			http.Error(w, "could not get next collection", http.StatusInternalServerError)
			return
		}

		respBytes, err := json.Marshal(collectionDates)
		if err != nil {
			logrus.WithError(err).Error("error marshalling response")
			http.Error(w, "error marshalling response", http.StatusInternalServerError)
			return
		}

		_, err = w.Write(respBytes)
		if err != nil {
			logrus.WithError(err).Error("error writing response")
			http.Error(w, "error writing response", http.StatusInternalServerError)
			return
		}
	}
}
