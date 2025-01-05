package handlers

import (
	"net/http"
	"weather/models/database"
	"weather/models/gateway"

	gin "github.com/gin-gonic/gin"
)

type request struct {
	Name     string `json:"food_name"`
	Location string `json:"where_eaten"`
	// the time the meal was eaten
	Date    string `json:"date_eaten"`
	Time    string `json:"time_eaten"`
	Periods string `json:"periods"`
}

func validRequest(r request) bool {
	if r.Name == "" || r.Location == "" || r.Time == "" {
		return false
	}
	return true
}
func findTimeSlots(t string) string {
	// t will be in the format "HH:MM:SS"
	if t < "12:00:00" {
		return "morning"
	}
	if t < "18:00:00" {
		return "afternoon"
	}
	return "evening"
}
func assignValue(meal request, timeslots string, nutritions gateway.Nutrition) (database.DBManager, error) {
	meals := database.NewDBManager("ingredients")
	ingredient, _ := meals.(*database.Ingredients)
	ingredient.Carolie = nutritions.Carolie
	ingredient.Protein = nutritions.Protein
	ingredient.Fat = nutritions.Fat
	ingredient.Carbohydrate = nutritions.Sugar

	diets := database.NewDBManager("diets")
	diet, _ := diets.(*database.Diets)
	diet.Name = meal.Name
	diet.Location = meal.Location
	diet.Date = meal.Date
	diet.Time = meal.Time
	diet.Periods = meal.Periods
	diet.TimeSlots = timeslots
	diet.Ingredients = ingredient

	return diets, nil
}

func RecordMeal(c *gin.Context) {
	var meal request
	if err := c.ShouldBindJSON(&meal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !validRequest(meal) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required fields"})
		return
	}
	timeslots := findTimeSlots(meal.Time)
	// get the nutrition analysis
	nutritions, err := gateway.GetNutritionAnalysis(meal.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	diet, err := assignValue(meal, timeslots, nutritions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = diet.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "meal recorded successfully"})
}
func GetDiets(c *gin.Context) ([]database.Diets, error) {
	// timeslots: morning, afternoon, evening
	// periods: how long the meal will last
	var diets []database.Diets
	timeslots := c.Param("timeslot")
	periods := c.Param("periods")
	diets, err := database.QueryDates(timeslots, periods)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, err
	}
	return diets, nil

}
