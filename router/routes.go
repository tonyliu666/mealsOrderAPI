package router

import (
	"net/http"
	"sync"
	"time"
	"weather/handlers"
	"weather/middleware"
	"weather/models/cache"
	"weather/models/gateway"

	"github.com/gin-gonic/gin"
	"googlemaps.github.io/maps"
)

func Init() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	diets := router.Group("/diets")
	{
		diets.GET("/:timeslot/:periods", func(c *gin.Context) {
			// search the meals that has been eaten in the morning with given period
			timeslots := c.Param("timeslot")
			periods := c.Param("periods")
			meals, err := handlers.GetDiets(timeslots, periods)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, meals)
		})

	}

	orders := router.Group("/orders")
	{
		orders.GET("/healthy/:timeslot/:periods", func(c *gin.Context) {
			timeslots := c.Param("timeslot")
			periods := c.Param("periods")
			var meals []string

			meals, err := cache.GetUtility(timeslots, periods)
			if err == nil {
				c.JSON(http.StatusOK, meals)
				return
			}
			diets, err := handlers.GetDiets(timeslots, periods)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			meals, err = handlers.Recommendation(diets)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			err = cache.StoreUtility(meals, timeslots, periods, 20*time.Second)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, meals)
		})
		// TODO: not yet started
		orders.GET("/enjoyable/:timeslot/:periods", func(c *gin.Context) {
			// get the recommendation for the given timeslot
			timeslots := c.Param("timeslot")
			periods := c.Param("periods")
			recommendation, err := handlers.GetDiets(timeslots, periods)
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
		promotions.GET("/healthy/:location/:timeslot/:periods", func(c *gin.Context) {
			var meals []string
			var shops []handlers.Shop
			var longlatude []maps.LatLng

			location := c.Param("location")
			timeslots := c.Param("timeslot")
			periods := c.Param("periods")
			diets, err := handlers.GetDiets(timeslots, periods)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			wg := sync.WaitGroup{}
			wg.Add(2)
			channel := make(chan error, 2)
			go func() {
				meals, err = handlers.Recommendation(diets)
				if err != nil {
					channel <- err
				}
				wg.Done()

			}()
			go func() {
				longlatude, err = gateway.GetCurrentLocation(location)
				if err != nil {
					channel <- err
				}
				wg.Done()
			}()
			wg.Wait()
			// get the error message from the channel
			if len(channel) > 0 {
				err := <-channel
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			shops, err = handlers.GetShops(meals, longlatude)
			if err != nil {
				channel <- err
				return
			}
			c.JSON(http.StatusOK, shops)
		})
	}

	publicRoutes := router.Group("/public")
	{
		publicRoutes.POST("/login", handlers.Login)
		publicRoutes.POST("/register", handlers.Register)
	}

	// // Protected routes (require authentication)
	protectedRoutes := router.Group("/protected")
	protectedRoutes.Use(middleware.AuthenticationMiddleware())
	{
		protectedRoutes.POST("/meals", func(c *gin.Context) {
			handlers.RecordMeal(c)
		})

	}

	return router
}
