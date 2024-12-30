package handlers

import (
	"log"
	"weather/models/gateway"

	gin "github.com/gin-gonic/gin"
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
	log.Println(result)
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

func GetShops(c *gin.Context, meals []string) ([]Shop, error) {
	// get the shops for the given meals
	location := c.Param("location")
	longlatude, err := gateway.GetCurrentLocation(location)
	if err != nil {
		return nil, err
	}
	shops, err := shopsRecommendation(longlatude, meals)
	if err != nil {
		return nil, err
	}
	return shops, nil
}
