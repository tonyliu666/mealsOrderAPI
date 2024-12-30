package gateway

import (
	"context"
	"os"

	"googlemaps.github.io/maps"
)

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

func NearbySearch(location maps.LatLng, keyword []string) ([]map[string]interface{}, error) {
	var resp []map[string]interface{}

	for _, key := range keyword {
		r := &maps.NearbySearchRequest{
			Keyword: key,
			// OpenNow: false,
		}

		r.RankBy = maps.RankByDistance
		r.Location = &location

		rsp, err := googleMapsClient.NearbySearch(context.Background(), r)
		if err != nil {
			return nil, err
		}
		information := extractInfo(rsp)
		resp = append(resp, information...)

	}

	return resp, nil
}

func extractInfo(response maps.PlacesSearchResponse) []map[string]interface{} {
	var extracted []map[string]interface{}

	for _, result := range response.Results {
		info := map[string]interface{}{
			"Location": result.Geometry.Location,
			"Name":     result.Name,
			"Rating":   result.Rating,
		}
		extracted = append(extracted, info)
	}
	return extracted
}
