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
	Carolie float64 `json:"carolie"`
	Sugar   float64 `json:"sugar"`
	Protein float64 `json:"protein"`
	Fat     float64 `json:"fat"`
}

func getNutrieAmount(totalNutrients map[string]interface{}, kind string) (float64, error) {
	// get the sugar nurtrient
	kinds, ok := totalNutrients[kind].(map[string]interface{})
	if !ok {
		return 0.0, fmt.Errorf("please enter a valid unit of your food or give concrete information")
	}
	quantity := kinds["quantity"].(float64)
	unit := kinds["unit"].(string)
	if unit == "mg" {
		quantity = quantity / 1000
	} else if unit == "kg" {
		quantity = quantity * 1000
	}
	return quantity, nil
}
func getCarolieAmount(totalNutrients map[string]interface{}) (float64, error) {
	// get the sugar nurtrient
	energy, ok := totalNutrients["ENERC_KCAL"].(map[string]interface{})
	if !ok {
		return 0.0, fmt.Errorf("cannot count the total calories")
	}
	quantity := energy["quantity"].(float64)
	return quantity * 1000, nil
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

func GetNutritionAnalysis(ingredient string) (Nutrition, error) {
	baseURL := "https://api.edamam.com/api/nutrition-data"
	appID := os.Getenv("APPID")
	appKey := os.Getenv("APPKeys")
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
		return Nutrition{}, fmt.Errorf("failed to send GET request: %v", err)
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Nutrition{}, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the response body to json format
	totalNutrients := unmarsalUtility(body)

	// Parse the JSON string into the map
	sugarAmount, err := getNutrieAmount(totalNutrients, "SUGAR")
	if err != nil {
		return Nutrition{}, err
	}
	proteinAmount, err := getNutrieAmount(totalNutrients, "PROCNT")
	if err != nil {
		return Nutrition{}, err
	}
	fatAmount, err := getNutrieAmount(totalNutrients, "FAT")
	if err != nil {
		return Nutrition{}, err
	}
	carolie, err := getCarolieAmount(totalNutrients)
	if err != nil {
		return Nutrition{}, err
	}

	nutrition := Nutrition{
		Carolie: carolie,
		Sugar:   sugarAmount,
		Protein: proteinAmount,
		Fat:     fatAmount,
	}

	return nutrition, nil
}
