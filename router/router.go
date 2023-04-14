package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/subkhiyoga/auth-jwt/controller"
	"github.com/subkhiyoga/auth-jwt/database"
	"github.com/subkhiyoga/auth-jwt/repository"
	"github.com/subkhiyoga/auth-jwt/usecase"
)

func Run() {
	db := repository.ConnectDB()
	secret_key := []byte(database.DotEnv("SECRET_KEY"))
	authMiddleware := controller.AuthMiddleware(secret_key)

	// gin router
	router := gin.Default()

	mahasiwaRepo := repository.NewMahasiswaRepo(db)
	loginUsecase := usecase.NewLoginUsecase(mahasiwaRepo)
	loginJwt := controller.NewCredentialsJwt(loginUsecase)

	// set up routes for login
	router.POST("/login", authMiddleware, loginJwt.Login)

	// other routes
	mahasiwaRouter := router.Group("/api/v1/mahasiswa")
	mahasiwaRouter.Use(authMiddleware, loginJwt.Login)

	mahasiwaRouter.GET("/:id/profile", authMiddleware, loginJwt.Login)

	err := router.Run(database.DotEnv("SERVER_PORT"))
	if err != nil {
		log.Fatal(err)
	}
}
