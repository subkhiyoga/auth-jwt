package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("hahaha")

type Mahasiswa struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	// buat gin router
	router := gin.Default()

	// set up routes for login
	router.POST("/auth/login", login)

	// other routes
	mahasiwaRouter := router.Group("/api/v1/mahasiswa")
	mahasiwaRouter.Use(authMiddleware())

	mahasiwaRouter.GET("/:id/profile", profile)

	// start server
	log.Fatal(router.Run(":8080"))
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("claims", claims)

		c.Next()
	}
}

func login(c *gin.Context) {
	var mahasiwa Mahasiswa

	err := c.ShouldBindJSON(&mahasiwa)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return //agar err berhenti
	}

	// authenticate mahasiswa (compare username dan password)
	if mahasiwa.Username == "hahaha" && mahasiwa.Password == "a12345." {
		// generate JWT token
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = mahasiwa.Username
		claims["exp"] = time.Now().Add(time.Minute * 3).Unix()

		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unregistered mahasiswa"})
	}
}

func profile(c *gin.Context) {
	// ambil username dari JWT token
	claims := c.MustGet("claims").(jwt.MapClaims)
	username := claims["username"].(string)

	// dapatkan informasi user dari database (dalam hal ini, return username)
	c.JSON(http.StatusOK, gin.H{
		"message":  "Welcome to profile",
		"username": username,
	})
}
