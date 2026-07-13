package controllers

import (
	"backend_golang/models"
	"backend_golang/setup"
	"backend_golang/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input struct {
		Name            string `json:"name" binding:"required"`
		Email           string `json:"email" binding:"required,email"`
		Password        string `json:"password" binding:"required,min=8"`
		ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
		ProfilePicture  string `json:"profile_picture"`
	}

	// Validate the request body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create a new user instance
	user := models.User{
		Name:           input.Name,
		Email:          input.Email,
		Password:       string(hashedPassword),
		ProfilePicture: input.ProfilePicture,
		Level:          "user", // You can modify this based on the user role
	}

	// Save the user in the database
	if err := setup.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func Login(c *gin.Context) {
	var input struct {
		Email      string `json:"email" binding:"required"`
		Password   string `json:"password" binding:"required"`
		RememberMe bool   `json:"remember_me"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := setup.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	var tokenDuration time.Duration
	if input.RememberMe {
		tokenDuration = 7 * 24 * time.Hour
	} else {
		tokenDuration = 24 * time.Hour
	}

	tokenString, err := utils.GenerateJWT(uint(user.Id), tokenDuration)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to generate token", "authenticated": false})
		return
	}

	c.SetCookie("Authorization", tokenString, int(tokenDuration.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"Name":          user.Name,
		"Level":         user.Level,
		"authenticated": true,
	})
}

func GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "status": false})
		return
	}

	var user models.User
	if err := setup.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "status": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"auth": true,
		"data": user,
	})
}

func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", "", false, true)

	// Kirim respon logout sukses
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
