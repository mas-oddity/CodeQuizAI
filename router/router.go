package router

import (
	"CodeQuizAI/controllers"
	"CodeQuizAI/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/api/auth/login", controllers.Login)
	r.POST("/api/auth/register", controllers.Register)

	r.GET("/api/users", middlewares.AuthMiddleware(), middlewares.AdminMiddleware(), controllers.ListUsers)
	r.PUT("/api/users/:id", middlewares.AuthMiddleware(), controllers.UpdateUser)

	questionGroup := r.Group("api/questions", middlewares.AuthMiddleware())
	questionGroup.POST("/generate", controllers.GenerateQuestions)
	questionGroup.POST("/confirm", controllers.ConfirmQuestions)
	questionGroup.GET("", controllers.GetQuestions)
	questionGroup.PUT("/:id", controllers.UpdateQuestion)
	questionGroup.DELETE("/:id", controllers.DeleteQuestion)

	paperGroup := r.Group("api/papers", middlewares.AuthMiddleware())
	paperGroup.GET("", controllers.GetPapers)
	paperGroup.POST("", controllers.CreatePaper)
	paperGroup.GET("/:id", controllers.GetPaperDetail)
	paperGroup.DELETE("/:id", controllers.DeletePaper)
	paperGroup.POST("/:id/questions", controllers.AddQuestionsToPaper)
	paperGroup.DELETE("/:id/questions/:questionID", controllers.RemoveQuestionFromPaper)
	paperGroup.PUT("/:id/questions/order", controllers.UpdateQuestionOrder)
	paperGroup.PUT("/:id", controllers.UpdatePaper)

	r.GET("/api/statistics/user/:id", middlewares.AuthMiddleware(), controllers.GetUserStatistics)
	r.GET("/api/statistics/overview", middlewares.AuthMiddleware(), middlewares.AdminMiddleware(), controllers.GetStatisticsOverview)
	return r
}
