package service

import (
	"errors"
	"go_final_project/internal/models"
	"go_final_project/pkg/api"
	"log"
	"strconv"
	"time"
)

var (
	ErrTaskIDEmpty     = errors.New("идентификатор задачи не указан")
	ErrTaskIDInvalid   = errors.New("идентификатор задачи содержит не правильный формат")
	ErrTaskTitleEmpty  = errors.New("не указан заголовок задачи")
	ErrTaskDateInvalid = errors.New("дата представлена в формате, отличном от " + api.DateFormat)
)

type ITaskStore interface {
	List(int) ([]models.Task, error)
	Get(string) (*models.Task, error)
	Add(*models.Task) (string, error)
	SetTitle(*models.Task) error
	SetComment(*models.Task) error
	SetDate(*models.Task) error
	SetRepeat(*models.Task) error
	Delete(string) error
}

type TaskService struct {
	store ITaskStore
	logs  *log.Logger
}

func NewTaskService(store ITaskStore, logs *log.Logger) *TaskService {
	return &TaskService{
		store: store,
		logs:  logs,
	}
}

func (s *TaskService) Add(task *models.Task) (string, error) {
	if task.Title == "" {
		s.logs.Println(ErrTaskTitleEmpty)

		return "", ErrTaskTitleEmpty
	}

	now := time.Now()

	if task.Date == "" {
		task.Date = api.FormatDate(now)
	}

	taskDate, err := api.ParseDate(task.Date)

	if err != nil {
		s.logs.Println(err)

		return "", ErrTaskDateInvalid
	}

	if task.Repeat != "" {
		date, err := api.NextDate(now, task.Date, task.Repeat)

		if err != nil {
			s.logs.Println(err)

			return "", err
		}

		task.Date = date
	}

	if taskDate.Before(now) {
		task.Date = api.FormatDate(now)
	}

	return s.store.Add(task)
}

func (s *TaskService) List(limit int) ([]models.Task, error) {
	return s.store.List(limit)
}

func (s *TaskService) Get(id string) (*models.Task, error) {
	if id == "" {
		s.logs.Println(ErrTaskIDEmpty)

		return nil, ErrTaskIDEmpty
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		s.logs.Println(err)

		return nil, ErrTaskIDInvalid
	}

	return s.store.Get(id)
}

func (s *TaskService) Edit(task *models.Task) error {
	_, err := s.Get(task.ID)
	if err != nil {
		return err
	}

	if task.Title == "" {
		s.logs.Println(ErrTaskTitleEmpty)

		return ErrTaskTitleEmpty
	}

	_, err = api.ParseDate(task.Date)
	if err != nil {
		s.logs.Println(err)

		return ErrTaskDateInvalid
	}

	now := time.Now()

	if task.Repeat != "" {
		_, err := api.NextDate(now, task.Date, task.Repeat)

		if err != nil {
			s.logs.Println(err)

			return err
		}
	}

	err = s.store.SetTitle(task)
	if err != nil {
		return err
	}

	err = s.store.SetComment(task)
	if err != nil {
		return err
	}

	err = s.store.SetDate(task)
	if err != nil {
		return err
	}

	err = s.store.SetRepeat(task)
	if err != nil {
		return err
	}

	return nil
}

func (s *TaskService) Delete(id string) error {
	task, err := s.Get(id)
	if err != nil {
		return err
	}

	return s.store.Delete(task.ID)
}

func (s *TaskService) Done(id string) error {
	task, err := s.Get(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		return s.Delete(task.ID)
	}

	taskDate, err := api.ParseDate(task.Date)
	if err != nil {
		s.logs.Println(err)

		return err
	}

	task.Date, err = api.NextDate(taskDate, task.Date, task.Repeat)
	if err != nil {
		s.logs.Println(err)

		return err
	}

	return s.store.SetDate(task)
}
