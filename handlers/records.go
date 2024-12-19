package handlers

import (
	"net/http"
	"weather/models"

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
func assignValue(meal request, timeslots string) models.Diets {
	var diet models.Diets
	diet.Name = meal.Name
	diet.Location = meal.Location
	diet.Date = meal.Date
	diet.Time = meal.Time
	diet.Periods = meal.Periods
	diet.TimeSlots = timeslots
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
	// assign the values to the diet struct
	diet := assignValue(meal, timeslots)

	if err := diet.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "meal recorded successfully"})
}
func Recommendation(c *gin.Context) ([]models.Diets, error) {
	// timeslots: morning, afternoon, evening
	// periods: how long the meal will last
	timeslots := c.Param("timeslot")
	periods := c.Param("period")
	diets, err := models.QueryDates(timeslots, periods)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, err
	}
	return diets, nil

}
