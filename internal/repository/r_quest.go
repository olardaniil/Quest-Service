package repository

import (
	"github.com/jmoiron/sqlx"
	"quest_service/internal/entity"
	"time"
)

type QuestRepo struct {
	db *sqlx.DB
}

func NewQuestRepo(db *sqlx.DB) *QuestRepo {
	return &QuestRepo{db: db}
}

func (r *QuestRepo) CreateQuest(quest *entity.QuestInput) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var questID int
	createQuestQuery := `INSERT INTO quests (name, cost, created_at) values ($1, $2, $3) RETURNING id`

	row := tx.QueryRow(createQuestQuery, quest.Name, quest.Cost, time.Now())
	err = row.Scan(&questID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createTaskQuery := `INSERT INTO tasks (quest_id, name, cost) values ($1, $2, $3)`
	for _, task := range quest.Tasks {
		_, err = tx.Exec(createTaskQuery, questID, task.Name, task.Cost)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	return questID, tx.Commit()
}

func (r *QuestRepo) GetQuestsAndTasks() ([]entity.Quest, error) {
	type QuestWithTasks struct {
		QuestID        int    `json:"quest_id,omitempty"`
		QuestName      string `json:"quest_name,omitempty"`
		QuestCost      int    `json:"quest_cost,omitempty"`
		TaskID         int    `json:"task_id,omitempty"`
		TaskName       string `json:"task_name,omitempty"`
		TaskIsReusable bool   `json:"task_is_reusable,omitempty"`
		TaskCost       int    `json:"task_cost,omitempty"`
	}
	var questWithTasks []QuestWithTasks
	// Получение всех квестов и их заданий
	questsQuery := `SELECT q.id, q.name, q.cost, t.id, t.name, t.is_reusable, t.cost FROM quests q JOIN tasks t ON q.id = t.quest_id WHERE q.deleted_at IS NULL AND t.deleted_at IS NULL`
	rows, err := r.db.Query(questsQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Запись данных в структуру
	for rows.Next() {
		var q QuestWithTasks
		err = rows.Scan(&q.QuestID, &q.QuestName, &q.QuestCost, &q.TaskID, &q.TaskName, &q.TaskIsReusable, &q.TaskCost)
		if err != nil {
			return nil, err
		}
		questWithTasks = append(questWithTasks, q)
	}

	var quests []entity.Quest
	var questsIDs []int

	for _, q := range questWithTasks {
		if !contains(questsIDs, q.QuestID) {
			quests = append(quests, entity.Quest{
				ID:    q.QuestID,
				Name:  q.QuestName,
				Cost:  q.QuestCost,
				Tasks: []entity.Task{},
			})
			questsIDs = append(questsIDs, q.QuestID)
		}
		quests[len(quests)-1].Tasks = append(quests[len(quests)-1].Tasks, entity.Task{
			ID:         q.TaskID,
			Name:       q.TaskName,
			IsReusable: q.TaskIsReusable,
			Cost:       q.TaskCost,
		})
	}

	return quests, nil
}

func (r *QuestRepo) UpdateNameQuest(questID int, quest *entity.QuestInput) error {
	// Обновление названия квеста
	_, err := r.db.Exec("UPDATE quests SET name = $1 WHERE id = $2", quest.Name, questID)
	if err != nil {
		return err
	}

	return nil
}

func (r *QuestRepo) UpdateCostQuest(questID int, quest *entity.QuestInput) error {
	// Обновление стоимости квеста
	_, err := r.db.Exec("UPDATE quests SET cost = $1 WHERE id = $2", quest.Cost, questID)
	if err != nil {
		return err
	}

	return nil
}

func (r *QuestRepo) DeleteQuest(questID int) error {
	// Удаление квеста
	_, err := r.db.Exec("UPDATE quests SET deleted_at = NOW() WHERE id = $1", questID)
	if err != nil {
		return err
	}

	return nil
}

//func (r *QuestRepo)

func contains(ds []int, id int) bool {
	for _, v := range ds {
		if v == id {
			return true
		}
	}
	return false
}
