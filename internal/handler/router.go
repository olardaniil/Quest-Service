package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	_ "quest_service/docs"
	"quest_service/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes(port string) {
	router := gin.New()
	gin.SetMode(gin.ReleaseMode)
	//
	api := router.Group("/api")
	{
		users := api.Group("/users")
		{
			// Создание пользователя
			users.POST("/", h.CreateUser)

			balance := users.Group(":id/balance")
			{
				// Получение баланса
				balance.GET("/", h.GetBalance)
			}
		}

		quests := api.Group("/quests")
		{
			// Создание квеста
			quests.POST("/", h.CreateQuest)
			//	Получение квестов
			quests.GET("/", h.GetQuests)
			//	Обновление квеста
			quests.PUT("/:id", h.UpdateQuest)
			//	Удаление квеста
			quests.DELETE("/:id", h.DeleteQuest)
		}

		tasks := api.Group("/tasks")
		{
			// Создание задания
			tasks.POST("/", h.CreateTask)
			//	Обновление задания
			tasks.PUT("/:id", h.UpdateTask)
			//	Удаление задания
			tasks.DELETE("/:id", h.DeleteTask)
		}

		taskProgress := api.Group("/task-progress")
		{
			// Завершение квеста
			taskProgress.POST("/", h.TaskCompletion)
		}

		//	Добавление тестовых данных
		quests.POST("/test", h.CreateTestQuestData)
	}
	router.GET("api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	fmt.Println("Swagger docs: http://127.0.0.1:" + port + "/api/swagger/")
	err := router.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
