package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func main() {
	baseURL := "https://api.edamam.com/api/nutrition-data"
	appID := "f6558e99"
	appKey := "78da3304835443dd714322492c7f62fe"
	nutritionType := "cooking"
	ingredient := "5 orz beef"

	// Build the query parameters
	params := url.Values{}
	params.Add("app_id", appID)
	params.Add("app_key", appKey)
	params.Add("nutrition-type", nutritionType)
	params.Add("ingr", ingredient)
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Send the GET request
	response, err := http.Get(fullURL)
	if err != nil {
		log.Fatalf("Failed to send GET request: %v", err)
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	var result map[string]interface{}

	// Parse the JSON string into the map
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	
	totalNutrients, ok := result["totalNutrients"].(map[string]interface{})
	if !ok {
		log.Fatal("Error: totalNutrients is not a map")
	}

	// Access the SUGAR field
	sugar, ok := totalNutrients["SUGAR"].(map[string]interface{})
	if !ok {
		log.Fatal("Error: SUGAR is not a map")
	}
	label := sugar["label"].(string)
	quantity := sugar["quantity"].(float64)
	unit := sugar["unit"].(string)

	// Print the values
	fmt.Println("Label:", label)
	fmt.Println("Quantity:", quantity)
	fmt.Println("Unit:", unit)
}
