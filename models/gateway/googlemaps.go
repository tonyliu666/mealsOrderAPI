package gateway

import (
	"context"
	"log"
	"os"

	"googlemaps.github.io/maps"
)

func check(err error) {
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
}

var googleMapsClient *maps.Client

func Init() error {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	var err error
	googleMapsClient, err = maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return err
	}
	return nil
}
func GetCurrentLocation(address string) ([]maps.LatLng, error) {
	var locationCollection []maps.LatLng
	r := &maps.GeocodingRequest{
		Address: address,
	}

	resp, err := googleMapsClient.Geocode(context.Background(), r)
	if err != nil {
		return nil, err
	}

	for _, location := range resp {
		locationCollection = append(locationCollection, location.Geometry.Location)
	}

	return locationCollection, nil
}

func NearbySearch(radius uint, location maps.LatLng, keyword []string) ([]maps.PlacesSearchResponse, error) {
	var resp []maps.PlacesSearchResponse
	log.Println("keyword", keyword[0], radius, location)
	for _, key := range keyword {
		r := &maps.NearbySearchRequest{
			Radius:  radius,
			Keyword: key,
			OpenNow: false,
		}

		r.RankBy = maps.RankByDistance
		r.Location = &location

		rsp, err := googleMapsClient.NearbySearch(context.Background(), r)
		if err != nil {
			return nil, err
		}
		log.Println("result", rsp, "for", key)
		resp = append(resp, rsp)

	}
	return resp, nil
}
