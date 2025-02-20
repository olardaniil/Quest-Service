package handler

import (
	"github.com/gin-gonic/gin"
	"quest_service/internal/entity"
	"strconv"
)

// @Summary		Создание пользователя
// @Tags			users
// @Description	Создание пользователя
// @ID				post-users
// @Accept			json
// @Produce		json
// @Param			input			body		entity.UserInput	true	"body"
// @Success		200				{object}	Response
// @Failure		400,401,403,404	{object}	Response
// @Failure		500				{object}	Response
// @Failure		default			{object}	Response
// @Router			/users/ [post]
func (h *Handler) CreateUser(ctx *gin.Context) {
	var input entity.UserInput
	// Получение тела запроса
	if err := ctx.BindJSON(&input); err != nil {
		resp := Response{
			Message: "Неверное тело запроса",
		}
		resp.SendError(ctx, err, 400)
		return
	}
	// Валидация
	if err := input.Validate(); err != nil {
		resp := Response{
			Message: err.Error(),
		}
		resp.Send(ctx, 400)
		return
	}
	// Создание пользователя
	userID, err := h.services.User.CreateUser(&input)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
			resp := Response{
				Message: "Пользователь с таким именем уже существует",
			}
			resp.Send(ctx, 409)
			return
		}
		resp := Response{
			Message: "Не удалось создать пользователя",
		}
		resp.SendError(ctx, err, 500)
		return
	}
	// Отправка ответа
	resp := Response{
		Message: "Пользователь успешно создан",
		Details: entity.User{
			ID:       userID,
			UserName: input.UserName,
		},
	}
	resp.Send(ctx, 201)
	return
}

// @Summary		Получить баланс пользователя
// @Tags			users
// @Description	Получить баланс пользователя
// @ID				get-users-id-balance
// @Accept			json
// @Produce		json
// @Param			user_id			path		int	true	"user_id"
// @Success		200				{object}	Response
// @Failure		400,401,403,404	{object}	Response
// @Failure		500				{object}	Response
// @Failure		default			{object}	Response
// @Router			/users/{user_id}/balance [get]
func (h *Handler) GetBalance(ctx *gin.Context) {
	// Получение userID из параметров запроса
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		resp := Response{
			Message: "Неверный ID пользователя",
		}
		resp.SendError(ctx, err, 400)
		return
	}
	// Получение баланса пользователя
	balance, tasks, err := h.services.User.GetBalanceAndHistoryTasks(userID)
	if err != nil {
		resp := Response{
			Message: "Не удалось получить баланс пользователя",
		}
		resp.SendError(ctx, err, 500)
		return
	}
	// Отправка ответа
	resp := Response{
		Message: "Баланс пользователя",
		Details: map[string]interface{}{
			"balance": balance,
			"tasks":   tasks,
		},
	}
	resp.Send(ctx, 200)
}
