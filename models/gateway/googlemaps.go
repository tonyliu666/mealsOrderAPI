package gateway

import (
	"context"
	"os"
	"sync"

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
	//call the google maps api
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

	// create the number of goroutines equal to the number of keywords
	// create wait group
	var wg sync.WaitGroup
	wg.Add(len(keyword))
	errChan := make(chan error, len(keyword))

	for _, key := range keyword {
		go func() {
			r := &maps.NearbySearchRequest{
				Keyword: key,
				// OpenNow: false,
			}

			r.RankBy = maps.RankByDistance
			r.Location = &location

			rsp, err := googleMapsClient.NearbySearch(context.Background(), r)
			if err != nil {
				errChan <- err
				return
			}

			information := extractInfo(rsp)
			resp = append(resp, information...)
			wg.Done()
		}()
	}
	wg.Wait()
	// check if there is any error
	select {
	case err := <-errChan:
		return nil, err
	default:
		return resp, nil
	}
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
