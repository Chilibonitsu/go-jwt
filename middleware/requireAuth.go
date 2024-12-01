package middleware

import (
	"fmt"
	"go-jwt/initializers"
	"go-jwt/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	var flag bool = true
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		fmt.Println("No access token cookie")
		flag = false
		//c.AbortWithStatus(http.StatusUnauthorized)
	}
	refreshTokenString, err := c.Cookie("RefreshToken")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	fmt.Println("auth", tokenString)
	fmt.Println("ref", refreshTokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {

		//c.AbortWithStatus(http.StatusUnauthorized)
		if flag {
			fmt.Println(err)
			fmt.Println("Invalid signing method")
			// c.AbortWithStatus(http.StatusUnauthorized)
			// return
		} else {
			fmt.Println(err)
		}

	}
	if flag { //по сути не нужно
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				fmt.Println("Access token has expired")
			}

		}
	}

	refreshToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRETREF")), nil
	})
	if err != nil {
		fmt.Println("here0")
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	fmt.Println("here")
	if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok {
		fmt.Println("here2")
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
			fmt.Println("here3")
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set("user", user)
		fmt.Println("here4")
		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
