package router

import (
	"log"
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
		diets.GET("/:timeslot/:periods", func(c *gin.Context) {
			// search the meals that has been eaten in the morning with given period
			diets, err := handlers.GetDiets(c)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, diets)
		})
		diets.POST("/meals", func(c *gin.Context) {
			handlers.RecordMeal(c)
		})
	}
	orders := router.Group("/orders")
	{
		orders.GET("/healthy/:timeslot/:periods", func(c *gin.Context) {
			// get the recommendation for the given timeslot
			diets, err := handlers.GetDiets(c)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			// meals I want to eat in a heakthy way
			meals, err := handlers.Recommendation(diets)

			c.JSON(http.StatusOK, meals)
		})
		orders.GET("/enjoyable/:timeslot/:periods", func(c *gin.Context) {
			// get the recommendation for the given timeslot
			recommendation, err := handlers.GetDiets(c)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, recommendation)
		})

	}
	promotions := router.Group("/shop")
	{
		promotions.GET("/healthy/:location/:radius/:timeslot/:periods", func(c *gin.Context) {
			diets, err := handlers.GetDiets(c)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			// meals I want to eat in a heakthy way
			meals, err := handlers.Recommendation(diets)
			log.Println(meals)

			shops, err := handlers.GetShops(c, meals)
			c.JSON(http.StatusOK, shops)
		})
	}

	return router
}
