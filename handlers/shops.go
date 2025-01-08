package handlers

import (
	"weather/models/gateway"

	"googlemaps.github.io/maps"
)

type Shop struct {
	Name     string      `json:"name"`
	Location maps.LatLng `json:"location"`
	Rating   float32     `json:"rating"`
}

func shopsRecommendation(loction []maps.LatLng, meals []string) ([]Shop, error) {
	// get the shops for the given meals
	var shops []Shop
	latlng := loction[0]
	result, err := gateway.NearbySearch(latlng, meals)

	if err != nil {
		return nil, err
	}
	// each element in the result is a map
	// I want to return a json string with the key same as the struct Name field and rest of the fields as the value

	for _, r := range result {
		shop := Shop{
			Name:     r["Name"].(string),
			Location: r["Location"].(maps.LatLng),
			Rating:   r["Rating"].(float32),
		}
		shops = append(shops, shop)
	}
	return shops, nil
}

func GetShops(meals []string, longlatude []maps.LatLng) ([]Shop, error) {
	// get the shops for the given meals
	shops, err := shopsRecommendation(longlatude, meals)
	if err != nil {
		return nil, err
	}
	return shops, nil
}
