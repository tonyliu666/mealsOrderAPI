package router

import (
	"net/http"
	"weather/handlers"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	diets := router.Group("/diets")
	{
		// get the weather of a city name
		diets.GET("/morning", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Good morning",
			})
		})
		diets.GET("/afternoon", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Good afternoon",
			})
		})
		diets.GET("/evening", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Good evening",
			})
		})
		
		diets.POST("/morning", func(c *gin.Context) {
			handlers.RecordMeal(c)
		})
	}
	orders := router.Group("/orders")
	{
		orders.GET("/:timeslot/:period", func(c *gin.Context) {
			// search the meals that has been eaten in the morning with given period
			diets,err := handlers.GetDiets(c)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, diets)
		})
	}


	return router
}
