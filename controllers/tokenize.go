package controllers

import (
	"go-jwt/initializers"
	"go-jwt/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Tokenize(c *gin.Context) {

	var body struct {
		Guid     string
		Username string
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}
	var user models.User

	initializers.DB.First(&user, "guid = ?", body.Guid)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": user.Guid,
		"exp": time.Now().Add(time.Second * 60).Unix(),
		"ip":  c.ClientIP(),
	})

	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create access token",
		})
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": user.Guid,
		"exp": time.Now().Add(time.Second * 60).Unix(),
		"ip":  c.ClientIP(),
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("SECRETREF")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create refresh token",
		})
	}
	initializers.DB.Model(&user).
		Where("guid = ?", body.Guid).
		Update("refresh_token", refreshTokenString)

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", accessTokenString, 120, "", "", true, true)
	c.SetCookie("RefreshToken", refreshTokenString, 360, "", "", true, true)

	c.JSON(http.StatusOK, gin.H{})

}

func RefreshTokens(c *gin.Context) {

	body := c.MustGet("body").(models.Body)

	var user models.User

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": body.Guid,
		"exp": time.Now().Add(time.Second * 60).Unix(),
		"ip":  c.ClientIP(),
	})

	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create access token",
		})
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": body.Guid,
		"exp": time.Now().Add(time.Second * 60).Unix(),
		"ip":  c.ClientIP(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("SECRETREF")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create refresh token",
		})
	}

	initializers.DB.Model(&user).
		Where("guid = ?", body.Guid).
		Update("refresh_token", refreshTokenString)

	c.SetCookie("Authorization", accessTokenString, 120, "/", "", true, true)
	c.SetCookie("RefreshToken", refreshTokenString, 3600*24*30, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{})

}

func SignUp(c *gin.Context) {

	var body struct {
		Guid     string
		Username string
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}
	user := models.User{Guid: body.Guid, Username: body.Username, Email: body.Email, PasswordHash: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
