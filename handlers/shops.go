package handlers

import (
	"log"
	"strconv"
	"weather/models/gateway"

	gin "github.com/gin-gonic/gin"
	"github.com/kr/pretty"
	"googlemaps.github.io/maps"
)

func shopsRecommendation(c *gin.Context, loction []maps.LatLng, meals []string) (string, error) {
	// get the shops for the given meals
	radiusStr := c.Param("radius")
	radius, err := strconv.ParseUint(radiusStr, 10, 32)
	if err != nil {
		return "", err
	}
	radiusParam := uint(radius)
	latlng := loction[0]

	result, err := gateway.NearbySearch(radiusParam, latlng, meals)
	log.Println("result", result)
	if err != nil {
		return "", err
	}
	return pretty.Sprintf("%# v", result), nil
}

func GetShops(c *gin.Context, meals []string) (string, error) {
	// get the shops for the given meals
	location := c.Param("location")
	longlatude, err := gateway.GetCurrentLocation(location)
	if err != nil {
		return "", err
	}
	shops, err := shopsRecommendation(c, longlatude, meals)
	if err != nil {
		return "", err
	}
	return shops, nil
}
