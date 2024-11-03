package main

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/gocarina/gocsv"
)

type Routes map[string]*Route

type Direction struct {
	Text    string `csv:"direction"`
	ID      uint32 `csv:"direction_id"`
	RouteID string `csv:"route_id"`
}

type Route struct {
	ID             string `csv:"route_id"`
	AgencyId       string `csv:"agency_id"`
	RouteShortName string `csv:"route_short_name"`
	RouteLongName  string `csv:"route_long_name"`
	Directions     map[uint32]Direction
}

func readRoutes(path string) (Routes, error) {
	routesTxt := filepath.Join(path, "routes.txt")

	f, err := os.Open(routesTxt)
	if err != nil {
		return nil, fmt.Errorf("failed to open routes file %s: %w",
			routesTxt, err)
	}
	defer f.Close()

	routes := []*Route{}
	if err := gocsv.UnmarshalFile(f, &routes); err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}

	directionsTxt := filepath.Join(path, "directions.txt")

	f, err = os.Open(directionsTxt)
	if err != nil {
		return nil, fmt.Errorf("failed to open routes file %s: %w",
			routesTxt, err)
	}
	defer f.Close()

	directions := []*Direction{}
	if err := gocsv.UnmarshalFile(f, &directions); err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}

	ret := Routes{}
	for _, route := range routes {
		//log.Printf("Route ID %s", route.ID)
		route.Directions = make(map[uint32]Direction)
		for _, direction := range directions {
			if route.ID == direction.RouteID {
				//log.Printf("Direction ID %d", direction.ID)
				route.Directions[direction.ID] = *direction
			} else {
				//log.Printf("No match for route %s direction %s/%d",
				//route.ID, direction.RouteID, direction.ID)
			}
		}

		ret[route.ID] = route
	}
	return ret, nil
}
