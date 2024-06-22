package service

import (
	"quest_service/internal/entity"
	"quest_service/internal/repository"
)

type User interface {
	CreateUser(user *entity.UserInput) (int, error)
	GetBalanceAndHistoryTasks(userID int) (int, []entity.Task, error)
}

type Quest interface {
	CreateQuest(quest *entity.QuestInput) (int, error)
	GetQuestsAndTasks() ([]entity.Quest, error)
	UpdateQuest(questID int, quest *entity.QuestInput) error
	DeleteQuest(questID int) error
	CreateTestQuestData() error
}

type Task interface {
	TaskCompletion(taskProgress *entity.TaskProgress) error
	CreateTask(task *entity.TaskInput) (int, error)
	UpdateTask(taskID int, name string, cost int, isReusable bool) error
	DeleteTask(taskID int) error
}

type Service struct {
	User
	Quest
	Task
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:  NewUserService(repos.User),
		Quest: NewQuestService(repos.Quest),
		Task:  NewTaskService(repos.Task),
	}
}
