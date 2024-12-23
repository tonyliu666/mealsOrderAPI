package handlers

import (
	"weather/models/database"
	"weather/models/gateway"
)

func nutritionCalculation(diet []database.Diets) []float64 {
	// calculate the nutrition for the given diet
	length := float64(len(diet))
	suagrIngestions := 0.0
	proteinIngestions := 0.0
	fatIngestions := 0.0
	carolieIngestions := 0.0
	for _, d := range diet {
		suagrIngestions += d.Ingredients.Carbohydrate
		proteinIngestions += d.Ingredients.Protein
		fatIngestions += d.Ingredients.Fat
		carolieIngestions += d.Ingredients.Carolie
	}
	averageSugar := suagrIngestions / length
	averageProtein := proteinIngestions / length
	averageFat := fatIngestions / length
	averageCarolie := carolieIngestions / length
	return []float64{averageSugar, averageProtein, averageFat, averageCarolie}

}

func Recommendation(diet []database.Diets) ([]string, error) {
	// get the recommendation for the given diet
	nutrition := nutritionCalculation(diet)
	//calculate the average ration amoung the average sugar, protein, fat
	if nutrition[1] == 0 {
		// recommend protein food
		meals,err := gateway.RecommendNutrition("protein")
		if err != nil {
			return nil, err
		}
		return meals, nil
	}else{
		ratio := nutrition[0] / nutrition[1]
		if ratio > 2 {
			meals, err := gateway.RecommendNutrition("protein")
			if err != nil {
				return nil, err
			}
			return meals, nil
		} else{
			if nutrition[2]==0{
				// recommend fat food
				meals,err := gateway.RecommendNutrition("fat")
				if err != nil {
					return nil, err
				}
				return meals, nil
			}else{
				if nutrition[2]/nutrition[1] > 2 {
					// recommend fat food
					meals,err := gateway.RecommendNutrition("Carbonhydrate")
					if err != nil {
						return nil, err
					}
					return meals, nil
				}else{
					// recommend protein food
					meals,err := gateway.RecommendNutrition("fat")
					if err != nil {
						return nil, err
					}
					return meals, nil
				}
			}
		}
	}
}
