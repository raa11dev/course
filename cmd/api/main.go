package main

import (
	"github.com/gin-gonic/gin"
	"github.com/raa11dev/course/internal/database"
	"github.com/raa11dev/course/internal/exercise"
	"github.com/raa11dev/course/internal/middleware"
	"github.com/raa11dev/course/internal/user"
)

func main() {
	route := gin.Default()
	route.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	db := database.NewDatabaseConn()

	exerciseService := exercise.NewExerciseService(db)
	userService := user.NewUserService(db)

	// exercises
	route.POST("/exercises", middleware.Authentication(userService), exerciseService.CreateExercise)
	route.GET("/exercises/:id", middleware.Authentication(userService), exerciseService.GetExercise)
	route.GET("/exercises/:id/score", middleware.Authentication(userService), exerciseService.GetUserScore)
	route.POST("/exercises/:id/questions", middleware.Authentication(userService), exerciseService.CreateQuestions)
	route.POST("/exercises/:id/questions/:id2/answer", middleware.Authentication(userService), exerciseService.CreateAnswer)

	// user
	route.POST("/register", userService.Register)
	route.POST("/login", userService.Login)
	route.Run(":8000")
}
