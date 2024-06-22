package service

import (
	"math/rand"
	"quest_service/internal/entity"
	"quest_service/internal/repository"
	"strconv"
)

type QuestService struct {
	questRepo repository.Quest
}

func NewQuestService(questRepo repository.Quest) *QuestService {
	return &QuestService{questRepo: questRepo}
}

func (s *QuestService) CreateQuest(quest *entity.QuestInput) (int, error) {
	return s.questRepo.CreateQuest(quest)
}

func (s *QuestService) GetQuestsAndTasks() ([]entity.Quest, error) {
	return s.questRepo.GetQuestsAndTasks()
}

func (s *QuestService) UpdateQuest(questID int, quest *entity.QuestInput) error {

	if quest.Name != "" {
		err := s.questRepo.UpdateNameQuest(questID, quest)
		if err != nil {
			return err
		}
	}

	err := s.questRepo.UpdateCostQuest(questID, quest)
	if err != nil {
		return err
	}

	return nil
}

func (s *QuestService) DeleteQuest(questID int) error {
	return s.questRepo.DeleteQuest(questID)
}

func (s *QuestService) CreateTestQuestData() error {
	//	Создаём 1000 квестов
	for i := 0; i < 1000; i++ {
		//	 Создаём рандомное количество заданий
		var tasks []entity.TaskInput
		for j := 0; j < 1+rand.Intn(4); j++ {
			tasks = append(tasks, entity.TaskInput{
				Name:       "task" + "_" + strconv.Itoa(i) + "_" + strconv.Itoa(j),
				IsReusable: false,
				Cost:       rand.Intn(250),
			})
		}

		_, err := s.questRepo.CreateQuest(&entity.QuestInput{
			Name:  "quest" + "_" + strconv.Itoa(i) + "_" + strconv.Itoa(len(tasks)),
			Cost:  rand.Intn(1000),
			Tasks: tasks,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
