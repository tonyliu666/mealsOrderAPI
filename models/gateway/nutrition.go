package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)
type Nutrition struct {
	Sugar float64 `json:"sugar"`
	Protein float64 `json:"protein"`
	Fat float64 `json:"fat"`
}
func getSugarAmount(totalNutrients map[string]interface{},kind string) float64 {
	// get the sugar nurtrient
	kinds, ok := totalNutrients[kind].(map[string]interface{})
	if !ok {
		log.Fatal("Error: SUGAR is not a map")
	}
	quantity := kinds["quantity"].(float64)
	unit := kinds["unit"].(string)
	if unit == "mg"{
		quantity = quantity / 1000
	}else if unit == "kg"{
		quantity = quantity * 1000
	}
	return quantity
}
func unmarsalUtility(body []byte) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	// totalNutrients 
	totalNutrients, ok := result["totalNutrients"].(map[string]interface{})
	if !ok {
		log.Fatal("Error: totalNutrients is not a map")
	}
	return totalNutrients
}

func GetNutritionAnalysis(ingredient string) Nutrition{
	baseURL := "https://api.edamam.com/api/nutrition-data"
	appID := os.Getenv("APPID")
	appKey := os.Getenv("APPKEY")
	nutritionType := "cooking"

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

	// Parse the response body to json format
	
	totalNutrients := unmarsalUtility(body)

	// Parse the JSON string into the map
	sugarAmount := getSugarAmount(totalNutrients,"SUGAR")
	proteinAmount := getSugarAmount(totalNutrients,"PROCNT")
	fatAmount := getSugarAmount(totalNutrients,"FAT")

	nutrition := Nutrition{
		Sugar: sugarAmount,
		Protein: proteinAmount,
		Fat: fatAmount,
	}
	return nutrition
}
