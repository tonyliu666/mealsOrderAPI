package handlers

import (
	"log"
	"math/rand/v2"
	"net/http"
	"time"
	"weather/models/cache"
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
	// deep copy the user
	realuser := user
	// check the request with session id exists in redis cache
	sessionID := c.GetHeader("Authorization")
	_, err := cache.Get(sessionID)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "User already logged in"})
		return
	}

	// check the user name exists in redis cache
	err = realuser.Read()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if !utils.VerifyPassword(user.Password, realuser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	// return a session id(unsigned int) to the user and generate jwt token

	jwt, err := utils.GenerateToken(realuser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}
	log.Println("jwt", jwt)
	// save the session id to the redis cache
	err = cache.Save(jwt, realuser.Username, 10*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "token": jwt})

}

// Function for registering a new user (for demonstration purposes)
func Register(c *gin.Context) {
	// realuser := database.NewDBManager("users")
	var user database.Client

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	// check the user name exists in database
	dbuser := database.Client{Username: user.Username}
	err := dbuser.Read()
	if err == nil {
		// return 409 if the user name exists
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
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
