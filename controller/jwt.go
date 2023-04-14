package controller

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/subkhiyoga/auth-jwt/model"
	"github.com/subkhiyoga/auth-jwt/usecase"
)

type LoginJwt struct {
	usecase usecase.LoginUsecase
	jwtKey  []byte
}

var jwtKey = []byte("test")

func generateToken(cre *model.Credentials) (string, error) {
	// set token claims
	claims := jwt.MapClaims{}
	claims["username"] = cre.Username
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // mengatur masa waktu token

	// create token with claims and secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func AuthMiddleware(jwtKey []byte) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")

		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		username := claims["username"].(string)

		ctx.Set("username", username)

		ctx.Next()
	}
}

func (l *LoginJwt) Login(ctx *gin.Context) {
	var c model.Credentials

	err := ctx.ShouldBindJSON(&c)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	credentials, err := l.usecase.Login(c.Username, c.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	if credentials == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := generateToken(credentials)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (l *LoginJwt) Profile(ctx *gin.Context) {
	// ambil username dari JWT token
	claims := ctx.MustGet("claims").(jwt.MapClaims)
	username := claims["username"].(string)

	// dapatkan informasi user dari database (dalam hal ini, return username)
	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Welcome to profile",
		"username": username,
	})
}

func NewCredentialsJwt(u usecase.LoginUsecase) *LoginJwt {
	loginjwt := LoginJwt{
		usecase: u,
		jwtKey:  jwtKey,
	}

	return &loginjwt
}
