package handlers

import (
	"math/rand/v2"
	"net/http"
	"weather/models/database"
	"weather/utils"

	"github.com/gin-gonic/gin"
)

// Function for logging in
func Login(c *gin.Context) {
	var user database.Client

	// Check user credentials and generate a JWT token
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Check if credentials are valid (replace this logic with real authentication)
	if user.Username == "user" && user.Password == "password" {
		// Generate a JWT token
		token, err := utils.GenerateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

// Function for registering a new user (for demonstration purposes)
func Register(c *gin.Context) {
	// realuser := database.NewDBManager("users")
	var user database.Client

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	user.ID = uint(rand.Uint32())
	// save the user to the database (replace this with real database logic)
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
	}
	user.Password = hashedPassword
	err = user.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
