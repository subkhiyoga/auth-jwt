package controll

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/subkhiyoga/auth-jwt/model"
	"github.com/subkhiyoga/auth-jwt/service"
)

type LoginJwt struct {
	service service.LoginService
	jwtKey  []byte
}

var jwtKey = []byte("test")

func generateToken(cre *model.Credentials) (string, error) {
	// set token claims
	claims := jwt.MapClaims{}
	claims["username"] = cre.Username
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix() // mengatur masa waktu token

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

	credentials, err := l.service.Login(c.Username, c.Password)
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

func NewCredentialsJwt(s service.LoginService) *LoginJwt {
	loginjwt := LoginJwt{
		service: s,
		jwtKey:  jwtKey,
	}

	return &loginjwt
}
