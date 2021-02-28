package dao

import (
	"context"

	"github.com/Yangruipis/gotd/pkg/core"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
)

type Task struct {
	db *gorm.DB
}

func NewTaskManager(db *gorm.DB) *Task {
	return &Task{
		db: db,
	}
}

var _ core.TaskManager = (*Task)(nil)

func (m *Task) CreateOrUpdate(ctx context.Context, task *core.Task) (*core.Task, error) {
	log.Info().Msgf("task: +v", task)
	if err := m.db.AutoMigrate(task).Error; err != nil {
		return task, err
	}
	err := m.db.Create(task).Error
	return task, err
}

func (m *Task) Get(ctx context.Context, id string) (*core.Task, error) {
	return nil, nil
}

func (m *Task) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *Task) List(ctx context.Context, filter core.TaskFilterParam) ([]*core.Task, error) {
	return nil, nil
}
