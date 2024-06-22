package handler

import (
	"github.com/gin-gonic/gin"
	"quest_service/internal/entity"
	"strconv"
)

// @Summary		Завершение задачи
// @Tags			tasks
// @Description	Завершение задачи. Задача может быть выполнена несколько раз - зависит от параметра is_reusable.
// @ID				post-tasks-progress
// @Accept			json
// @Produce		json
// @Param			input			body		entity.TaskProgress	true	"body"
// @Success		200				{object}	Response
// @Failure		400,401,403,404	{object}	Response
// @Failure		500				{object}	Response
// @Failure		default			{object}	Response
// @Router			/task-progress/ [post]
func (h *Handler) TaskCompletion(ctx *gin.Context) {
	var input entity.TaskProgress
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
	// Завершение задания
	err = h.services.Task.TaskCompletion(&input)
	if err != nil {
		if err.Error() == "Вы уже выполнили это задание" || err.Error() == "Задание не найдено" {
			resp := Response{
				Message: err.Error(),
			}
			resp.Send(ctx, 200)
			return
		}
		resp := Response{
			Message: "Не удалось завершить задание",
		}
		resp.SendError(ctx, err, 500)
		return
	}
	// Отправка ответа
	resp := Response{
		Message: "Задание успешно завершено",
	}
	resp.Send(ctx, 200)
	return
}

// @Summary		Создание задания
// @Tags			tasks
// @Description	Создание задания
// @ID				post-tasks
// @Accept			json
// @Produce		json
// @Param			input			body		entity.TaskInput	true	"body"
// @Success		200				{object}	Response
// @Failure		400,401,403,404	{object}	Response
// @Failure		500				{object}	Response
// @Failure		default			{object}	Response
// @Router			/tasks/ [post]
func (h *Handler) CreateTask(ctx *gin.Context) {
	// 	Получение тела запроса
	var input entity.TaskInput
	if err := ctx.BindJSON(&input); err != nil {
		resp := Response{
			Message: "Неверное тело запроса",
		}
		resp.SendError(ctx, err, 400)
		return
	}
	// Валидация
	err := input.ValidateForCreate()
	if err != nil {
		resp := Response{
			Message: err.Error(),
		}
		resp.Send(ctx, 400)
		return
	}
	// Создание задания
	taskID, err := h.services.Task.CreateTask(&input)
	if err != nil {
		resp := Response{
			Message: "Не удалось создать задание",
		}
		resp.SendError(ctx, err, 500)
		return
	}
	// Отправка ответа
	resp := Response{
		Message: "Задание успешно создано",
		Details: taskID,
	}
	resp.Send(ctx, 200)
	return
}

// @Summary		Обновление задания
// @Tags			tasks
// @Description	Обновление задания
// @ID				put-tasks
// @Accept			json
// @Produce		json
// @Param			id				path		int					true	"ID квеста"
// @Param			input			body		entity.TaskInput	true	"body"
// @Success		200				{object}	Response
// @Failure		400,401,403,404	{object}	Response
// @Failure		500				{object}	Response
// @Failure		default			{object}	Response
// @Router			/tasks/{id} [put]
func (h *Handler) UpdateTask(ctx *gin.Context) {
	// Получение questID из параметров запроса
	questID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		resp := Response{
			Message: "Неверный ID квеста",
		}
		resp.SendError(ctx, err, 400)
		return
	}
	// 	Получение тела запроса
	var input entity.TaskInput
	if err = ctx.BindJSON(&input); err != nil {
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
	// Обновление задания
	err = h.services.Task.UpdateTask(questID, input.Name, input.Cost, input.IsReusable)
	if err != nil {
		resp := Response{
			Message: "Не удалось обновить задание",
		}
		resp.SendError(ctx, err, 500)
		return
	}
	// Отправка ответа
	resp := Response{
		Message: "Задание успешно обновлено",
	}
	resp.Send(ctx, 200)
	return
}

// @Summary		Удаление задания
// @Tags			tasks
// @Description	Удаление задания
// @ID				delete-task
// @Accept			json
// @Produce		json
// @Param			id				path		int	true	"ID задания"
// @Success		200				{object}	Response
// @Failure		400,401,403,404	{object}	Response
// @Failure		500				{object}	Response
// @Failure		default			{object}	Response
// @Router			/tasks/{id} [delete]
func (h *Handler) DeleteTask(ctx *gin.Context) {
	// Получение questID из параметров запроса
	questID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		resp := Response{
			Message: "Неверный ID квеста",
		}
		resp.SendError(ctx, err, 400)
		return
	}
	// Удаление задания
	err = h.services.Task.DeleteTask(questID)
	if err != nil {
		resp := Response{
			Message: "Не удалось удалить задание",
		}
		resp.SendError(ctx, err, 500)
		return
	}
	// Отправка ответа
	resp := Response{
		Message: "Задание успешно удалено",
	}
	resp.Send(ctx, 200)
	return
}
