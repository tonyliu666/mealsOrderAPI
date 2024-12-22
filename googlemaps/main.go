package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kr/pretty"
	"googlemaps.github.io/maps"
)

var (
	apiKey    = flag.String("key", "", "API Key for using Google Maps API.")
	clientID  = flag.String("client_id", "", "ClientID for Maps for Work API access.")
	signature = flag.String("signature", "", "Signature for Maps for Work API access.")
	location  = flag.String("location", "", "The latitude/longitude around which to retrieve place information. This must be specified as latitude,longitude.")
	radius    = flag.Uint("radius", 0, "Defines the distance (in meters) within which to bias place results. The maximum allowed radius is 50,000 meters.")
	keyword   = flag.String("keyword", "", "A term to be matched against all content that Google has indexed for this place, including but not limited to name, type, and address, as well as customer reviews and other third-party content.")
	language  = flag.String("language", "", "The language in which to return results.")
	minPrice  = flag.String("minprice", "", "Restricts results to only those places within the specified price level.")
	maxPrice  = flag.String("maxprice", "", "Restricts results to only those places within the specified price level.")
	// name      = flag.String("name", "", "One or more terms to be matched against the names of places, separated with a space character.")
	openNow   = flag.Bool("open_now", false, "Restricts results to only those places that are open for business at the time the query is sent.")
	rankBy    = flag.String("rankby", "", "Specifies the order in which results are listed. Valid values are prominence or distance.")
	placeType = flag.String("type", "", "Restricts the results to places matching the specified type.")
	pageToken = flag.String("pagetoken", "", "Set to retrieve the next page of results.")
)

func usageAndExit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	fmt.Println("Flags:")
	flag.PrintDefaults()
	os.Exit(2)
}

func check(err error) {
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
}

func main() {
	flag.Parse()

	var client *maps.Client
	var err error
	if *apiKey != "" {
		client, err = maps.NewClient(maps.WithAPIKey(*apiKey))
	} else if *clientID != "" || *signature != "" {
		client, err = maps.NewClient(maps.WithClientIDAndSignature(*clientID, *signature))
	} else {
		usageAndExit("Please specify an API Key, or Client ID and Signature.")
	}
	check(err)

	r := &maps.NearbySearchRequest{
		Radius:    *radius,
		Keyword:   *keyword,
		Language:  *language,
		// Name:      *name,
		OpenNow:   *openNow,
		PageToken: *pageToken,
	}

	parseLocation(*location, r)
	parsePriceLevels(*minPrice, *maxPrice, r)
	parseRankBy(*rankBy, r)
	parsePlaceType(*placeType, r)

	resp, err := client.NearbySearch(context.Background(), r)
	check(err)

	pretty.Println(resp)
}

func parseLocation(location string, r *maps.NearbySearchRequest) {
	if location != "" {
		l, err := maps.ParseLatLng(location)
		check(err)
		r.Location = &l
	}
}

func parsePriceLevel(priceLevel string) maps.PriceLevel {
	switch priceLevel {
	case "0":
		return maps.PriceLevelFree
	case "1":
		return maps.PriceLevelInexpensive
	case "2":
		return maps.PriceLevelModerate
	case "3":
		return maps.PriceLevelExpensive
	case "4":
		return maps.PriceLevelVeryExpensive
	default:
		usageAndExit(fmt.Sprintf("Unknown price level: '%s'", priceLevel))
	}
	return maps.PriceLevelFree
}

func parsePriceLevels(minPrice string, maxPrice string, r *maps.NearbySearchRequest) {
	if minPrice != "" {
		r.MinPrice = parsePriceLevel(minPrice)
	}

	if maxPrice != "" {
		r.MaxPrice = parsePriceLevel(minPrice)
	}
}

func parseRankBy(rankBy string, r *maps.NearbySearchRequest) {
	switch rankBy {
	case "prominence":
		r.RankBy = maps.RankByProminence
		return
	case "distance":
		r.RankBy = maps.RankByDistance
		return
	case "":
		return
	default:
		usageAndExit(fmt.Sprintf("Unknown rank by: \"%v\"", rankBy))
	}
}

func parsePlaceType(placeType string, r *maps.NearbySearchRequest) {
	if placeType != "" {
		t, err := maps.ParsePlaceType(placeType)
		if err != nil {
			usageAndExit(fmt.Sprintf("Unknown place type \"%v\"", placeType))
		}

		r.Type = t
	}
}

// package main

// import (
// 	"context"
// 	"flag"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strconv"
// 	"strings"

// 	"github.com/kr/pretty"
// 	"googlemaps.github.io/maps"
// )

// var (
// 	apiKey       = flag.String("key", "", "API Key for using Google Maps API.")
// 	clientID     = flag.String("client_id", "", "ClientID for Maps for Work API access.")
// 	signature    = flag.String("signature", "", "Signature for Maps for Work API access.")
// 	address      = flag.String("address", "", "The street address that you want to geocode, in the format used by the national postal service of the country concerned.")
// 	components   = flag.String("components", "", "A component filter for which you wish to obtain a geocode.")
// 	bounds       = flag.String("bounds", "", "The bounding box of the viewport within which to bias geocode results more prominently.")
// 	language     = flag.String("language", "", "The language in which to return results.")
// 	region       = flag.String("region", "", "The region code, specified as a ccTLD two-character value.")
// 	latlng       = flag.String("latlng", "", "The textual latitude/longitude value for which you wish to obtain the closest, human-readable address.")
// 	resultType   = flag.String("result_type", "", "One or more address types, separated by a pipe (|).")
// 	locationType = flag.String("location_type", "", "One or more location types, separated by a pipe (|).")
// 	//enableAddressDescriptor = flag.String("enable_address_descriptor", "", "True or False.  Whether to return the Address Descriptors in the response.")
// )

// func usageAndExit(msg string) {
// 	fmt.Fprintln(os.Stderr, msg)
// 	fmt.Println("Flags:")
// 	flag.PrintDefaults()
// 	os.Exit(2)
// }

// func check(err error) {
// 	if err != nil {
// 		log.Fatalf("fatal error: %s", err)
// 	}
// }

// func main() {
// 	flag.Parse()

// 	var client *maps.Client
// 	var err error
// 	if *apiKey != "" {
// 		client, err = maps.NewClient(maps.WithAPIKey(*apiKey))
// 	} else if *clientID != "" || *signature != "" {
// 		client, err = maps.NewClient(maps.WithClientIDAndSignature(*clientID, *signature))
// 	} else {
// 		usageAndExit("Please specify an API Key, or Client ID and Signature.")
// 	}
// 	check(err)

// 	r := &maps.GeocodingRequest{
// 		Address:  *address,
// 		Language: *language,
// 		Region:   *region,
// 	}

// 	parseComponents(*components, r)
// 	parseBounds(*bounds, r)
// 	parseLatLng(*latlng, r)
// 	parseResultType(*resultType, r)
// 	parseLocationType(*locationType, r)
	

// 	resp, err := client.Geocode(context.Background(), r)
// 	check(err)

// 	pretty.Println(resp)
// }

// func parseComponents(components string, r *maps.GeocodingRequest) {
// 	if components == "" {
// 		return
// 	}
// 	if r.Components == nil {
// 		r.Components = make(map[maps.Component]string)
// 	}

// 	c := strings.Split(components, "|")
// 	for _, cf := range c {
// 		i := strings.Split(cf, ":")
// 		switch i[0] {
// 		case "route":
// 			r.Components[maps.ComponentRoute] = i[1]
// 		case "locality":
// 			r.Components[maps.ComponentLocality] = i[1]
// 		case "administrative_area":
// 			r.Components[maps.ComponentAdministrativeArea] = i[1]
// 		case "postal_code":
// 			r.Components[maps.ComponentPostalCode] = i[1]
// 		case "country":
// 			r.Components[maps.ComponentCountry] = i[1]
// 		default:
// 			log.Fatalf("parseComponents: component name %#v unknown", i[0])
// 		}
// 	}
// }

// func parseBounds(bounds string, r *maps.GeocodingRequest) {
// 	if bounds != "" {
// 		b := strings.Split(bounds, "|")
// 		sw := strings.Split(b[0], ",")
// 		ne := strings.Split(b[1], ",")

// 		swLat, err := strconv.ParseFloat(sw[0], 64)
// 		if err != nil {
// 			log.Fatalf("Couldn't parse bounds: %#v", err)
// 		}
// 		swLng, err := strconv.ParseFloat(sw[1], 64)
// 		if err != nil {
// 			log.Fatalf("Couldn't parse bounds: %#v", err)
// 		}
// 		neLat, err := strconv.ParseFloat(ne[0], 64)
// 		if err != nil {
// 			log.Fatalf("Couldn't parse bounds: %#v", err)
// 		}
// 		neLng, err := strconv.ParseFloat(ne[1], 64)
// 		if err != nil {
// 			log.Fatalf("Couldn't parse bounds: %#v", err)
// 		}

// 		r.Bounds = &maps.LatLngBounds{
// 			NorthEast: maps.LatLng{Lat: neLat, Lng: neLng},
// 			SouthWest: maps.LatLng{Lat: swLat, Lng: swLng},
// 		}
// 	}
// }

// func parseLatLng(latlng string, r *maps.GeocodingRequest) {
// 	if latlng != "" {
// 		l := strings.Split(latlng, ",")
// 		lat, err := strconv.ParseFloat(l[0], 64)
// 		if err != nil {
// 			log.Fatalf("Couldn't parse latlng: %#v", err)
// 		}
// 		lng, err := strconv.ParseFloat(l[1], 64)
// 		if err != nil {
// 			log.Fatalf("Couldn't parse latlng: %#v", err)
// 		}
// 		r.LatLng = &maps.LatLng{
// 			Lat: lat,
// 			Lng: lng,
// 		}
// 	}
// }

// func parseResultType(resultType string, r *maps.GeocodingRequest) {
// 	if resultType != "" {
// 		r.ResultType = strings.Split(resultType, "|")
// 	}
// }

// func parseLocationType(locationType string, r *maps.GeocodingRequest) {
// 	if locationType != "" {
// 		for _, l := range strings.Split(locationType, "|") {
// 			switch l {
// 			case "ROOFTOP":
// 				r.LocationType = append(r.LocationType, maps.GeocodeAccuracyRooftop)
// 			case "RANGE_INTERPOLATED":
// 				r.LocationType = append(r.LocationType, maps.GeocodeAccuracyRangeInterpolated)
// 			case "GEOMETRIC_CENTER":
// 				r.LocationType = append(r.LocationType, maps.GeocodeAccuracyGeometricCenter)
// 			case "APPROXIMATE":
// 				r.LocationType = append(r.LocationType, maps.GeocodeAccuracyApproximate)
// 			}
// 		}

// 	}
// }

