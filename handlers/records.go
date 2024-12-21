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
func assignValue(meal request, timeslots string, nutritions gateway.Nutrition) database.Diets {
	var diet database.Diets
	ingredients := &database.Ingredients{
		Carolie:      nutritions.Carolie,
		Protein:      nutritions.Protein,
		Fat:          nutritions.Fat,
		Carbohydrate: nutritions.Sugar,
	}
	diet = database.Diets{
		Name:        meal.Name,
		Location:    meal.Location,
		Date:        meal.Date,
		Time:        meal.Time,
		Periods:     meal.Periods,
		TimeSlots:   timeslots,
		Ingredients: ingredients,
	}
	return diet
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
	nutritions := gateway.GetNutritionAnalysis(meal.Name)
	diet := assignValue(meal, timeslots, nutritions)
	if err := diet.Ingredients.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := diet.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "meal recorded successfully"})
}
func GetDiets(c *gin.Context) ([]database.Diets, error) {
	// timeslots: morning, afternoon, evening
	// periods: how long the meal will last
	timeslots := c.Param("timeslot")
	periods := c.Param("period")
	diets, err := database.QueryDates(timeslots, periods)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, err
	}
	return diets, nil

}
