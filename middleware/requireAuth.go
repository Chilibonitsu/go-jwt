package middleware

import (
	"errors"
	"fmt"
	"go-jwt/initializers"
	"go-jwt/models"
	"net/http"
	"os"
	"time"

	"go-jwt/testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func BindBody(c *gin.Context) (bodyStr models.Body, err error) {

	var body models.Body

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to redfsgdfgdfggad bodyss",
		})

		return body, fmt.Errorf("failed to read body")

	}
	c.Set("body", body)
	return body, nil

}
func RequireAuth(c *gin.Context) {

	var accessCookieExp bool = false

	body, err := BindBody(c)
	if err != nil {
		fmt.Println(err)
	}

	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		fmt.Println("No access token cookie")
		//flag = false
		//c.AbortWithStatus(http.StatusUnauthorized)
	}
	refreshTokenString, err := c.Cookie("RefreshToken")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	fmt.Println("Access token: ", tokenString)
	fmt.Println("Refresh token: ", refreshTokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil

	})

	if err != nil {

		fmt.Println(err)

		// c.AbortWithStatus(http.StatusUnauthorized)
		// return

		accessCookieExp = true

	}
	if !accessCookieExp {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				fmt.Println("Access token has expired")
			}
			err, isChanged := CheckIP(c, claims)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			if isChanged {
				//fmt.Println(err)
				fmt.Println("Sending email")
				testing.SendEmail()
			}

			err = CheckUser(c, claims, body)

			if err != nil {

				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

		}
	}

	refreshToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRETREF")), nil
	})

	if err != nil {

		c.AbortWithStatus(http.StatusUnauthorized)
	}

	if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok {

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			fmt.Println("Refresh token has expired")
		}
		if claims["ip"].(string) != c.ClientIP() {
			fmt.Println("Ip's are different")
		}

		var user models.User
		fmt.Println("GUID From claims", claims["sub"].(string))
		initializers.DB.First(&user, "guid = ?", claims["sub"])

		if user.ID == 0 {

			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set("user", user)

		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
func CheckIP(c *gin.Context, claims jwt.MapClaims) (err error, isChanged bool) {

	fmt.Println(c.ClientIP() + " " + claims["ip"].(string))
	if claims["ip"].(string) != c.ClientIP() {
		fmt.Printf("IP mismatch: token IP %s, request IP %s", claims["ip"].(string), c.ClientIP())

		return nil, true
	}
	return nil, false

}
func CheckUser(c *gin.Context, claim jwt.MapClaims, body models.Body) (err error) {

	var user models.User
	if result := initializers.DB.Where("guid = ?", body.Guid).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {

			return result.Error
		}
	}
	return nil
}
