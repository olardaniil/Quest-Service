package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"quest_service/internal/entity"
	"time"
)

type TaskRepo struct {
	db *sqlx.DB
}

func NewTaskRepo(db *sqlx.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) GetTaskByID(taskID int) (*entity.Task, error) {
	var task entity.Task
	taskQuery := `SELECT * FROM tasks WHERE id = $1`
	err := r.db.Get(&task, taskQuery, taskID)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepo) GetCountTaskProgress(task *entity.TaskProgress) (int, error) {
	var countTaskProgress int
	countTaskProgressQuery := fmt.Sprintf("SELECT count(*) FROM tasks_complete WHERE user_id = $1 AND task_id = $2")
	err := r.db.Get(&countTaskProgress, countTaskProgressQuery, task.UserID, task.TaskID)
	if err != nil {
		return 0, err
	}
	return countTaskProgress, nil
}

func (r *TaskRepo) GetTaskStatusesByQuestAndUser(questID, userID int) ([]entity.TaskStatus, error) {
	query := `
        SELECT t.id, CASE WHEN tp.task_id IS NULL THEN false ELSE true END AS is_completed
        FROM tasks t
        LEFT JOIN tasks_complete tp ON t.id = tp.task_id AND tp.user_id = $1
        WHERE t.quest_id = $2
    `
	rows, err := r.db.Query(query, userID, questID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taskStatuses []entity.TaskStatus
	for rows.Next() {
		var taskStatus entity.TaskStatus
		if err := rows.Scan(&taskStatus.TaskID, &taskStatus.IsCompleted); err != nil {
			return nil, err
		}
		taskStatuses = append(taskStatuses, taskStatus)
	}
	return taskStatuses, nil
}

func (r *TaskRepo) TaskCompletion(userID, taskID, taskCost, questID int, isLastTask bool) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	// Записываем данные о выполнении задания
	taskCompletionQuery := `INSERT INTO tasks_complete (user_id, task_id, completed_at) values ($1, $2, $3)`
	_, err = tx.Exec(taskCompletionQuery, userID, taskID, time.Now())
	if err != nil {
		tx.Rollback()
		return err
	}

	if isLastTask == true {
		// Завершаем квест
		questCompletionQuery := `INSERT INTO quests_complete (user_id, quest_id, completed_at) values ($1, $2, $3)`
		_, err = tx.Exec(questCompletionQuery, userID, questID, time.Now())
		if err != nil {
			tx.Rollback()
			return err
		}

		// Обновляем баланс пользователя
		updateBalanceQuery := `UPDATE users SET balance = balance + $2 + (SELECT cost FROM quests WHERE id = $3) WHERE id = $1;`
		_, err = tx.Exec(updateBalanceQuery, userID, taskCost, questID)
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
		return nil
	}

	// Обновляем баланс пользователя
	updateBalanceQuery := `UPDATE users SET balance = balance + $2  WHERE id = $1;`
	_, err = tx.Exec(updateBalanceQuery, userID, taskCost)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil

}

func (r *TaskRepo) CreateTask(task *entity.TaskInput) (int, error) {
	var taskID int
	taskQuery := `INSERT INTO tasks (name, cost, quest_id, is_reusable) values ($1, $2, $3, $4) RETURNING id`
	log.Println(task.QuestID)
	err := r.db.Get(&taskID, taskQuery, task.Name, task.Cost, task.QuestID, task.IsReusable)
	if err != nil {
		return 0, err
	}
	return taskID, nil
}

func (r *TaskRepo) UpdateNameTask(taskID int, name string) error {
	_, err := r.db.Exec("UPDATE tasks SET name = $1 WHERE id = $2", name, taskID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepo) UpdateCostTask(taskID int, cost int) error {
	_, err := r.db.Exec("UPDATE tasks SET cost = $1 WHERE id = $2", cost, taskID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepo) UpdateIsReusableTask(taskID int, isReusable bool) error {
	_, err := r.db.Exec("UPDATE tasks SET is_reusable = $1 WHERE id = $2", isReusable, taskID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepo) DeleteTask(taskID int) error {
	// Удаление задания
	_, err := r.db.Exec("UPDATE tasks SET deleted_at = NOW() WHERE id = $1", taskID)
	if err != nil {
		return err
	}

	return nil
}
