package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"quest_service/internal/entity"
	"strconv"
	"time"
)

// @Summary		Создание квеста
// @Tags			quests
// @Description	Создание квеста. В квесте может быть несколько задач. Каждая задача может быть выполнена один или несколько раз в квесте - зависит от параметра is_reusable.
// @ID				post-quests
// @Accept			json
// @Produce		json
// @Param			input			body		entity.QuestInput	true	"body"
// @Success		200				{object}	Response
// @Failure		400,401,403,404	{object}	Response
// @Failure		500				{object}	Response
// @Failure		default			{object}	Response
// @Router			/quests/ [post]
func (h *Handler) CreateQuest(ctx *gin.Context) {
	var input entity.QuestInput
	// Получение тела запроса
	if err := ctx.BindJSON(&input); err != nil {
		resp := Response{
			Message: "Неверное тело запроса",
		}
		resp.SendError(ctx, err, 400)
		return
	}
	// Валидация
	err := input.Validate()
	if err != nil {
		resp := Response{
			Message: err.Error(),
		}
		resp.Send(ctx, 400)
		return
	}
	// Создание квеста
	_, err = h.services.Quest.CreateQuest(&input)
	if err != nil {
		resp := Response{
			Message: "Не удалось создать квест",
		}
		resp.SendError(ctx, err, 500)
		return
	}
	// Отправка ответа
	resp := Response{
		Message: "Квест успешно создан",
	}
	resp.Send(ctx, 201)
	return
}

// @Summary		Получить квесты
// @Tags			quests
// @Description	Получить квесты
// @ID				get-quests
// @Accept			json
// @Produce		json
// @Success		200				{object}	Response
// @Failure		400,401,403,404	{object}	Response
// @Failure		500				{object}	Response
// @Failure		default			{object}	Response
// @Router			/quests/ [get]
func (h *Handler) GetQuests(ctx *gin.Context) {
	// Получение квестов и их заданий
	quests, err := h.services.Quest.GetQuestsAndTasks()
	if err != nil {
		resp := Response{
			Message: "Не удалось получить квесты",
		}
		resp.SendError(ctx, err, 500)
		return
	}
	// Отправка ответа
	resp := Response{
		Message: "Квесты",
		Details: quests,
	}
	resp.Send(ctx, 200)
	return
}

// @Summary		Обновление квеста
// @Tags			quests
// @Description	Обновление квеста
// @ID				put-quests
// @Accept			json
// @Produce		json
// @Param			id				path		int					true	"ID квеста"
// @Param			input			body		entity.QuestInputForUpdate	true	"body"
// @Success		200				{object}	Response
// @Failure		400,401,403,404	{object}	Response
// @Failure		500				{object}	Response
// @Failure		default			{object}	Response
// @Router			/quests/{id} [put]
func (h *Handler) UpdateQuest(ctx *gin.Context) {
	// Получение questID из параметров запроса
	questID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		resp := Response{
			Message: "Неверный ID квеста",
		}
		resp.SendError(ctx, err, 400)
		return
	}
	var input entity.QuestInput
	// Получение тела запроса
	if err := ctx.BindJSON(&input); err != nil {
		resp := Response{
			Message: "Неверное тело запроса",
		}
		resp.SendError(ctx, err, 400)
		return
	}
	// Валидация
	err = input.ValidateForUpdate()
	if err != nil {
		resp := Response{
			Message: err.Error(),
		}
		resp.Send(ctx, 400)
		return
	}
	log.Printf("input: %v", input)
	// Обновление квеста
	err = h.services.Quest.UpdateQuest(questID, &input)
	if err != nil {
		resp := Response{
			Message: "Не удалось обновить квест",
		}
		resp.SendError(ctx, err, 500)
		return
	}
	// Отправка ответа
	resp := Response{
		Message: "Квест успешно обновлен",
	}
	resp.Send(ctx, 200)
	return
}

// @Summary		Удаление квеста
// @Tags			quests
// @Description	Удаление квеста
// @ID				delete-quests
// @Accept			json
// @Produce		json
// @Param			id				path		int	true	"ID квеста"
// @Success		200				{object}	Response
// @Failure		400,401,403,404	{object}	Response
// @Failure		500				{object}	Response
// @Failure		default			{object}	Response
// @Router			/quests/{id} [delete]
func (h *Handler) DeleteQuest(ctx *gin.Context) {
	// Получение questID из параметров запроса
	questID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		resp := Response{
			Message: "Неверный ID квеста",
		}
		resp.SendError(ctx, err, 400)
		return
	}
	// Удаление квеста
	err = h.services.Quest.DeleteQuest(questID)
	if err != nil {
		resp := Response{
			Message: "Не удалось удалить квест",
		}
		resp.SendError(ctx, err, 500)
		return
	}
	// Отправка ответа
	resp := Response{
		Message: "Квест успешно удален",
	}
	resp.Send(ctx, 200)
	return
}

// @Summary		Создание тестовых данных
// @Tags			quests
// @Description	Создание тестовых данных
// @ID				post-quests-test
// @Accept			json
// @Produce		json
// @Success		200				{object}	Response
// @Failure		400,401,403,404	{object}	Response
// @Failure		500				{object}	Response
// @Failure		default			{object}	Response
// @Router			/quests/test [post]
func (h *Handler) CreateTestQuestData(ctx *gin.Context) {
	timeStart := time.Now()
	err := h.services.Quest.CreateTestQuestData()
	if err != nil {
		resp := Response{
			Message: "Не удалось создать тестовые данные",
		}
		resp.SendError(ctx, err, 500)
		return
	}
	log.Printf("Время выполнения: %v", time.Since(timeStart))
	// Отправка ответа
	resp := Response{
		Message: "Тестовые данные успешно созданы",
	}
	resp.Send(ctx, 200)
	return

}
