package handlers

import "weather/models/database"

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

func Recommendation(diet []database.Diets) error {
	// get the recommendation for the given diet
	nutrition := nutritionCalculation(diet)
	//calculate the average ration amoung the average sugar, protein, fat
	if nutrition[1] == 0 {
		// recommend protein food
	}else{
		ratio := nutrition[0] / nutrition[1]
		if ratio > 2 {
			// recommend protein food
		} else {
			// recommend sugar food
		}
	}
	// recommend food contains the fat 

}
