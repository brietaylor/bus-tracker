package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/alecthomas/kingpin/v2"
	gtfs "github.com/brietaylor/online-bus-tracker/proto"
	"google.golang.org/protobuf/proto"
)

const (
	apiKey     = "***REMOVED***"
	staticData = "gtfs-static/translink-2024-11-01/"
)

func getLiveData() (*gtfs.FeedMessage, error) {
	v := url.Values{}
	v.Add("apikey", apiKey)
	url := url.URL{
		Scheme:   "https",
		Host:     "gtfsapi.translink.ca",
		Path:     "/v3/gtfsposition",
		RawQuery: v.Encode(),
	}
	res, err := http.Get(url.String())
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer res.Body.Close()

	log.Printf("API returned %d\n", res.StatusCode)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	log.Printf("Read %d bytes\n", len(body))

	pb := gtfs.FeedMessage{}
	err = proto.Unmarshal(body, &pb)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling pb: %w", err)
	}

	return &pb, nil
}

type handler struct {
	routes Routes
}

func (h *handler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Error handling request at %s: %s",
		r.URL.Path, err)

	w.WriteHeader(500)
	io.WriteString(w, "Internal server error")
}

func (h *handler) handleGetRoutes(w http.ResponseWriter, r *http.Request) {
	type respRoute struct {
		ShortName string `json:"short_name"`
		LongName  string `json:"long_name"`
	}

	routes := make([]respRoute, 0, len(h.routes))
	for _, route := range h.routes {
		routes = append(routes, respRoute{
			ShortName: route.RouteShortName,
			LongName:  route.RouteLongName,
		})
	}

	slices.SortFunc(routes, func(a, b respRoute) int {
		return strings.Compare(a.ShortName, b.ShortName)
	})

	respBody, err := json.Marshal(routes)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.Header().Add("access-control-allow-origin", "*")
	w.WriteHeader(200)
	w.Write(respBody)
}

func (h *handler) handleGetVehicles(w http.ResponseWriter, r *http.Request) {
	liveData, err := getLiveData()
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	routeSelection := query.Get("route")

	type vehicle struct {
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		RouteShortName string  `json:"route_short_name"`
		RouteLongName  string  `json:"route_long_name"`
		Direction      string  `json:"direction"`
	}

	type respJSON struct {
		Vehicles []vehicle `json:"vehicles"`
	}

	resp := respJSON{}

	// Convert live feed to JSON response
	for _, entity := range liveData.Entity {
		routeId := *entity.Vehicle.Trip.RouteId

		// Find the route name by picking the route ID from the static data.
		route, ok := h.routes[routeId]
		if !ok {
			// There's no matching route between the live and static data.
			// Probably Translink has added a new route and we need to update
			// the static data.  For now let's ignore it.
			continue
		}

		if routeSelection != "" && route.RouteShortName != routeSelection {
			continue
		}

		// Find the direction name
		directionID := *entity.Vehicle.Trip.DirectionId
		direction, ok := route.Directions[directionID]
		if !ok {
			h.handleError(w, r,
				fmt.Errorf("couldn't find matching direction for route id %s, direction id %d",
					routeId, directionID))
			return
		}

		vehicle := vehicle{
			Lat:            float64(*entity.Vehicle.Position.Latitude),
			Lon:            float64(*entity.Vehicle.Position.Longitude),
			RouteShortName: route.RouteShortName,
			RouteLongName:  route.RouteLongName,
			Direction:      direction.Text,
		}
		resp.Vehicles = append(resp.Vehicles, vehicle)
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.Header().Add("access-control-allow-origin", "*")
	w.WriteHeader(200)
	w.Write(respBody)
}

var (
	listenAddr = kingpin.Flag("listen", "").Default("localhost:8080").String()
)

func main() {
	kingpin.Parse()

	routes, err := readRoutes("gtfs-static/translink-2024-11-01/")
	if err != nil {
		log.Fatalf("Failed to read routes: %s", err)
	}

	h := handler{
		routes: routes,
	}
	http.HandleFunc("/getVehicles", h.handleGetVehicles)
	http.HandleFunc("/getRoutes", h.handleGetRoutes)

	log.Printf("Listening on %s", *listenAddr)
	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		log.Fatalf("Failed to start http server: %s", err)
	}
}
